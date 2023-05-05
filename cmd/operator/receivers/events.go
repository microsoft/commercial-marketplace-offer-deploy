package receivers

import (
	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/handlers"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
)

// creates the new events message receiver, wiring up the events message handler
func NewEventsMessageReceiver(appConfig *config.AppConfig, credential azcore.TokenCredential) messaging.MessageReceiver {
	options := getOptions(appConfig)
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
	}

	handler, err := handlers.NewEventsMessageHandler(appConfig, credential)
	if err != nil {
		log.Fatal(err)
	}

	receiver, err := messaging.NewServiceBusReceiver(handler, credential, options)
	if err != nil {
		log.Fatal(err)
	}
	return receiver
}

func getOptions(appConfig *config.AppConfig) messaging.ServiceBusMessageReceiverOptions {
	queueName := string(messaging.QueueNameEvents)
	options := messaging.ServiceBusMessageReceiverOptions{
		MessageReceiverOptions:  messaging.MessageReceiverOptions{QueueName: queueName},
		FullyQualifiedNamespace: appConfig.Azure.GetFullQualifiedNamespace(),
	}
	return options
}
