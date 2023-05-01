package hook

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"gorm.io/gorm"
)

type EventHooksProvider interface {
	// Gets the subscriptions for an event type
	Get() ([]*data.EventHook, error)
}

type provider struct {
	db *gorm.DB
}

func NewEventHooksProvider(db *gorm.DB) EventHooksProvider {
	provider := &provider{db}
	return provider
}

func (p provider) Get() ([]*data.EventHook, error) {
	data := []*data.EventHook{}
	tx := p.db.Find(&data)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return data, nil
}
