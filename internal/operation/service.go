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
)

type operationService struct {
	ctx                 context.Context
	db                  *gorm.DB
	publishNotification hook.Notify
	id                  uuid.UUID
	sender              messaging.MessageSender
	log                 *log.Entry
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

	if notify {
		service.notify() // if the notification failes, save still happened
	}

	if tx.Error != nil {
		service.log.Errorf("saveChanges failed to commit transaction: %v", tx.Error)
		tx.Rollback()
	}

	return tx.Error
}

// triggers a notification of the invoked operation's state
func (service *operationService) notify() error {
	message := service.getMessage()

	err := service.publishNotification(service.Context(), message)
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

func (service *operationService) first(id uuid.UUID) (*model.InvokedOperation, error) {
	record := &model.InvokedOperation{}
	result := service.db.First(record, service.id)

	if result.Error != nil || result.RowsAffected == 0 {
		err := result.Error
		if err == nil {
			err = fmt.Errorf("failed to get invoked operation: %v", result.Error)
		}
		return nil, err
	}
	return record, nil
}

// initializes the context and returns the single instance of InvokedOperation by the context's id
// if the id is invalid and an instance cannot be found, returns an error
func (service *operationService) init(ctx context.Context, id uuid.UUID) (*Operation, error) {
	service.id = id
	service.ctx = ctx
	service.log = log.WithFields(log.Fields{
		"invokedOperationId": id,
	})

	invokedOperation, err := service.first(id)
	if err != nil {
		return nil, err
	}

	service.operation = &Operation{
		InvokedOperation: *invokedOperation,
		service:          service,
	}

	service.log = log.WithFields(log.Fields{
		"deploymentId": invokedOperation.DeploymentId,
	})

	return service.operation, nil
}

// encapsulates the conversion of an invoked operation to an event hook message
func (service *operationService) getMessage() *sdk.EventHookMessage {
	return mapToMessage(&service.operation.InvokedOperation)
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

	return &operationService{
		db:                  data.NewDatabase(appConfig.GetDatabaseOptions()).Instance(),
		sender:              sender,
		publishNotification: hook.Add,
	}, nil
}

func getAzureCredential() (azcore.TokenCredential, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	return credential, err
}
