package app

import (
	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/middleware"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/routes"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
)

func BuildApp(configurationFilePath string) *hosting.App {
	builder := hosting.NewAppBuilder()

	appConfig := config.AppConfig{}
	hosting.LoadConfiguration(configurationFilePath, nil, appConfig)
	builder.AddConfig(appConfig)

	builder.AddRoutes(func(options *hosting.RouteOptions) {
		databaseOptions := appConfig.GetDatabaseOptions()
		routes := routes.GetRoutes(databaseOptions)

		*options.Routes = routes
	})

	app := builder.Build(func(e *echo.Echo) {
		e.Group("/eventgrid", middleware.EventGridWebHookSubscriptionValidation())
	})
	return app
}
