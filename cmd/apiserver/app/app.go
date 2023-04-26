package app

import (
	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/middleware"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/routes"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
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

	app := builder.Build(func(e *echo.Echo) {
		e.Use(middleware.EventGridWebHookSubscriptionValidation())
	})

	builder.AddTask(newReadinessTask(appConfig, app.IsReady))
	builder.AddTask(newEventGridRegistrationTask(appConfig, app.IsReady))

	return app
}
