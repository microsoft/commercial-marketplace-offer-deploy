package messaging

import (
	"context"
	"encoding/json"
	"strings"

	//"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	//"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	//"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type OperationsHandler struct {
	running bool
	database data.Database
}

func NewOperationsHandler() *OperationsHandler {
	return &OperationsHandler{
		running: false,
	}
}

func (h *OperationsHandler) Handle(ctx context.Context, message *azservicebus.ReceivedMessage) (*deployment.AzureDeploymentResult, error) {
	var operation data.InvokedOperation
	err := json.Unmarshal(message.Body, &operation)
	if err != nil {
		return nil, err
	}
	db := h.database.Instance()

	deployment := &data.Deployment{}
	db.First(&deployment, operation.DeploymentId)

	startedStatus := strings.Replace(string(events.DeploymentStartedEventType), "deployment.", "", 1)
	caser := cases.Title(language.English)
	deployment.Status = caser.String(startedStatus)
	db.Save(deployment)

	azureDeployment := h.mapAzureDeployment(deployment, &operation)
	res, err := h.Deploy(ctx, azureDeployment)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (h *OperationsHandler) mapAzureDeployment(d *data.Deployment, io *data.InvokedOperation) *deployment.AzureDeployment {
	return &deployment.AzureDeployment{
		SubscriptionId: d.SubscriptionId,
		ResourceGroupName: d.ResourceGroup,
		DeploymentName: io.DeploymentName,
		Template: d.Template,
		Params: io.Params,
	}
}

func (h *OperationsHandler) Deploy(ctx context.Context, azureDeployment *deployment.AzureDeployment) (*deployment.AzureDeploymentResult, error)  {
	
	
	return deployment.Create(*azureDeployment)
	// h.running = true
	// cred, err := azidentity.NewDefaultAzureCredential(nil)
	// if err != nil {
	// 	return nil
	// }
	// deploymentsClient, err := armresources.NewDeploymentsClient(azureDeployment.SubscriptionId, cred, nil)
	// if err != nil {
	// 	return nil
	// }
	// deploymentsClient.BeginCreateOrUpdate(
	// 	ctx, 
	// 	azureDeployment.ResourceGroupName, 
	// 	azureDeployment.DeploymentName, 
	// 	armresources.Deployment{
	// 		Properties: &armresources.DeploymentProperties{
	// 			Template: azureDeployment.Template,
	// 			Parameters: azureDeployment.Params,
	// 			Mode: to.Ptr(armresources.DeploymentModeIncremental),
	// 		},
	// 	},
	// 	nil,
	// )
	// return nil
}



