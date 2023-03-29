package main

import (
	"log"
	"strconv"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/routes"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
)

var (
	configurationFilePath string = "."
	port                  int    = 8180
)

func main() {
	formattedPort := ":" + strconv.Itoa(port)
	log.Printf("Server starting on %s", formattedPort)

	app := buildApp()
	log.Fatal(app.Start(port))
}

func buildApp() *hosting.App {
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

	app := builder.Build()
	return app
}
