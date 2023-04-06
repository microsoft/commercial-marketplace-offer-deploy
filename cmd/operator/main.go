package main

import (
	"log"
	"strconv"
	operator "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/app"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/messaging"
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

	appConfig := config.GetAppConfig()
	options := appConfig.GetDatabaseOptions()

	namespace := "bobjacmodm.servicebus.windows.net"
	operationsQueue := "deployoperationsqueue"
	sbConfig := messaging.ServiceBusConfig{
		Namespace: namespace,
		QueueName: operationsQueue,
	}

	db := data.NewDatabase(options)

	handler := messaging.NewOperationsHandler(db)
	receiver, err := messaging.NewServiceBusReceiver(sbConfig, handler)
	if err != nil {
		log.Fatal(err)
	}

	go receiver.Start()
	select{}
}
