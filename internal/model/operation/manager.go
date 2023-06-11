package operation

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/mapper"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OperationManager struct {
	ctx    context.Context
	db     *gorm.DB
	notify hook.NotifyFunc
	id     uuid.UUID
	sender messaging.MessageSender
	log    *log.Entry
	// the reference of the operation
	operation *Operation
}

func (service *OperationManager) Context() context.Context {
	return service.ctx
}

func (service *OperationManager) saveChanges(notify bool) error {
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
		id, err := service.publish(snapshot) // if the notification fails, save still happened
		if err == nil {
			service.log.Infof("notification published [%v]", id)
			service.db.Save(service.operation.InvokedOperation)
		}
	}

	return nil
}

// triggers a notification of the invoked operation's state
func (service *OperationManager) publish(snapshot model.InvokedOperation) (uuid.UUID, error) {
	message := service.getMessage(&snapshot)
	if message == nil {
		service.log.Warnf("no message to publish for operation [%v]", snapshot.ID)
		return uuid.Nil, errors.New("no message to publish")
	}
	notificationId, err := service.notify(service.Context(), message)
	if err != nil {
		service.log.Errorf("failed to publish event notification: %v", err)
		return notificationId, err
	}

	return notificationId, nil
}

func (service *OperationManager) dispatch() error {
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

func (service *OperationManager) first() (*model.InvokedOperation, error) {
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

func (service *OperationManager) new(i *model.InvokedOperation) (*model.InvokedOperation, error) {
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
func (service *OperationManager) initialize(id uuid.UUID) (*Operation, error) {
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

	service.log.Trace("operation service initialized.")
	return service.operation, nil
}

func (service *OperationManager) withContext(ctx context.Context) {
	service.ctx = ctx
	service.log = service.log.WithContext(ctx)
}

// encapsulates the conversion of an invoked operation to an event hook message
func (service *OperationManager) getMessage(io *model.InvokedOperation) *sdk.EventHookMessage {
	if io.Status == string(sdk.StatusNone) {
		return nil
	}
	return mapper.MapInvokedOperation(io)
}

func (service *OperationManager) deployment() *model.Deployment {
	deployment := &model.Deployment{}
	result := service.db.Preload("Stages").First(deployment, service.operation.DeploymentId)

	if result.RowsAffected == 0 {
		return nil
	}
	return deployment
}

// constructor factory of operation service
func NewManager(db *gorm.DB, sender messaging.MessageSender, notify hook.NotifyFunc) (*OperationManager, error) {

	ctx := context.Background()

	return &OperationManager{
		ctx:    ctx,
		db:     db,
		sender: sender,
		notify: notify,
		log:    log.WithContext(ctx),
	}, nil
}

func CloneManager(manager *OperationManager) *OperationManager {
	return &OperationManager{
		ctx:    manager.ctx,
		db:     manager.db,
		sender: manager.sender,
		notify: manager.notify,
		log:    log.WithContext(manager.ctx),
	}
}
