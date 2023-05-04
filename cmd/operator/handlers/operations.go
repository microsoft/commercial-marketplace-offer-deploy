package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/operations"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operation"
	log "github.com/sirupsen/logrus"
)

type operationMessageHandler struct {
	database  data.Database
	evaluator *invokedOperationEvaluator
	factory   operations.ExecutorFactory
}

func (h *operationMessageHandler) Handle(message *messaging.ExecuteInvokedOperation, context messaging.MessageHandlerContext) error {
	id, err := uuid.Parse(message.OperationId)
	if err != nil {
		return err
	}

	invokedOperation, err := h.getInvokedOperation(id)
	if err != nil {
		return err
	}

	if !h.shouldExecute(invokedOperation) {
		return nil
	}

	err = h.executeOperation(context.Context(), invokedOperation)

	return err
}

func (h *operationMessageHandler) getInvokedOperation(operationId uuid.UUID) (*data.InvokedOperation, error) {
	db := h.database.Instance()
	invokedOperation := &data.InvokedOperation{}
	db.First(&invokedOperation, operationId)

	if db.Error != nil {
		log.Errorf("Error retrieving invoked operation: %s", db.Error)
		return nil, db.Error
	}

	log.WithField("invokedOperationId", operationId).Infof("Retrieved invoked operation")

	return invokedOperation, nil
}

func (h *operationMessageHandler) shouldExecute(invokedOperation *data.InvokedOperation) bool {
	reasons, ok := h.evaluator.IsExecutable(invokedOperation)
	if !ok {
		log.Infof("Operation '%s' is not executable", invokedOperation.Name)
		for _, reason := range reasons {
			log.Info(reason)
		}
	}
	return ok
}

func (h *operationMessageHandler) executeOperation(ctx context.Context, invokedOperation *data.InvokedOperation) error {
	executor, err := h.factory.Create(operation.OperationType(invokedOperation.Name))
	if err != nil {
		return err
	}

	exec := operations.Trace(executor.Execute)
	return exec(ctx, invokedOperation)
}

func NewOperationsMessageHandler(appConfig *config.AppConfig) *operationMessageHandler {
	return &operationMessageHandler{
		database:  data.NewDatabase(appConfig.GetDatabaseOptions()),
		evaluator: &invokedOperationEvaluator{},
		factory:   operations.NewExecutorFactory(appConfig),
	}
}

//region mediator

// mediator evaluation of the invoked operation
type invokedOperationEvaluator struct {
	invokedOperation *data.InvokedOperation
	reasons          []string
}

func (e *invokedOperationEvaluator) IsExecutable(invokedOperation *data.InvokedOperation) ([]string, bool) {
	e.reasons = []string{}
	e.invokedOperation = invokedOperation

	isRunning := e.isRunning()
	reachMaxRetries := e.reachedMaxRetries()
	return e.reasons, (!isRunning && !reachMaxRetries)
}

func (e *invokedOperationEvaluator) isRunning() bool {
	isRunning := operation.Status(e.invokedOperation.Status) == operation.StatusRunning
	if isRunning {
		e.reasons = append(e.reasons, "'%s' is already running [%s]", e.invokedOperation.Name, e.invokedOperation.ID.String())
	}
	return isRunning
}

func (e *invokedOperationEvaluator) reachedMaxRetries() bool {
	reachedMaxRetries := e.invokedOperation.Attempts >= e.invokedOperation.Retries
	if reachedMaxRetries {
		e.reasons = append(e.reasons, "'%s' reached max retries [%s]", e.invokedOperation.Name, e.invokedOperation.ID.String())
	}
	return reachedMaxRetries
}

//endregion mediator
