package main

import (
	"log"
	"strconv"

	operator "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/app"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/messaging/receivers"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
)

var (
	configurationFilePath string = "."
	port                  int    = 8180
)

func main() {
	formattedPort := ":" + strconv.Itoa(port)
	log.Printf("Server starting on %s", formattedPort)

	app := operator.BuildApp(configurationFilePath)
	go app.Start(port, nil)

	events, operations := getMessageReceivers()
	go events.Start()
	go operations.Start()

	select {}
}

func getMessageReceivers() (messaging.MessageReceiver, messaging.MessageReceiver) {
	appConfig := config.GetAppConfig()
	credential := hosting.GetAzureCredential()

	eventsReceiver := receivers.NewEventsMessageReceiver(appConfig, credential)
	operationsReceiver := receivers.NewOperationsMessageReceiver(appConfig, credential)

	return eventsReceiver, operationsReceiver
}
