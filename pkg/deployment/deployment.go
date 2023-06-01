package deployment

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type DeploymentType int64

const (
	DeploymentTypeARM DeploymentType = iota
	// terraform not supported yet. stubbed only
	DeploymentTypeTerraform
)

// deploys templates to Azure
type Deployer interface {
	Begin(ctx context.Context, azureDeployment AzureDeployment) (*BeginAzureDeploymentResult, error)
	Wait(ctx context.Context, resumeToken *ResumeToken) (*AzureDeploymentResult, error)

	Redeploy(ctx context.Context, d AzureRedeployment) (*AzureDeploymentResult, error)

	// gets the template type the deployer is set to
	Type() DeploymentType
}

func NewDeployer(deploymentType DeploymentType, subscriptionId string) (Deployer, error) {
	client, err := createAzureDeploymentClient(subscriptionId)
	if err != nil {
		return nil, err
	}

	if deploymentType == DeploymentTypeARM {
		return &ArmDeployer{
			templateType: DeploymentTypeARM,
			client:       client,
		}, nil
	}

	return nil, errors.New("template type not supported")
}

func createAzureDeploymentClient(subscriptionId string) (*armresources.DeploymentsClient, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	client, err := armresources.NewDeploymentsClient(subscriptionId, cred, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
