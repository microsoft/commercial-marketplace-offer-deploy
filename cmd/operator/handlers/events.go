package handlers

import (
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type eventsMessageHandler struct {
	publisher hook.Publisher
}

func (h *eventsMessageHandler) Handle(message *sdk.EventHookMessage, context messaging.MessageHandlerContext) error {
	bytes, _ := json.Marshal(message)
	log.WithField("eventHookMessage", string(bytes)).Trace("Events Handler excuting")

	err := h.publisher.Publish(message)
	return err
}

//region factory

func NewEventsMessageHandler(appConfig *config.AppConfig, credential azcore.TokenCredential) (*eventsMessageHandler, error) {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
	publisher := newWebHookPublisher(db)

	return &eventsMessageHandler{
		publisher: publisher,
	}, nil
}

func newWebHookPublisher(db *gorm.DB) hook.Publisher {
	subscriptionsProvider := hook.NewEventHooksProvider(db)
	recorder := hook.NewRecorder(db)
	publisher := hook.NewEventHookPublisher(subscriptionsProvider, recorder)
	return publisher
}

//endregion factory
