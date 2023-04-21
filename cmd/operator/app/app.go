package app

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/receivers"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
)

func BuildApp(configurationFilePath string) *hosting.App {
	builder := hosting.NewAppBuilder()

	appConfig := &config.AppConfig{}
	config.LoadConfiguration(configurationFilePath, nil, appConfig)
	builder.AddConfig(appConfig)

	eventsReceiver, operationsReceiver := getMessageReceivers(appConfig)
	builder.AddService(eventsReceiver)
	builder.AddService(operationsReceiver)

	app := builder.Build(nil)
	return app
}

func getMessageReceivers(appConfig *config.AppConfig) (messaging.MessageReceiver, messaging.MessageReceiver) {
	credential := hosting.GetAzureCredential()

	eventsReceiver := receivers.NewEventsMessageReceiver(appConfig, credential)

	operationsReceiver := receivers.NewOperationsMessageReceiver(appConfig, credential)

	return eventsReceiver, operationsReceiver
}
