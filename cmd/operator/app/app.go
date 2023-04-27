package app

import (
	"time"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/receivers"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/diagnostics"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
)

func BuildApp(configurationFilePath string) *hosting.App {
	builder := hosting.NewAppBuilder("operator")

	appConfig := &config.AppConfig{}
	config.LoadConfiguration(configurationFilePath, nil, appConfig)
	builder.AddConfig(appConfig)

	addReadinessChecks(builder, appConfig)
	addMessageReceivers(builder, appConfig)

	app := builder.Build(nil)
	return app
}

func addReadinessChecks(builder *hosting.AppBuilder, appConfig *config.AppConfig) {
	defaultTimeout := time.Duration(3 * time.Minute)

	azureCredentialCheck := diagnostics.NewAzureCredentialHealthCheck(diagnostics.AzureCredentialHealthCheckOptions{
		Timeout: defaultTimeout,
	})

	serviceBusCheck := diagnostics.NewServiceBusHealthCheck(diagnostics.ServiceBusHealthCheckOptions{
		FullyQualifiedNamespace: appConfig.Azure.GetFullQualifiedNamespace(),
		QueueName:               diagnostics.HealthCheckQueueName,
		Timeout:                 defaultTimeout,
	})

	builder.AddReadinessCheck(azureCredentialCheck)
	builder.AddReadinessCheck(serviceBusCheck)
}

func addMessageReceivers(builder *hosting.AppBuilder, appConfig *config.AppConfig) {
	eventsReceiver, operationsReceiver := getMessageReceivers(appConfig)
	builder.AddService(eventsReceiver)
	builder.AddService(operationsReceiver)
}

func getMessageReceivers(appConfig *config.AppConfig) (messaging.MessageReceiver, messaging.MessageReceiver) {
	credential := hosting.GetAzureCredential()

	eventsReceiver := receivers.NewEventsMessageReceiver(appConfig, credential)
	operationsReceiver := receivers.NewOperationsMessageReceiver(appConfig, credential)

	return eventsReceiver, operationsReceiver
}
