package handlers

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	e "github.com/microsoft/commercial-marketplace-offer-deploy/internal/events"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"gorm.io/gorm"
)

type eventsMessageHandler struct {
	publisher e.EventHookPublisher
}

func (h *eventsMessageHandler) Handle(message *events.EventHookMessage, context messaging.MessageHandlerContext) error {
	err := h.publisher.Publish(message)
	return err
}

//region factory

func NewEventsMessageHandler(appConfig *config.AppConfig) *eventsMessageHandler {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
	publisher := newWebHookPublisher(db)
	return &eventsMessageHandler{
		publisher: publisher,
	}
}

func newWebHookPublisher(db *gorm.DB) e.EventHookPublisher {
	subscriptionsProvider := e.NewEventHooksProvider(db)
	publisher := e.NewWebHookPublisher(subscriptionsProvider)
	return publisher
}

//endregion factory
