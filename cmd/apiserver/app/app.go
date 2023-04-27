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

	builder.AddTask(newReadinessTask(appConfig, app.IsReady))
	builder.AddTask(newEventGridRegistrationTask(appConfig, app.IsReady))

	return app
}

func addReadinessChecks(builder *hosting.AppBuilder, appConfig *config.AppConfig) {
	defaultTimeout := time.Duration(2 * time.Minute)

	azureCredentialCheck := diagnostics.NewAzureCredentialHealthCheck(diagnostics.AzureCredentialHealthCheckOptions{
		Timeout: defaultTimeout,
	})

	azureRoleAssignmentsHealthCheck := diagnostics.NewRoleAssignmentsHealthCheck(diagnostics.AzureRoleAssignmentsHealthCheckOptions{
		SubscriptionId: appConfig.Azure.SubscriptionId,
		Timeout:        defaultTimeout,
	})

	builder.AddReadinessCheck(azureCredentialCheck)
	builder.AddReadinessCheck(azureRoleAssignmentsHealthCheck)
}
