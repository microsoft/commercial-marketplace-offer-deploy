package hook

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"gorm.io/gorm"
)

type EventHookMessageRecorder struct {
	db *gorm.DB
}

func (r *EventHookMessageRecorder) Record(message *sdk.EventHookMessage) error {
	return r.db.Create(&model.EventHookMessage{
		Message: message,
	}).Error
}

func NewRecorder(db *gorm.DB) *EventHookMessageRecorder {
	return &EventHookMessageRecorder{db: db}
}
