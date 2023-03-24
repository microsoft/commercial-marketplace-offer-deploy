package main

import (
	"log"
	"strconv"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/routes"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/runtime"
)

const configurationFilePath string = "."

var (
	port int = 8080
)

func main() {
	formattedPort := ":" + strconv.Itoa(port)
	log.Printf("Server started on %s", formattedPort)

	builder := runtime.NewAppBuilder()
	builder.AddConfig(configureConfig)

	routes := routes.GetRoutes()
	builder.AddRoutes(&routes)

	app := builder.Build()
	log.Fatal(app.Start(port))
}

func configureConfig(c *config.Configuration) {
	configuration, err := config.LoadConfiguration(configurationFilePath, nil)

	if err != nil {
		log.Fatal()
	}
	c = configuration
}
