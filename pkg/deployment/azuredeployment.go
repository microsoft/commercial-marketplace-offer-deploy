package deployment

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type AzureDeployment struct {
	subscriptionId string
	location string
	resourceGroupName string
	deploymentName string
	deploymentType DeploymentType
	template Template
	params TemplateParams
	resumeToken string
}

func (ad *AzureDeployment) GetDeploymentType() DeploymentType {
	return ad.deploymentType
}

func (ad *AzureDeployment) GetTemplate() map[string]interface{} {
	return ad.template
}

func (ad *AzureDeployment) GetTemplateParams() map[string]interface{} {
	return ad.params
}

type AzureDeploymentResult struct {
	code string
}

type Deployer interface {
	Deploy(d *AzureDeployment) (*AzureDeploymentResult, error)
}

type ArmTemplateDeployer struct {
	deployerType DeploymentType
}

func (armDeployer *ArmTemplateDeployer) Deploy(ad *AzureDeployment) (*AzureDeploymentResult, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Print(err)
	}
	ctx := context.Background()
	deploymentsClient, err := armresources.NewDeploymentsClient(ad.subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}

	log.Printf("About to Create a deployment")

	deploymentPollerResp, err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		ad.resourceGroupName,
		ad.deploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Template:   ad.template,
				Parameters: ad.params,
				Mode:       to.Ptr(armresources.DeploymentModeIncremental),
			},
		},
		nil)

	if err != nil {
		return nil, fmt.Errorf("cannot create deployment: %v", err)
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

func (armDeployer *ArmTemplateDeployer) mapDeploymentResult(resp armresources.DeploymentsClientCreateOrUpdateResponse)	(*AzureDeploymentResult, error) {
	result := &AzureDeploymentResult{
		code: "Success",
	}
	return result, nil
}

func CreateNewDeployer(deployment AzureDeployment) Deployer {
	return &ArmTemplateDeployer{
		deployerType: deployment.deploymentType,
	}
}

func (azureDeployment *AzureDeployment) validate() (error) {
	if len(azureDeployment.subscriptionId) == 0 {
		return errors.New("subscriptionId is not set on azureDeployment input struct")
	}
	if len(azureDeployment.location) == 0 {
		return errors.New("location is not set on azureDeployment input struct")
	}
	if len(azureDeployment.resourceGroupName) == 0 {
		return errors.New("resourceGroupName is not set on azureDeployment input struct")
	}
	if len(azureDeployment.resourceGroupName) == 0 {
		return errors.New("resourceGroupName is not set on azureDeployment input struct")
	}
	if azureDeployment.template == nil {
		return errors.New("template is not set on deployment azureDeployment struct")
	}
	// allow params to be empty to support all default params
	return nil
}

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
	// READ-ONLY; The additional info.
	Info interface{} `json:"info,omitempty" azure:"ro"`

	// READ-ONLY; The additional info type.
	Type *string `json:"type,omitempty" azure:"ro"`
}
