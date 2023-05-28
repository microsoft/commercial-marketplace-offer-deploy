package operation

import (
	"context"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
)

// Operation factory
type Factory interface {
	Create(ctx context.Context, id uuid.UUID) (*Operation, error)
}

type factory struct {
	context *operationService
}

func (iof *factory) Create(ctx context.Context, id uuid.UUID) (*Operation, error) {
	invokedOperation, err := iof.context.begin(ctx, id)

	if err != nil {
		return nil, err
	}

	instance := &Operation{
		InvokedOperation: *invokedOperation, //copy left
		context:          iof.context,
	}

	return instance, nil
}

func NewOperationFactory(appConfig *config.AppConfig) (Factory, error) {
	context, err := newOperationService(appConfig)
	if err != nil {
		return nil, err
	}

	factory := &factory{
		context: context,
	}

	return factory, nil
}
