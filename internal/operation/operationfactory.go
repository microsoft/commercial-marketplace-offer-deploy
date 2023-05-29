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
	service *operationService
}

func (iof *factory) Create(ctx context.Context, id uuid.UUID) (*Operation, error) {
	invokedOperation, err := iof.service.begin(ctx, id)

	if err != nil {
		return nil, err
	}

	instance := &Operation{
		InvokedOperation: *invokedOperation, //copy left
		service:          iof.service,
	}

	return instance, nil
}

func NewOperationFactory(appConfig *config.AppConfig) (Factory, error) {
	context, err := newOperationService(appConfig)
	if err != nil {
		return nil, err
	}

	factory := &factory{
		service: context,
	}

	return factory, nil
}
