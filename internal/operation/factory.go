package operation

import (
	"context"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

type OperationFuncProvider interface {
	Get(operationType sdk.OperationType) (OperationFunc, error)
}

// Operation factory
type Factory interface {
	Create(ctx context.Context, id uuid.UUID) (*Operation, error)
}

type factory struct {
	service  *operationService
	provider OperationFuncProvider
}

func (iof *factory) Create(ctx context.Context, id uuid.UUID) (*Operation, error) {
	operation, err := iof.service.init(ctx, id)
	if err != nil {
		return nil, err
	}
	do, err := iof.provider.Get(sdk.OperationType(operation.Name))
	if err != nil {
		return nil, err
	}
	operation.do = do

	return operation, nil
}

func NewOperationFactory(appConfig *config.AppConfig, provider OperationFuncProvider) (Factory, error) {
	service, err := newOperationService(appConfig)
	if err != nil {
		return nil, err
	}

	factory := &factory{
		service: service,
	}

	return factory, nil
}
