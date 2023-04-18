package main

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	apiserver "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/app"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
	logger "github.com/microsoft/commercial-marketplace-offer-deploy/internal/log"
)

var (
	configurationFilePath string = "."
	myLogger = logger.NewLoggerPublisher()
)

func main() {
	app := apiserver.BuildApp(configurationFilePath)
	startOptions := &hosting.AppStartOptions{
		Port:      to.Ptr(8080),
		WebServer: true,
	}

	log.Fatal(app.Start(startOptions))
}
