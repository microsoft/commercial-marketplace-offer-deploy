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
	return iof.service.init(ctx, id)
}

func NewOperationFactory(appConfig *config.AppConfig) (Factory, error) {
	service, err := newOperationService(appConfig)
	if err != nil {
		return nil, err
	}

	factory := &factory{
		service: service,
	}

	return factory, nil
}
