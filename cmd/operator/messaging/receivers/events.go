package receivers

import (
	"log"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
)

type eventsMessageHandler struct {
}

func (h *eventsMessageHandler) Handle(message *events.WebHookEventMessage, context messaging.MessageHandlerContext) error {
	return nil
}

//region factory

func newEventsMessageHandler() *eventsMessageHandler {
	return nil
}

func NewEventsMessageReceiver(appConfig *config.AppConfig) messaging.MessageReceiver {
	options := getOptions(appConfig)
	handler := newEventsMessageHandler()
	receiver, err := messaging.NewServiceBusReceiver(handler, options)
	if err != nil {
		log.Fatal(err)
	}
	return receiver
}

func getOptions(appConfig *config.AppConfig) messaging.ServiceBusMessageReceiverOptions {
	queueName := string(messaging.QueueNameEvents)
	options := messaging.ServiceBusMessageReceiverOptions{
		MessageReceiverOptions: messaging.MessageReceiverOptions{QueueName: queueName},
		NamespaceName:          appConfig.ServiceBusNamespace,
	}
	return options
}

//endregion factory
