package receivers

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/google/uuid"
	ops "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/operations"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operations"
)

type operationMessageHandler struct {
	database data.Database
	factory  ops.DeploymentOperationFactory
}

func (h *operationMessageHandler) Handle(message *messaging.InvokedOperationMessage, context messaging.MessageHandlerContext) error {
	db := h.database.Instance()
	var invokedOperation *data.InvokedOperation
	db.First(&invokedOperation, uuid.MustParse(message.OperationId))

	log.Printf("operation id: %s", message.OperationId)
	log.Printf("Invoked Operation from DB: %v", invokedOperation)

	operationType, err := operations.Type(invokedOperation.Name)
	if err != nil {
		log.Println("Error getting operation type: ", err)
	}

	operation, err := h.factory.Create(operationType)
	if err != nil {
		log.Println("Error creating operation: ", err)
		return err
	}

	return operation.Invoke(invokedOperation)
}

//region factory

func NewOperationsMessageReceiver(appConfig *config.AppConfig, credential azcore.TokenCredential) messaging.MessageReceiver {
	db := data.NewDatabase(appConfig.GetDatabaseOptions())

	handler := &operationMessageHandler{
		database: db,
		factory:  ops.NewDeploymentOperationFactory(appConfig),
	}

	options := messaging.ServiceBusMessageReceiverOptions{
		MessageReceiverOptions:  messaging.MessageReceiverOptions{QueueName: string(messaging.QueueNameOperations)},
		FullyQualifiedNamespace: appConfig.Azure.GetFullQualifiedNamespace(),
	}
	receiver, err := messaging.NewServiceBusReceiver(handler, credential, options)
	if err != nil {
		log.Fatal(err)
	}
	return receiver
}

//endregion receiver factory
