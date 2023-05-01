package handlers

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"gorm.io/gorm"
)

type eventsMessageHandler struct {
	publisher hook.Publisher
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

func newWebHookPublisher(db *gorm.DB) hook.Publisher {
	subscriptionsProvider := hook.NewEventHooksProvider(db)
	publisher := hook.NewEventHookPublisher(subscriptionsProvider)
	return publisher
}

//endregion factory
