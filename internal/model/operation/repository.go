package operation

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

// configure the operation
type Configure func(i *model.InvokedOperation) error

type RepositoryFactory func() (Repository, error)

type OperationFuncProvider interface {
	Get(operationType sdk.OperationType) (OperationTask, error)
}

// Operation factory
type Repository interface {
	New(operationType sdk.OperationType, configure Configure) (*Operation, error)
	First(id uuid.UUID) (*Operation, error)
	Any(id uuid.UUID) bool
	Provider(provider OperationFuncProvider) error
	WithContext(ctx context.Context)
}

type repository struct {
	manager  *OperationManager //clone managers from this instance. It acts as the base manager
	provider OperationFuncProvider
}

func (repo *repository) Provider(provider OperationFuncProvider) error {
	if provider == nil {
		return fmt.Errorf("provider cannot be nil")
	}
	repo.provider = provider
	return nil
}

func (repo *repository) WithContext(ctx context.Context) {
	repo.manager.withContext(ctx)
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

	_, err := repo.manager.new(instance)
	if err != nil {
		return nil, err
	}

	return repo.First(id)
}

// Gets the instance of an operation by id, otherwise, nil and an error
func (repo *repository) First(id uuid.UUID) (*Operation, error) {
	manager := CloneManager(repo.manager)

	operation, err := manager.initialize(id)
	if err != nil {
		return nil, err
	}

	if repo.provider != nil {
		task, err := repo.provider.Get(sdk.OperationType(operation.Name))

		if err != nil {
			return nil, err
		}
		operation.task = task
	}

	return operation, nil
}

func (repo *repository) Any(id uuid.UUID) bool {
	return repo.manager.any(id)
}

// NewRepository creates a new operation repository
// appConfig: application configuration
// provider: operation function provider, optional if the operation is not going to be executed and you want to interact with the operation
func NewRepository(manager *OperationManager, provider OperationFuncProvider) (Repository, error) {
	repo := &repository{
		manager:  manager,
		provider: provider,
	}

	return repo, nil
}

func NewRepositoryFactory(appConfig *config.AppConfig) RepositoryFactory {
	return func() (Repository, error) {
		db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()

		credential, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return nil, err
		}

		sender, err := messaging.NewServiceBusMessageSender(credential, messaging.MessageSenderOptions{
			SubscriptionId:          appConfig.Azure.SubscriptionId,
			Location:                appConfig.Azure.Location,
			ResourceGroupName:       appConfig.Azure.ResourceGroupName,
			FullyQualifiedNamespace: appConfig.Azure.GetFullQualifiedNamespace(),
		})

		if err != nil {
			return nil, err
		}

		service, err := NewManager(db, sender, hook.Notify)
		if err != nil {
			return nil, err
		}

		repository, err := NewRepository(service, nil)
		if err != nil {
			return nil, err
		}
		return repository, nil
	}
}
