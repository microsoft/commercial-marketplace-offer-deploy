package deployment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	log "github.com/sirupsen/logrus"
)

type ArmDeployer struct {
	templateType DeploymentType
	client       *armresources.DeploymentsClient
}

func (deployer *ArmDeployer) Type() DeploymentType {
	return deployer.templateType
}

func (deployer *ArmDeployer) getParamsMapFromTemplate(template map[string]interface{}, params map[string]interface{}) map[string]interface{} {
	paramValues := make(map[string]interface{})

	templateParams := getParams(template)
	for k := range templateParams {
		valueMap := make(map[string]interface{})
		templateValueMap := params[k].(map[string]interface{})

		valueMap["value"] = templateValueMap["value"]
		paramValues[k] = valueMap
	}

	return paramValues
}

func (deployer *ArmDeployer) Cancel(ctx context.Context, acd AzureCancelDeployment) (*AzureCancelDeploymentResult, error) {
	log.Trace("Beginning Azure deployment cancellation")

	_, err := deployer.client.Cancel(ctx, acd.ResourceGroupName, acd.DeploymentName, nil)
	if err != nil {
		return nil, err
	}

	return &AzureCancelDeploymentResult{
		CancelSubmitted: true,
	}, nil
}

func (deployer *ArmDeployer) Redeploy(ctx context.Context, ad AzureRedeployment) (*AzureDeploymentResult, error) {
	b, err := json.MarshalIndent(ad, "", "  ")
	if err != nil {
		log.Error(err)
	}
	log.Tracef("Inside Redeploy in deployment package with a value of %s", string(b))

	deployment, err := deployer.client.Get(ctx, ad.ResourceGroupName, ad.DeploymentName, nil)
	if err != nil {
		return nil, err
	}

	if deployment.Properties == nil || deployment.Properties.Parameters == nil {
		return nil, errors.New("deployment.Properties.Parameters  not found")
	}

	params := (*deployment.DeploymentExtended.Properties).Parameters
	if params == nil {
		return nil, errors.New("unable to get the parameters from the deployment")
	}

	castParams := params.(map[string]interface{})
	if castParams == nil {
		return nil, errors.New("unable to cast parameters to map[string]interface{}")
	}

	template, err := deployer.client.ExportTemplate(ctx, ad.ResourceGroupName, ad.DeploymentName, nil)
	if err != nil {
		return nil, err
	}

	if template.Template == nil {
		return nil, errors.New("unable to get the template from the deployment")
	}

	castTemplate := template.Template.(map[string]interface{})
	paramValuesMap := deployer.getParamsMapFromTemplate(castTemplate, castParams)

	log.Debugf("About to call BeginCreateOrUpdate in Redeploy in deployment package with a resourceGroupName of %s and a deploymentName of %s", ad.ResourceGroupName, ad.DeploymentName)
	deploymentPollerResp, err := deployer.client.BeginCreateOrUpdate(
		ctx,
		ad.ResourceGroupName,
		ad.DeploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Template:   castTemplate,
				Parameters: paramValuesMap,
				Mode:       to.Ptr(armresources.DeploymentModeIncremental),
			},
			Tags: ad.Tags,
		},
		nil)

	if err != nil {
		return nil, errors.New("unable to redeploy the deployment")
	}

	resp, err := deploymentPollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot get the create deployment future respone: %v", err)
	}

	mappedResult, err := deployer.mapDeploymentResult(resp)
	if err != nil {
		return nil, err
	}

	return mappedResult, nil
}

func (deployer *ArmDeployer) Begin(ctx context.Context, azureDeployment AzureDeployment) (*BeginAzureDeploymentResult, error) {
	logger := azureDeployment.logger()
	logger.Trace("Beginning Azure deployment")

	poller, err := deployer.client.BeginCreateOrUpdate(
		ctx,
		azureDeployment.ResourceGroupName,
		azureDeployment.DeploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Template:   azureDeployment.Template,
				Parameters: azureDeployment.GetParameters(),
				Mode:       to.Ptr(armresources.DeploymentModeIncremental),
			},
			Tags: azureDeployment.Tags,
		},
		nil)

	if err != nil {
		return nil, fmt.Errorf("the Azure deployment creation failed: %v", err)
	}

	resumeTokenValue, err := poller.ResumeToken()
	if err != nil {
		return nil, fmt.Errorf("cannot get the create deployment future response: %v", err)
	}

	result := &BeginAzureDeploymentResult{
		ResumeToken: ResumeToken{
			SubscriptionId: azureDeployment.SubscriptionId,
			Value:          resumeTokenValue,
		},
	}

	//try and get the correlationId
	info, err := deployer.client.Get(ctx, azureDeployment.ResourceGroupName, azureDeployment.DeploymentName, nil)
	if err == nil {
		result.CorrelationID = info.Properties.CorrelationID
	}

	return result, nil
}

func (deployer *ArmDeployer) Wait(ctx context.Context, resumeToken *ResumeToken) (*AzureDeploymentResult, error) {
	if resumeToken == nil {
		return nil, errors.New("resume token is nil. cannot wait for deployment")
	}

	logger := log.WithFields(log.Fields{
		"subscriptionId": resumeToken.SubscriptionId,
		"resumeToken":    resumeToken.Value,
	})

	logger.Trace("Beginning to wait on Azure deployment")

	poller, err := deployer.client.BeginCreateOrUpdate(ctx, "", "", armresources.Deployment{}, &armresources.DeploymentsClientBeginCreateOrUpdateOptions{
		ResumeToken: resumeToken.Value,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to resume waiting on azure deployment: %v", err)
	}

	response, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get the deployment future respone: %v", err)
	}

	mappedResult, err := deployer.mapDeploymentResult(response)
	if err != nil {
		return nil, err
	}

	return mappedResult, nil
}

func (deployer *ArmDeployer) mapDeploymentResult(resp armresources.DeploymentsClientCreateOrUpdateResponse) (*AzureDeploymentResult, error) {
	var status ExecutionStatus
	deploymentExtended := resp.DeploymentExtended
	provisioningState := *deploymentExtended.Properties.ProvisioningState

	switch provisioningState {
	case armresources.ProvisioningStateSucceeded:
		status = Succeeded
	case armresources.ProvisioningStateCanceled:
		status = Canceled
	default:
		status = Failed
	}

	// make sure response outputs are always there, even if empty
	var responseOutputs map[string]interface{}
	if deploymentExtended.Properties.Outputs != nil {
		responseOutputs = deploymentExtended.Properties.Outputs.(map[string]interface{})
	} else {
		responseOutputs = make(map[string]interface{})
	}

	res := AzureDeploymentResult{}
	if deploymentExtended.Properties.ProvisioningState != nil {
		res.ProvisioningState = string(*deploymentExtended.Properties.ProvisioningState)
	}
	if deploymentExtended.ID != nil {
		res.ID = *deploymentExtended.ID
	}
	if deploymentExtended.Properties.CorrelationID != nil {
		res.CorrelationID = *deploymentExtended.Properties.CorrelationID
	}
	if deploymentExtended.Properties.Duration != nil {
		res.Duration = *deploymentExtended.Properties.Duration
	}
	if deploymentExtended.Properties.Timestamp != nil {
		res.Timestamp = *deploymentExtended.Properties.Timestamp
	}
	res.Status = status
	res.Outputs = responseOutputs
	return &res, nil
}
