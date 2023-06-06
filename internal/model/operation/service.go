package operation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/mapper"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/notification"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OperationService struct {
	ctx    context.Context
	db     *gorm.DB
	notify hook.NotifyFunc
	id     uuid.UUID
	sender messaging.MessageSender
	log    *log.Entry
	// the reference of the operation
	operation     *Operation
	stageNotifier notification.StageNotifier
}

func (service *OperationService) Context() context.Context {
	return service.ctx
}

func (service *OperationService) saveChanges(notify bool) error {
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

func (service *OperationService) notifyForStages() error {
	op := service.operation.InvokedOperation

	// only handle scheduled deployment operations, where we want to notify of stages getting scheduled
	if !op.IsRunning() && !op.IsFirstAttempt() {
		return errors.New("not a running deployment operation or not first attempt")
	}

	correlationId, err := op.CorrelationId()
	if err != nil {
		return err
	}

	notification := &model.StageNotification{
		OperationId:   op.ID,
		CorrelationId: *correlationId,
		Entries:       []model.StageNotificationEntry{},
		Done:          false,
	}

	deployment := service.deployment()
	if deployment == nil {
		return errors.New("deployment not found")
	}

	for _, stage := range deployment.Stages {
		notification.Entries = append(notification.Entries, model.StageNotificationEntry{
			StageId: stage.ID,
			Message: sdk.EventHookMessage{
				Id:     uuid.New(),
				Type:   string(sdk.EventTypeStageStarted),
				Status: sdk.StatusRunning.String(),
				Data: sdk.DeploymentEventData{
					DeploymentId:  int(deployment.ID),
					StageId:       &stage.ID,
					OperationId:   op.ID,
					CorrelationId: correlationId,
					Attempts:      1,
					StartedAt:     time.Now().UTC(),
				},
			},
		})
	}

	err = service.stageNotifier.Notify(service.ctx, notification)
	return err
}

// triggers a notification of the invoked operation's state
func (service *OperationService) publish(snapshot model.InvokedOperation) (uuid.UUID, error) {
	message := service.getMessage(&snapshot)
	notificationId, err := service.notify(service.Context(), message)
	if err != nil {
		service.log.Errorf("failed to publish event notification: %v", err)
		return notificationId, err
	}

	return notificationId, nil
}

func (service *OperationService) dispatch() error {
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

func (service *OperationService) first() (*model.InvokedOperation, error) {
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

func (service *OperationService) new(i *model.InvokedOperation) (*model.InvokedOperation, error) {
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
func (service *OperationService) initialize(id uuid.UUID) (*Operation, error) {
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

func (service *OperationService) withContext(ctx context.Context) {
	service.ctx = ctx
	service.log = service.log.WithContext(ctx)
}

// encapsulates the conversion of an invoked operation to an event hook message
func (service *OperationService) getMessage(io *model.InvokedOperation) *sdk.EventHookMessage {
	return mapper.MapInvokedOperation(io)
}

func (service *OperationService) deployment() *model.Deployment {
	deployment := &model.Deployment{}
	result := service.db.Preload("Stages").First(deployment, service.operation.DeploymentId)

	if result.RowsAffected == 0 {
		return nil
	}
	return deployment
}

// constructor factory of operation service
func NewService(db *gorm.DB, sender messaging.MessageSender, notify hook.NotifyFunc, stageNotifier notification.StageNotifier) (*OperationService, error) {

	ctx := context.Background()

	return &OperationService{
		ctx:           ctx,
		db:            db,
		sender:        sender,
		notify:        notify,
		log:           log.WithContext(ctx),
		stageNotifier: stageNotifier,
	}, nil
}
