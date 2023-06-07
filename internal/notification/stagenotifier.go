package notification

import (
	"context"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StageNotifyFunc func(ctx context.Context, notification *model.StageNotification) error

// dispatch notifications for stages
type StageNotifier interface {
	Notify(ctx context.Context, notification *model.StageNotification) error
}

// default implementation of the stage notifier
type stageNotifier struct {
	db *gorm.DB
}

func NewStageNotifier(db *gorm.DB) StageNotifier {
	return &stageNotifier{db: db}
}

func (s *stageNotifier) Notify(ctx context.Context, notification *model.StageNotification) error {
	return trace(s.notify)(ctx, notification)
}

func (s *stageNotifier) notify(ctx context.Context, notification *model.StageNotification) error {
	tx := s.db.WithContext(ctx).Begin()

	tx.Create(notification)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()

	return nil
}

func trace(notify StageNotifyFunc) StageNotifyFunc {
	return func(ctx context.Context, notification *model.StageNotification) error {
		logger := log.WithFields(
			log.Fields{
				"operationId":   notification.OperationId.String(),
				"correlationId": notification.CorrelationId.String(),
			})
		logger.Trace("Executing operation")
		err := notify(ctx, notification)
		logger.Trace("Execution of operation done")

		if err != nil {
			logger.WithError(err).Error("stage notification failed")
		}
		return err
	}
}
