package messaging

import (
	"context"
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
)

type OperationsHandler struct {
	running bool
}

func NewOperationsHandler() *OperationsHandler {
	return &OperationsHandler{
		running: false,
	}
}

func (h *OperationsHandler) Handle(ctx context.Context, message *azservicebus.ReceivedMessage) error {
	var azureDeployment deployment.AzureDeployment
	err := json.Unmarshal(message.Body, &azureDeployment)
	if err != nil {
		return err
	}
	h.deploy(ctx, &azureDeployment)
	return nil
}

func (h *OperationsHandler) deploy(ctx context.Context, azureDeployment *deployment.AzureDeployment) error  {
	h.running = true
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil
	}
	deploymentsClient, err := armresources.NewDeploymentsClient(azureDeployment.SubscriptionId, cred, nil)
	if err != nil {
		return nil
	}
	deploymentsClient.BeginCreateOrUpdate(
		ctx, 
		azureDeployment.ResourceGroupName, 
		azureDeployment.DeploymentName, 
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Template: azureDeployment.Template,
				Parameters: azureDeployment.Params,
				Mode: to.Ptr(armresources.DeploymentModeIncremental),
			},
		},
		nil,
	)
	return nil
}



