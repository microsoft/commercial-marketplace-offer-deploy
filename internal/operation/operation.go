package operation

import (
	"context"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
)

// a executable operation with an execution context
type OperationFunc func(context *ExecutionContext) error

// remarks: Invoked Operation decorator+visitor
type Operation struct {
	model.InvokedOperation
	service *operationService
	do      OperationFunc
}

func (io *Operation) Context() context.Context {
	return io.service.ctx
}

func (io *Operation) Running() error {
	err, running := io.InvokedOperation.Running()

	if err != nil {
		return err
	}
	if !running {
		return io.service.saveChanges(true)
	}
	return nil
}

func (io *Operation) Attribute(key model.AttributeKey, v any) error {
	io.InvokedOperation.Attribute(key, v)
	return io.saveChangesWithoutNotification()
}

func (io *Operation) Value(v any) error {
	io.InvokedOperation.Value(v)
	return io.saveChangesWithoutNotification()
}

func (io *Operation) Failed() error {
	io.InvokedOperation.Failed()
	return io.service.saveChanges(true)
}

func (io *Operation) Success() error {
	io.InvokedOperation.Success()
	return io.service.saveChanges(true)
}

func (io *Operation) SaveChanges() error {
	return io.saveChangesWithoutNotification()
}

// Attempts to trigger a retry of the operation, if the operation has a retriable state
func (io *Operation) Retry() error {
	if !io.InvokedOperation.IsRetriable() {
		return nil
	}

	err := io.service.retry()
	if err != nil {
		return err
	}
	return nil
}

// provides access to latest instance of associated deployment
func (io *Operation) Deployment() *model.Deployment {
	return io.service.deployment()
}

// executes the operation
func (io *Operation) Execute() error {
	context := newExecutionContext(io)
	executor := NewExecutor(io.do)

	return executor.Execute(context)
}

func (o *Operation) saveChangesWithoutNotification() error {
	return o.service.saveChanges(false)
}
