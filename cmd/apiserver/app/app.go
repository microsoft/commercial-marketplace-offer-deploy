package app

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/routes"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
)

func GetApp(configurationFilePath string) *hosting.App {
	builder := hosting.NewAppBuilder()

	builder.AddConfig(func(c *any) {
		appSettings := &config.AppSettings{}
		hosting.LoadConfiguration(configurationFilePath, nil, &appSettings)
		*c = *appSettings
	})

	builder.AddRoutes(func(options *hosting.RouteOptions) {
		appSettings := options.AppConfig.(config.AppSettings)
		databaseOptions := appSettings.GetDatabaseOptions()
		routes := routes.GetRoutes(databaseOptions)

		*options.Routes = routes
	})

	app := builder.Build(nil)
	return app
}
