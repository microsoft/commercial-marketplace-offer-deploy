package app

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/middleware"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/routes"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/diagnostics"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
)

func BuildApp(configurationFilePath string) *hosting.App {
	builder := hosting.NewAppBuilder("apiserver")

	appConfig := &config.AppConfig{}
	config.LoadConfiguration(configurationFilePath, nil, appConfig)
	builder.AddConfig(appConfig)

	builder.AddRoutes(func(options *hosting.RouteOptions) {
		routes := routes.GetRoutes(appConfig)

		*options.Routes = routes
	})

	addReadinessChecks(builder, appConfig)

	app := builder.Build(func(e *echo.Echo) {
		e.Use(middleware.EventGridWebHookSubscriptionValidation())
	})

	builder.AddTask(newEventGridRegistrationTask(appConfig))

	return app
}

func addReadinessChecks(builder *hosting.AppBuilder, appConfig *config.AppConfig) {
	defaultTimeout := time.Duration(5 * time.Minute)

	azureCredentialCheck := diagnostics.NewAzureCredentialHealthCheck(diagnostics.AzureCredentialHealthCheckOptions{
		Timeout: defaultTimeout,
	})

	// TODO: get queue name, pull from config -> env var
	serviceBusCheck := diagnostics.NewServiceBusHealthCheck(diagnostics.ServiceBusHealthCheckOptions{
		FullyQualifiedNamespace: appConfig.Azure.GetFullQualifiedNamespace(),
		QueueName:               diagnostics.HealthCheckQueueName,
		Timeout:                 defaultTimeout,
	})

	publicUrlHealthCheck := diagnostics.NewUrlHealthCheck(diagnostics.UrlHealthCheckOptions{
		ReadinessFilePath: appConfig.GetReadinessFilePath(),
		Url:               appConfig.GetPublicBaseUrl().String(),
		Timeout:           defaultTimeout,
	})

	// make sure that the last thing we do is verify public url health which will signal that MODM is ready
	// to begin receiving traffic. The ready file will be created which the operator is waiting for

	builder.AddReadinessCheck(azureCredentialCheck)
	builder.AddReadinessCheck(serviceBusCheck)
	builder.AddReadinessCheck(publicUrlHealthCheck)
}
