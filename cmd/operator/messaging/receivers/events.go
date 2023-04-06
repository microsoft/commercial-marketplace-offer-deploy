package receivers

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	e "github.com/microsoft/commercial-marketplace-offer-deploy/internal/events"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"gorm.io/gorm"
)

type eventsMessageHandler struct {
	publisher e.WebHookPublisher
}

func (h *eventsMessageHandler) Handle(message *events.WebHookEventMessage, context messaging.MessageHandlerContext) error {
	err := h.publisher.Publish(message)
	return err
}

//region factory

func newEventsMessageHandler(appConfig *config.AppConfig) *eventsMessageHandler {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
	publisher := newWebHookPublisher(db)
	return &eventsMessageHandler{
		publisher: publisher,
	}
}

func NewEventsMessageReceiver(appConfig *config.AppConfig, credential azcore.TokenCredential) messaging.MessageReceiver {
	options := getOptions(appConfig)

	handler := newEventsMessageHandler(appConfig)
	receiver, err := messaging.NewServiceBusReceiver(handler, credential, options)
	if err != nil {
		log.Fatal(err)
	}
	return receiver
}

func newWebHookPublisher(db *gorm.DB) e.WebHookPublisher {
	subscriptionsProvider := e.NewGormSubscriptionsProvider(db)
	publisher := e.NewWebHookPublisher(subscriptionsProvider)
	return publisher
}

func getOptions(appConfig *config.AppConfig) messaging.ServiceBusMessageReceiverOptions {
	queueName := string(messaging.QueueNameEvents)
	options := messaging.ServiceBusMessageReceiverOptions{
		MessageReceiverOptions:  messaging.MessageReceiverOptions{QueueName: queueName},
		FullyQualifiedNamespace: appConfig.ServiceBusNamespace,
	}
	return options
}

//endregion factory
