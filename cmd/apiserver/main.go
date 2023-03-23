package main

import (
	"log"
	"strconv"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/routes"
)

const configurationFilePath string = "."

var (
	port          int = 8080
	configuration *config.Configuration
)

func main() {
	formattedPort := ":" + strconv.Itoa(port)
	log.Printf("Server started on %s", formattedPort)

	loadConfiguration()

	app := NewApp()

	routes := routes.GetRoutes()
	app.AddRoutes(&routes)
	app.Start(port)
}

func loadConfiguration() {
	var err error
	configuration, err = config.LoadConfiguration(configurationFilePath, nil)
	if err != nil {
		log.Fatal()
	}
}
