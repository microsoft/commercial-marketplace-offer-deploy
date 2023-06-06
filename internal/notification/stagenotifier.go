package notification

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"gorm.io/gorm"
)

// dispatch notifications for stages
type StageNotifier interface {
	Notify(notification *model.StageNotification) error
}

// default implementation of the stage notifier
type stageNotifier struct {
	db *gorm.DB
}

func (s *stageNotifier) Notify(notification *model.StageNotification) error {
	return nil
}
