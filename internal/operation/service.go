package operation

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type operationService struct {
	ctx    context.Context
	db     *gorm.DB
	notify hook.NotifyFunc
	id     uuid.UUID
	sender messaging.MessageSender
	log    *log.Entry
	// the reference of the operation
	operation *Operation
}

func (service *operationService) Context() context.Context {
	return service.ctx
}

func (service *operationService) saveChanges(notify bool) error {
	tx := service.db.WithContext(service.ctx).Begin()

	// could be an issue with starting the tx
	if tx.Error != nil {
		service.log.Errorf("saveChanges transaction aborted: %v", tx.Error)
		return tx.Error
	}

	tx.Save(service.operation.InvokedOperation)

	if tx.Error != nil {
		tx.Rollback()
		service.log.Errorf("saveChanges failed to save: %v", tx.Error)
		return tx.Error
	}

	tx.Commit()

	if tx.Error != nil {
		service.log.Errorf("saveChanges failed to commit transaction: %v", tx.Error)
		tx.Rollback()
		return tx.Error
	}

	if notify {
		snapshot := service.operation.InvokedOperation
		go func() {
			err := service.notification(snapshot) // if the notification fails, save still happened
			if err != nil {
				service.log.Errorf("notify failed: %v", err)
			}
		}()
	}

	return nil
}

// triggers a notification of the invoked operation's state
func (service *operationService) notification(snapshot model.InvokedOperation) error {
	message := service.getMessage(&snapshot)
	err := service.notify(service.Context(), message)
	if err != nil {
		service.log.Errorf("failed to add event message: %v", err)
		return err
	}
	return nil
}

func (service *operationService) dispatch() error {
	message := messaging.ExecuteInvokedOperation{OperationId: service.id}

	results, err := service.sender.Send(service.Context(), string(messaging.QueueNameOperations), message)
	if err != nil {
		return err
	}

	if len(results) == 1 && results[0].Error != nil {
		return results[0].Error
	}
	return nil
}

func (service *operationService) first() (*model.InvokedOperation, error) {
	record := &model.InvokedOperation{}
	result := service.db.Preload(clause.Associations).First(record, service.id)

	if result.Error != nil || result.RowsAffected == 0 {
		err := result.Error
		if err == nil {
			err = fmt.Errorf("failed to get invoked operation: %v", result.Error)
		}
		return nil, err
	}
	return record, nil
}

func (service *operationService) new(i *model.InvokedOperation) (*model.InvokedOperation, error) {
	tx := service.db.Begin()
	result := tx.Create(i)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return i, nil
}

// initializes and returns the single instance of InvokedOperation by the context's id
// if the id is invalid and an instance cannot be found, returns an error
func (service *operationService) initialize(id uuid.UUID) (*Operation, error) {
	service.id = id
	service.log = service.log.WithFields(log.Fields{
		"invokedOperationId": id,
	})

	invokedOperation, err := service.first()
	if err != nil {
		return nil, err
	}

	service.log = service.log.WithFields(log.Fields{
		"deploymentId": invokedOperation.DeploymentId,
	})

	service.operation = &Operation{
		InvokedOperation: *invokedOperation,
		service:          service,
	}

	return service.operation, nil
}

func (service *operationService) withContext(ctx context.Context) {
	service.ctx = ctx
	service.log = service.log.WithContext(ctx)
}

// encapsulates the conversion of an invoked operation to an event hook message
func (service *operationService) getMessage(io *model.InvokedOperation) *sdk.EventHookMessage {
	return mapToMessage(io)
}

func (service *operationService) deployment() *model.Deployment {
	deployment := &model.Deployment{}
	result := service.db.First(deployment, service.operation.DeploymentId)

	if result.RowsAffected == 0 {
		return nil
	}
	return deployment
}

// constructor factory of operation service
func newOperationService(appConfig *config.AppConfig) (*operationService, error) {
	credential, err := getAzureCredential()
	if err != nil {
		log.Errorf("Error creating Azure credential for hook.Queue: %v", err)
		return nil, err
	}

	sender, err := messaging.NewServiceBusMessageSender(credential, messaging.MessageSenderOptions{
		SubscriptionId:          appConfig.Azure.SubscriptionId,
		Location:                appConfig.Azure.Location,
		ResourceGroupName:       appConfig.Azure.ResourceGroupName,
		FullyQualifiedNamespace: appConfig.Azure.GetFullQualifiedNamespace(),
	})

	if err != nil {
		log.Errorf("Error creating message sender: %v", err)
		return nil, err
	}

	ctx := context.Background()

	return &operationService{
		ctx:    ctx,
		db:     data.NewDatabase(appConfig.GetDatabaseOptions()).Instance(),
		sender: sender,
		notify: hook.Notify,
		log:    log.WithContext(ctx),
	}, nil
}

func getAzureCredential() (azcore.TokenCredential, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	return credential, err
}
