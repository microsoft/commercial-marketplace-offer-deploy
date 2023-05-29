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

type operationService struct {
	ctx                 context.Context
	db                  *gorm.DB
	publishNotification hook.Notify
	id                  uuid.UUID
	sender              messaging.MessageSender
	log                 *log.Entry
	// the reference of the invokedOperation
	invokedOperation *model.InvokedOperation
}

func (ioc *operationService) Context() context.Context {
	return ioc.ctx
}

func (ioc *operationService) saveChanges() error {
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
func (ioc *operationService) notify() error {
	message := ioc.getMessage()

	err := ioc.publishNotification(ioc.Context(), message)
	if err != nil {
		ioc.log.Errorf("failed to add event message: %v", err)
		return err
	}
	return nil
}

func (ioc *operationService) retry() error {
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
func (ioc *operationService) begin(ctx context.Context, id uuid.UUID) (*model.InvokedOperation, error) {
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
func (ioc *operationService) getMessage() *sdk.EventHookMessage {
	return mapToMessage(ioc.invokedOperation)
}

func (ioc *operationService) deployment() *model.Deployment {
	deployment := &model.Deployment{}
	ioc.db.First(deployment, ioc.invokedOperation.DeploymentId)

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
		log.Errorf("Error creating message sender for hook.Queue: %v", err)
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
