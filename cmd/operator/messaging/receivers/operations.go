package receivers

import (
	"context"
	"log"
	"strings"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type operationsMessageHandler struct {
	running  bool
	database data.Database
}

func (h *operationsMessageHandler) Handle(message *messaging.InvokedOperationMessage, context messaging.MessageHandlerContext) error {
	db := h.database.Instance()

	var operation data.InvokedOperation
	db.First(&operation, message.OperationId)

	log.Println("Unmarshalled message: ", operation)
	pulledOperationId := operation.ID
	log.Println("Pulled operationId: ", pulledOperationId)
	log.Println("Pulled params: ", operation.Parameters)

	deployment := &data.Deployment{}
	db.First(&deployment, operation.DeploymentId)
	log.Println("Found deployment: ", deployment)

	startedStatus := strings.Replace(string(events.DeploymentStartedEventType), "deployment.", "", 1)
	caser := cases.Title(language.English)
	deployment.Status = caser.String(startedStatus)

	db.Save(deployment)
	log.Println("Updated deployment: ", deployment)

	azureDeployment := h.mapAzureDeployment(deployment, &operation)
	log.Println("Mapped deployment: ", azureDeployment)
	log.Println("Calling deployment.Create")

	go func() {
		_, err := h.deploy(context.Context(), azureDeployment)
		if err != nil {
			log.Println("Error calling deployment.Create: ", err)
		}
	}()

	return nil
}

func (h *operationsMessageHandler) mapAzureDeployment(d *data.Deployment, io *data.InvokedOperation) *deployment.AzureDeployment {
	return &deployment.AzureDeployment{
		SubscriptionId:    d.SubscriptionId,
		ResourceGroupName: d.ResourceGroup,
		DeploymentName:    d.GetAzureDeploymentName(),
		Template:          d.Template,
		Params:            io.Parameters,
	}
}

func (h *operationsMessageHandler) deploy(ctx context.Context, azureDeployment *deployment.AzureDeployment) (*deployment.AzureDeploymentResult, error) {
	return deployment.Create(*azureDeployment)
}

//region factory

func newOperationsMessageHandler(db data.Database) *operationsMessageHandler {
	return &operationsMessageHandler{
		running:  false,
		database: db,
	}
}

func NewOperationsMessageReceiver(appConfig *config.AppConfig) messaging.MessageReceiver {
	db := data.NewDatabase(appConfig.GetDatabaseOptions())

	handler := newOperationsMessageHandler(db)
	options := messaging.ServiceBusMessageReceiverOptions{
		MessageReceiverOptions:  messaging.MessageReceiverOptions{QueueName: string(messaging.QueueNameOperations)},
		FullyQualifiedNamespace: appConfig.Azure.GetFullQualifiedNamespace(),
	}
	receiver, err := messaging.NewServiceBusReceiver(handler, options)
	if err != nil {
		log.Fatal(err)
	}
	return receiver
}

//endregion receiver factory
