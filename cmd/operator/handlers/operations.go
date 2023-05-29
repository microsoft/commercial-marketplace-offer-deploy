package handlers

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/operations"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/operation"
	log "github.com/sirupsen/logrus"
)

type operationMessageHandler struct {
	operationFactory operation.Factory
}

func (h *operationMessageHandler) Handle(message *messaging.ExecuteInvokedOperation, context messaging.MessageHandlerContext) error {
	operation, err := h.operationFactory.Create(context.Context(), message.OperationId)
	if err != nil {
		return err
	}
	return operation.Execute()
}

func NewOperationsMessageHandler(appConfig *config.AppConfig) *operationMessageHandler {
	handler := &operationMessageHandler{}

	operationFactory, err := operation.NewOperationFactory(appConfig, &operations.OperationFuncProvider{})
	if err != nil {
		log.Errorf("Error creating operations message handler: %s", err)
		return nil
	}

	handler.operationFactory = operationFactory

	return &operationMessageHandler{}
}
