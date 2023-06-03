package handlers

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/operations"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
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

	operationFactory, err := operation.NewRepository(appConfig, &operations.OperationFuncProvider{})
	if err != nil {
		log.Errorf("Error creating operations message handler: %s", err)
		return nil
	}

	handler.operationFactory = operationFactory

	return handler
}
