package main

import (
	"log"
	"strconv"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/app"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/routes"
)

const configurationFilePath string = "."

var (
	port int = 8080
)

func main() {
	formattedPort := ":" + strconv.Itoa(port)
	log.Printf("Server started on %s", formattedPort)

	builder := app.NewAppBuilder()
	builder.AddConfig(configureConfig)
	builder.AddRoutes(routes.GetRoutes)

	app := builder.Build()
	log.Fatal(app.Start(port))
}

// lint:ignore this value of c is never used (SA4006)
func configureConfig(c *config.Configuration) {
	var err error
	c, err = config.LoadConfiguration(configurationFilePath, nil)

	if err != nil {
		log.Fatal()
	}
}
