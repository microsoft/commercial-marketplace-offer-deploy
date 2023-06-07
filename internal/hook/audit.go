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
			Hash:    message.GetHash(),
		})
	}()
}

func (r *EventHookAudit) IsDuplicate(message *sdk.EventHookMessage) bool {
    var count int64 = 0
    var countPtr *int64 = &count
    r.db.Model(&model.EventHookAuditEntry{}).Where("hash = ?", message.GetHash()).Count(countPtr)
    return countPtr != nil && *countPtr > 0
}


func NewAudit(db *gorm.DB) *EventHookAudit {
	return &EventHookAudit{db: db}
}


