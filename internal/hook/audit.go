package hook

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"gorm.io/gorm"
)

type EventHookAudit struct {
	db *gorm.DB
}

func (r *EventHookAudit) Log(message *sdk.EventHookMessage) {
	go func() {
		r.db.Create(&model.EventHookAuditEntry{
			Message: message,
		})
	}()
}

func NewAudit(db *gorm.DB) *EventHookAudit {
	return &EventHookAudit{db: db}
}
