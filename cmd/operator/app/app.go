package app

import (
	"time"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/receivers"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/diagnostics"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/notification"
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
	// have the operator wait for the readiness of the app which will be controlled via file check by the apiserver's
	// trigger that creates the ready file
	readyCheck := diagnostics.NewFuncHealthCheck(diagnostics.FuncHealthCheckOptions{
		Timeout: time.Duration(10 * time.Minute),
		Ready: func() bool {
			return hosting.GetApp().IsReady()
		},
	})
	builder.AddReadinessCheck(readyCheck)
}

func addMessageReceivers(builder *hosting.AppBuilder, appConfig *config.AppConfig) {
	eventsReceiver, operationsReceiver := getMessageReceivers(appConfig)
	builder.AddService(eventsReceiver)
	builder.AddService(operationsReceiver)

	stageNotificationService := notification.NewStageNotificationService()
	builder.AddService(stageNotificationService)

}

func getMessageReceivers(appConfig *config.AppConfig) (messaging.MessageReceiver, messaging.MessageReceiver) {
	credential := hosting.GetAzureCredential()

	eventsReceiver := receivers.NewEventsMessageReceiver(appConfig, credential)
	operationsReceiver := receivers.NewOperationsMessageReceiver(appConfig, credential)

	return eventsReceiver, operationsReceiver
}
