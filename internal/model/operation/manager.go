package operation

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/mapper"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OperationManager struct {
	ctx       context.Context
	db        *gorm.DB
	notify    hook.NotifyFunc
	id        uuid.UUID
	scheduler Scheduler
	log       *log.Entry
	// the reference of the operation
	operation *Operation
}

func (manager *OperationManager) Context() context.Context {
	return manager.ctx
}

func (manager *OperationManager) saveChanges(notify bool) error {
	tx := manager.db.WithContext(manager.ctx).Begin()

	// could be an issue with starting the tx
	if tx.Error != nil {
		manager.log.Errorf("saveChanges transaction aborted: %v", tx.Error)
		return tx.Error
	}

	tx.Save(manager.operation.InvokedOperation)

	if tx.Error != nil {
		tx.Rollback()
		manager.log.Errorf("saveChanges failed to save: %v", tx.Error)
		return tx.Error
	}

	tx.Commit()

	if tx.Error != nil {
		manager.log.Errorf("saveChanges failed to commit transaction: %v", tx.Error)
		tx.Rollback()
		return tx.Error
	}

	if notify {
		snapshot := manager.operation.InvokedOperation
		id, err := manager.publish(snapshot) // if the notification fails, save still happened
		if err == nil {
			manager.log.Infof("notification published [%v]", id)
			manager.db.Save(manager.operation.InvokedOperation)
		}
	}

	return nil
}

// triggers a notification of the invoked operation's state
func (manager *OperationManager) publish(snapshot model.InvokedOperation) (uuid.UUID, error) {
	message := manager.getMessage(&snapshot)
	if message == nil {
		manager.log.Warnf("no message to publish for operation [%v]", snapshot.ID)
		return uuid.Nil, errors.New("no message to publish")
	}
	notificationId, err := manager.notify(manager.Context(), message)
	if err != nil {
		manager.log.Errorf("failed to publish event notification: %v", err)
		return notificationId, err
	}

	return notificationId, nil
}

func (manager *OperationManager) dispatch() error {
	return manager.scheduler.Schedule(manager.ctx, manager.id)
}

func (manager *OperationManager) first() (*model.InvokedOperation, error) {
	record := &model.InvokedOperation{}
	result := manager.db.Preload(clause.Associations).First(record, manager.id)

	if result.Error != nil || result.RowsAffected == 0 {
		err := result.Error
		if err == nil {
			err = fmt.Errorf("failed to get invoked operation: %v", result.Error)
		}
		return nil, err
	}
	return record, nil
}

func (manager *OperationManager) new(i *model.InvokedOperation) (*model.InvokedOperation, error) {

	if i.Results == nil {
		i.Results = make(map[uint]*model.InvokedOperationResult)
	}

	tx := manager.db.Begin()
	result := tx.Create(i)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return i, nil
}

func (manager *OperationManager) any(id uuid.UUID) bool {
	var count int64
	manager.db.Model(&model.InvokedOperation{}).Where("id = ?", id).Count(&count)
	return count > 0
}

// initializes and returns the single instance of InvokedOperation by the context's id
// if the id is invalid and an instance cannot be found, returns an error
func (manager *OperationManager) initialize(id uuid.UUID) (*Operation, error) {
	manager.id = id
	manager.log = manager.log.WithFields(log.Fields{
		"invokedOperationId": id,
	})

	invokedOperation, err := manager.first()
	if err != nil {
		return nil, err
	}

	manager.log = manager.log.WithFields(log.Fields{
		"deploymentId": invokedOperation.DeploymentId,
	})

	manager.operation = &Operation{
		InvokedOperation: *invokedOperation,
		manager:          manager,
	}

	return manager.operation, nil
}

func (manager *OperationManager) withContext(ctx context.Context) {
	manager.ctx = ctx
	manager.log = manager.log.WithContext(ctx)
}

// encapsulates the conversion of an invoked operation to an event hook message
func (manager *OperationManager) getMessage(io *model.InvokedOperation) *sdk.EventHookMessage {
	if io.Status == string(sdk.StatusNone) {
		return nil
	}
	return mapper.MapInvokedOperation(io)
}

func (manager *OperationManager) deployment() *model.Deployment {
	deployment := &model.Deployment{}
	result := manager.db.Preload("Stages").First(deployment, manager.operation.DeploymentId)

	if result.RowsAffected == 0 {
		return nil
	}
	return deployment
}

// constructor factory of operation service
func NewManager(db *gorm.DB, scheduler Scheduler, notify hook.NotifyFunc) (*OperationManager, error) {

	ctx := context.Background()

	return &OperationManager{
		ctx:       ctx,
		db:        db,
		scheduler: scheduler,
		notify:    notify,
		log:       log.WithContext(ctx),
	}, nil
}

func CloneManager(manager *OperationManager) *OperationManager {
	return &OperationManager{
		ctx:       manager.ctx,
		db:        manager.db,
		scheduler: manager.scheduler,
		notify:    manager.notify,
		log:       log.WithContext(manager.ctx),
	}
}
