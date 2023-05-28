package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/operations"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
)

type operationMessageHandler struct {
	operationFactory operation.Factory
	executorFactory  operation.ExecutorFactory
}

func (h *operationMessageHandler) Handle(message *messaging.ExecuteInvokedOperation, context messaging.MessageHandlerContext) error {
	executionContext, err := h.createContext(context.Context(), message.OperationId)
	if err != nil {
		return err
	}

	executor, err := h.getExecutor(executionContext)
	if err != nil {
		return err
	}
	return executor.Execute(executionContext)
}

func (h *operationMessageHandler) createContext(ctx context.Context, id uuid.UUID) (*operation.ExecutionContext, error) {
	invokedOperation, err := h.operationFactory.Create(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return operation.NewExecutionContext(ctx, invokedOperation), nil
}

func (h *operationMessageHandler) getExecutor(context *operation.ExecutionContext) (operation.Executor, error) {
	operationType := sdk.OperationType(context.InvokedOperation().Name)
	executor, err := h.executorFactory.Create(operationType)

	if err != nil {
		log.Errorf("Error creating executor for operation '%s': %s", operationType, err)
		return nil, err
	}

	return executor, nil
}

func NewOperationsMessageHandler(appConfig *config.AppConfig) *operationMessageHandler {
	invokedOperationFactory, err := operation.NewOperationFactory(appConfig)
	if err != nil {
		log.Errorf("Error creating operations message handler: %s", err)
		return nil
	}

	return &operationMessageHandler{
		executorFactory:  operations.NewExecutorFactory(appConfig),
		operationFactory: invokedOperationFactory,
	}
}
