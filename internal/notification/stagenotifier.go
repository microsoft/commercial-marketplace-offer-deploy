package notification

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"gorm.io/gorm"
)

type StageNotifier interface {
	Notify(notification *model.StageNotification) error
}

type stageNotifier struct {
	db *gorm.DB
}

func (s *stageNotifier) Notify(notification *model.StageNotification) error {
	return nil
}
