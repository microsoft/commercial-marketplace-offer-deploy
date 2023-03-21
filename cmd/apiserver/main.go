package main

import (
	"log"
	"net/http"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/routes"
)

const configurationFilePath string = "."

var (
	configuration *config.Configuration
)

func main() {
	log.Printf("Server started")

	loadConfiguration()

	router := routes.NewRouter(configuration)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func loadConfiguration() {
	var err error
	configuration, err = config.LoadConfiguration(configurationFilePath)
	if err != nil {
		log.Fatal()
	}
}
