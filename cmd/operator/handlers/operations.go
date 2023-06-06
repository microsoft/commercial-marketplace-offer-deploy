package handlers

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/operations"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	log "github.com/sirupsen/logrus"
)

type operationMessageHandler struct {
	operationFactory operation.Repository
}

func (h *operationMessageHandler) Handle(message *messaging.ExecuteInvokedOperation, context messaging.MessageHandlerContext) error {
	h.operationFactory.WithContext(context.Context())
	operation, err := h.operationFactory.First(message.OperationId)
	if err != nil {
		return err
	}
	return operation.Execute()
}

func NewOperationsMessageHandler(appConfig *config.AppConfig) *operationMessageHandler {
	handler := &operationMessageHandler{}
	service, err := newOperationService(appConfig)
	if err != nil {
		log.Errorf("Error creating operations message handler: %s", err)
		return nil
	}

	operationFactory, err := operation.NewRepository(service, &operations.OperationFuncProvider{})
	if err != nil {
		log.Errorf("Error creating operations message handler: %s", err)
		return nil
	}

	handler.operationFactory = operationFactory

	return handler
}

func newOperationService(appConfig *config.AppConfig) (*operation.OperationService, error) {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	sender, err := messaging.NewServiceBusMessageSender(credential, messaging.MessageSenderOptions{
		SubscriptionId:          appConfig.Azure.SubscriptionId,
		Location:                appConfig.Azure.Location,
		ResourceGroupName:       appConfig.Azure.ResourceGroupName,
		FullyQualifiedNamespace: appConfig.Azure.GetFullQualifiedNamespace(),
	})

	if err != nil {
		return nil, err
	}

	service, err := operation.NewService(db, sender, hook.Notify)
	if err != nil {
		return nil, err
	}
	return service, nil
}
