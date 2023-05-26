package operation

import (
	"context"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
)

// remarks: decorator of InvokedOperation
type InvokedOperation struct {
	model.InvokedOperation
	context *InvokedOperationContext
}

func (io *InvokedOperation) Running() error {
	err := io.InvokedOperation.Running()

	if err != nil {
		return err
	}
	return io.SaveChanges()
}

func (io *InvokedOperation) SaveChanges() error {
	err := io.context.saveChanges()

	if err != nil {
		return err
	}
	io.context.notify() // if the notification failes, save still happened

	return nil
}

// retries the operation
func (io *InvokedOperation) Retry() error {
	err := io.context.retry()
	if err != nil {
		return err
	}
	io.context.notify()

	return nil
}

type InvokedOperationFactory struct {
	context *InvokedOperationContext
}

func (iof *InvokedOperationFactory) Create(ctx context.Context, id uuid.UUID) (*InvokedOperation, error) {
	invokedOperation, err := iof.context.begin(ctx, id)

	if err != nil {
		return nil, err
	}

	instance := &InvokedOperation{
		InvokedOperation: *invokedOperation, //copy left
		context:          iof.context,
	}

	return instance, nil
}

func NewInvokedOperationFactory(appConfig *config.AppConfig) (*InvokedOperationFactory, error) {
	context, err := NewInvokedOperationContext(appConfig)
	if err != nil {
		return nil, err
	}

	factory := &InvokedOperationFactory{
		context: context,
	}

	return factory, nil
}
