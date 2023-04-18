package main

import (
	"log"
	"strconv"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
	testharness "github.com/microsoft/commercial-marketplace-offer-deploy/tools/testharness/app"
)

var (
	configurationFilePath string = "."
	port                  int    = 8280
)

func main() {
	formattedPort := ":" + strconv.Itoa(port)
	log.Printf("Server starting on %s", formattedPort)

	app := testharness.BuildApp(configurationFilePath)

	options := &hosting.AppStartOptions{
		Port: &port,
	}
	log.Fatal(app.Start(options))
}
