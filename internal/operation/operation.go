package operation

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
)

// a executable operation with an execution context
type OperationFunc func(context *ExecutionContext) error

// remarks: decorator+visitor of Invoked Operation
type Operation struct {
	model.InvokedOperation
	service *operationService
}

func (io *Operation) Running() error {
	err, running := io.InvokedOperation.Running()

	if err != nil {
		return err
	}
	if !running {
		return io.SaveChanges()
	}
	return nil
}

func (io *Operation) Failed() error {
	io.InvokedOperation.Failed()
	return io.SaveChanges()
}

func (io *Operation) Success() error {
	io.InvokedOperation.Success()
	return io.SaveChanges()
}

func (io *Operation) SaveChanges() error {
	err := io.service.saveChanges()

	if err != nil {
		return err
	}
	io.service.notify() // if the notification failes, save still happened

	return nil
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
