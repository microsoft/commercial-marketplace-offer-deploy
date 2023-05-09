package deployment

import (
	"context"
	"errors"
	"fmt"
	"time"
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	log "github.com/sirupsen/logrus"
)

type ExecutionStatus string

const (
	Started           ExecutionStatus = "Started"
	Failed            ExecutionStatus = "Failed"
	PermanentlyFailed ExecutionStatus = "PermanentlyFailed"
	Succeeded         ExecutionStatus = "Succeeded"
	Restart           ExecutionStatus = "Restart"
	Restarted         ExecutionStatus = "Restarted"
	RestartTimedOut   ExecutionStatus = "RestartTimedOut"
	Canceled          ExecutionStatus = "Canceled"
)

type AzureDeployment struct {
	SubscriptionId    string         `json:"subscriptionId"`
	Location          string         `json:"location"`
	ResourceGroupName string         `json:"resourceGroupName"`
	DeploymentName    string         `json:"deploymentName"`
	DeploymentType    DeploymentType `json:"deploymentType"`
	Template          Template       `json:"template"`
	Params            TemplateParams `json:"templateParams"`
	ResumeToken       string         `json:"resumeToken"`
}

type AzureRedeployment struct {
	SubscriptionId    string         `json:"subscriptionId"`
	Location          string         `json:"location"`
	ResourceGroupName string         `json:"resourceGroupName"`
	DeploymentName    string         `json:"deploymentName"`
}

type AzureDeploymentResult struct {
	ID                string                 `json:"id"`
	CorrelationID     string                 `json:"correlationId"`
	Duration          string                 `json:"duration"`
	Timestamp         time.Time              `json:"timestamp"`
	ProvisioningState string                 `json:"provisioningState"`
	Outputs           map[string]interface{} `json:"outputs"`
	Status            ExecutionStatus
}

func (ad *AzureDeployment) GetDeploymentType() DeploymentType {
	return ad.DeploymentType
}

func (ad *AzureDeployment) GetTemplate() map[string]interface{} {
	return ad.Template
}

func (ad *AzureDeployment) GetTemplateParams() map[string]interface{} {
	return ad.Params
}

type Deployer interface {
	Deploy(ctx context.Context, d *AzureDeployment) (*AzureDeploymentResult, error)
	Redeploy(ctx context.Context, d *AzureRedeployment) (*AzureDeploymentResult, error)
}

type ArmTemplateDeployer struct {
	deployerType DeploymentType
}

func (armDeployer *ArmTemplateDeployer) getParamsMapFromTemplate(template map[string]interface{}, params map[string]interface{}) map[string]interface{} { 
	paramValues := make(map[string]interface{})
	

	var templateParams map[string]interface{}
	if template != nil {
		if p, ok := template["parameters"]; ok {
			templateParams = p.(map[string]interface{})
		} else {
			templateParams = template
		}
	}

	for k := range templateParams {
		valueMap := make(map[string]interface{})
		templateValueMap := params[k].(map[string]interface{})

		valueMap["value"] = templateValueMap["value"]
		paramValues[k] = valueMap
	}

	return paramValues
}


func (armDeployer *ArmTemplateDeployer) Redeploy(ctx context.Context, ad *AzureRedeployment) (*AzureDeploymentResult, error) {
	b, err := json.MarshalIndent(ad, "", "  ")
    if err != nil {
        log.Error(err)
    }
	log.Debugf("Inside Redeploy in deployment package with a value of %s", string(b))
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	
	deploymentsClient, err := armresources.NewDeploymentsClient(ad.SubscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}

	deployment, err := deploymentsClient.Get(ctx, ad.ResourceGroupName, ad.DeploymentName, nil)
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

	template, err := deploymentsClient.ExportTemplate(ctx, ad.ResourceGroupName, ad.DeploymentName, nil)
	if err != nil {
		return nil, err
	}

	if template.Template == nil {
		return nil, errors.New("unable to get the template from the deployment")
	}

	castTemplate := template.Template.(map[string]interface{})
	paramValuesMap := armDeployer.getParamsMapFromTemplate(castTemplate, castParams)


	log.Debugf("About to call BeginCreateOrUpdate in Redeploy in deployment package with a resourceGroupName of %s and a deploymentName of %s", ad.ResourceGroupName, ad.DeploymentName)
	deploymentPollerResp, err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		ad.ResourceGroupName,
		ad.DeploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Template:   castTemplate,
				Parameters: paramValuesMap,
				Mode:       to.Ptr(armresources.DeploymentModeIncremental),
			},
		},
		nil)
		
	if err != nil {
		return nil, errors.New("unable to redeploy the deployment")
	}

	resp, err := deploymentPollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot get the create deployment future respone: %v", err)
	}

	mappedResult, err := armDeployer.mapDeploymentResult(resp)
	if err != nil {
		return nil, err
	}

	return mappedResult, nil
}

func (armDeployer *ArmTemplateDeployer) Deploy(ctx context.Context, ad *AzureDeployment) (*AzureDeploymentResult, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Error(err)
	}

	deploymentsClient, err := armresources.NewDeploymentsClient(ad.SubscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}

	log.Error("About to Create a deployment")

	var templateParams map[string]interface{}
	if ad.Params != nil {
		if p, ok := ad.Params["parameters"]; ok {
			templateParams = p.(map[string]interface{})
		} else {
			templateParams = ad.Params
		}
	}

	deploymentPollerResp, err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		ad.ResourceGroupName,
		ad.DeploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Template:   ad.Template,
				Parameters: templateParams,
				Mode:       to.Ptr(armresources.DeploymentModeIncremental),
			},
		},
		nil)

	if err != nil {
		return nil, fmt.Errorf("cannot create deployment: %v", err)
	}

	//todo: capture the state of the started deployment
	resp, err := deploymentPollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot get the create deployment future respone: %v", err)
	}

	mappedResult, err := armDeployer.mapDeploymentResult(resp)
	if err != nil {
		return nil, err
	}

	return mappedResult, nil
}

func (armDeployer *ArmTemplateDeployer) mapDeploymentResult(resp armresources.DeploymentsClientCreateOrUpdateResponse) (*AzureDeploymentResult, error) {
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
	if &provisioningState != nil {
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

func CreateNewDeployer(deploymentType DeploymentType) Deployer {
	return &ArmTemplateDeployer{
		deployerType: deploymentType,
	}
}

func (azureDeployment *AzureDeployment) validate() error {
	if len(azureDeployment.SubscriptionId) == 0 {
		return errors.New("subscriptionId is not set on azureDeployment input struct")
	}
	if len(azureDeployment.Location) == 0 {
		return errors.New("location is not set on azureDeployment input struct")
	}
	if len(azureDeployment.ResourceGroupName) == 0 {
		return errors.New("resourceGroupName is not set on azureDeployment input struct")
	}
	if len(azureDeployment.ResourceGroupName) == 0 {
		return errors.New("resourceGroupName is not set on azureDeployment input struct")
	}
	if azureDeployment.Template == nil {
		return errors.New("template is not set on deployment azureDeployment struct")
	}
	// allow params to be empty to support all default params
	return nil
}

func FindResourcesByType(template Template, resourceType string) []string {
	deploymentResources := []string{}
	if template != nil && template["resources"] != nil {
		resources := template["resources"].([]interface{})
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			if resourceMap["type"] != nil && resourceMap["type"].(string) == resourceType {
				deploymentResources = append(deploymentResources, resourceMap["name"].(string))
			}
		}
	}
	return deploymentResources
}

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
	// READ-ONLY; The additional info.
	Info interface{} `json:"info,omitempty" azure:"ro"`

	// READ-ONLY; The additional info type.
	Type *string `json:"type,omitempty" azure:"ro"`
}
