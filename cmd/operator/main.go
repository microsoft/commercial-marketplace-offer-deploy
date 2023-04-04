package main

import (
	"log"
	"strconv"

	operator "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/app"
)

var (
	configurationFilePath string = "."
	port                  int    = 8180
)

func main() {
	formattedPort := ":" + strconv.Itoa(port)
	log.Printf("Server starting on %s", formattedPort)

	app := operator.BuildApp(configurationFilePath)
	log.Fatal(app.Start(port, nil))
}
