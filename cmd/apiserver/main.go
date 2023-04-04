package main

import (
	"log"
	"strconv"

	apiserver "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/app"
)

var (
	configurationFilePath string = "."
	port                  int    = 8080
)

func main() {
	formattedPort := ":" + strconv.Itoa(port)
	log.Printf("Server starting on %s", formattedPort)

	app := apiserver.BuildApp(configurationFilePath)
	log.Fatal(app.Start(port, nil))
}
