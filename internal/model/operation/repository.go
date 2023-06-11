package operation

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

// configure the operation
type Configure func(i *model.InvokedOperation) error

type OperationFuncProvider interface {
	Get(operationType sdk.OperationType) (OperationFunc, error)
}

// Operation factory
type Repository interface {
	New(operationType sdk.OperationType, configure Configure) (*Operation, error)
	First(id uuid.UUID) (*Operation, error)
	Provider(provider OperationFuncProvider) error
	WithContext(ctx context.Context)
}

type repository struct {
	managerTemplate *OperationManager //clone managers from this instance
	provider        OperationFuncProvider
}

func (repo *repository) Provider(provider OperationFuncProvider) error {
	if provider == nil {
		return fmt.Errorf("provider cannot be nil")
	}
	repo.provider = provider
	return nil
}

func (repo *repository) WithContext(ctx context.Context) {
	repo.managerTemplate.withContext(ctx)
}

// creates a new operation instance by type
func (repo *repository) New(operationType sdk.OperationType, configure Configure) (*Operation, error) {
	id := uuid.New()

	instance := &model.InvokedOperation{}
	if configure != nil {
		err := configure(instance)
		if err != nil {
			return nil, err
		}
	}

	// no matter what, these values will override anything set by the configure function
	instance.Name = operationType.String()
	instance.ID = id

	_, err := repo.managerTemplate.new(instance)
	if err != nil {
		return nil, err
	}

	return repo.First(id)
}

// Gets the instance of an operation by id, otherwise, nil and an error
func (repo *repository) First(id uuid.UUID) (*Operation, error) {
	manager := CloneManager(repo.managerTemplate)

	operation, err := manager.initialize(id)
	if err != nil {
		return nil, err
	}

	if repo.provider != nil {
		do, err := repo.provider.Get(sdk.OperationType(operation.Name))

		if err != nil {
			return nil, err
		}
		operation.do = do
	}

	return operation, nil
}

// NewRepository creates a new operation repository
// appConfig: application configuration
// provider: operation function provider, optional if the operation is not going to be executed and you want to interact with the operation
func NewRepository(manager *OperationManager, provider OperationFuncProvider) (Repository, error) {
	repo := &repository{
		managerTemplate: manager,
		provider:        provider,
	}

	return repo, nil
}
