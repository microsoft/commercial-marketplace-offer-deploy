package events

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	model "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"gorm.io/gorm"
)

type SubscriptionsProvider interface {
	// Gets the subscriptions for an event type
	GetSubscriptions(eventType model.EventType) ([]*data.EventSubscription, error)
}

type gormSubscriptionsProvider struct {
	db *gorm.DB
}

func NewGormSubscriptionsProvider(db *gorm.DB) SubscriptionsProvider {
	provider := &gormSubscriptionsProvider{db}
	return provider
}

func (p gormSubscriptionsProvider) GetSubscriptions(eventType model.EventType) ([]*data.EventSubscription, error) {
	data := []*data.EventSubscription{}
	tx := p.db.Where("eventType = ?", eventType.String()).Find(&data)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return data, nil
}
