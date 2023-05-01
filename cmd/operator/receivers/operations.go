package receivers

import (
	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/handlers"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
)

func NewOperationsMessageReceiver(appConfig *config.AppConfig, credential azcore.TokenCredential) messaging.MessageReceiver {
	handler := handlers.NewOperationsMessageHandler(appConfig)

	options := messaging.ServiceBusMessageReceiverOptions{
		MessageReceiverOptions:  messaging.MessageReceiverOptions{QueueName: string(messaging.QueueNameOperations)},
		FullyQualifiedNamespace: appConfig.Azure.GetFullQualifiedNamespace(),
	}
	receiver, err := messaging.NewServiceBusReceiver(handler, credential, options)
	if err != nil {
		log.Fatal(err)
	}
	return receiver
}
