package operation

import (
	"context"

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

type InvokedOperationContext struct {
	ctx                 context.Context
	db                  *gorm.DB
	publishNotification hook.Notify
	id                  uuid.UUID
	sender              messaging.MessageSender
	log                 *log.Entry
	// the reference of the invoked operation that is being tracked
	invokedOperation *InvokedOperation
}

func (ioc *InvokedOperationContext) Context() context.Context {
	return ioc.ctx
}

func (ioc *InvokedOperationContext) saveChanges() error {
	tx := ioc.db.WithContext(ioc.ctx).Begin()

	// could be an issue with starting the tx
	if tx.Error != nil {
		ioc.log.Errorf("saveChanges transaction aborted: %v", tx.Error)
		return tx.Error
	}

	tx.Save(ioc.invokedOperation)

	if tx.Error != nil {
		tx.Rollback()
		ioc.log.Errorf("saveChanges failed to save: %v", tx.Error)
		return tx.Error
	}

	tx.Commit()

	if tx.Error != nil {
		ioc.log.Errorf("saveChanges failed to commit transaction: %v", tx.Error)
		tx.Rollback()
	}

	return tx.Error
}

// triggers a notification of the invoked operation's state
func (ioc *InvokedOperationContext) notify() error {
	message, err := ioc.getMessage()
	if err != nil {
		ioc.log.Errorf("could not get event message: %v", err)
		return err
	}
	err = ioc.publishNotification(ioc.Context(), message)
	if err != nil {
		ioc.log.Errorf("failed to add event message: %v", err)
		return err
	}
	return nil
}

func (ioc *InvokedOperationContext) retry() error {
	message := messaging.ExecuteInvokedOperation{OperationId: ioc.id}

	results, err := ioc.sender.Send(ioc.Context(), string(messaging.QueueNameOperations), message)
	if err != nil {
		return err
	}

	if len(results) == 1 && results[0].Error != nil {
		return results[0].Error
	}
	return nil
}

// initializes the context and returns the single instance of InvokedOperation by the context's id
// if the id is invalid and an instance cannot be found, returns an error
func (ioc *InvokedOperationContext) begin(ctx context.Context, id uuid.UUID) (*model.InvokedOperation, error) {
	ioc.id = id
	ioc.ctx = ctx
	ioc.log = log.WithFields(log.Fields{
		"invokedOperationId": id,
	})

	instance := &model.InvokedOperation{}
	err := ioc.db.First(instance, ioc.id).Error

	if err != nil {
		ioc.log.Errorf("failed to get invoked operation: %v", err)
		return nil, err
	}

	ioc.log = log.WithFields(log.Fields{
		"deploymentId": instance.DeploymentId,
	})

	return instance, err
}

// encapsulates the conversion of an invoked operation to an event hook message
func (ioc *InvokedOperationContext) getMessage() (*sdk.EventHookMessage, error) {
	return mapToMessage(ioc.invokedOperation)
}

// constructor factory of InvokedOperationContext
func NewInvokedOperationContext(appConfig *config.AppConfig) (*InvokedOperationContext, error) {
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
		log.Errorf("Error creating message sender for hook.Queue: %v", err)
		return nil, err
	}

	return &InvokedOperationContext{
		db:                  data.NewDatabase(appConfig.GetDatabaseOptions()).Instance(),
		sender:              sender,
		publishNotification: hook.Add,
	}, nil
}

func getAzureCredential() (azcore.TokenCredential, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	return credential, err
}
