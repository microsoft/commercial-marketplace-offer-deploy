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

	builder.AddTask(newEventGridRegistrationTask(appConfig, app.IsReady))

	return app
}

func addReadinessChecks(builder *hosting.AppBuilder, appConfig *config.AppConfig) {
	defaultTimeout := time.Duration(2 * time.Minute)

	publicUrlHealthCheck := diagnostics.NewUrlHealthCheck(diagnostics.UrlHealthCheckOptions{
		ReadinessFilePath: appConfig.GetReadinessFilePath(),
		Url:               appConfig.GetPublicBaseUrl(),
		Timeout:           defaultTimeout,
	})

	azureCredentialCheck := diagnostics.NewAzureCredentialHealthCheck(diagnostics.AzureCredentialHealthCheckOptions{
		Timeout: defaultTimeout,
	})

	builder.AddReadinessCheck(publicUrlHealthCheck)
	builder.AddReadinessCheck(azureCredentialCheck)
}
