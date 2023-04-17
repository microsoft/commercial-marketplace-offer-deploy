package main

import (
	"strconv"

	apiserver "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/app"
	logger "github.com/microsoft/commercial-marketplace-offer-deploy/internal/log"
	"github.com/sirupsen/logrus"
)

var (
	configurationFilePath string = "."
	port                  int    = 8081
)

func main() {
	formattedPort := ":" + strconv.Itoa(port)

	myLogger := logger.NewLoggerPublisher()

	myLogger.Publish(&logger.LogMessage{
		Message: "apiserver: Server starting on " + formattedPort,
		Level:   logrus.InfoLevel,
	})

	app := apiserver.BuildApp(configurationFilePath)
	err := app.Start(port, nil)

	if err != nil {
		myLogger.Publish(&logger.LogMessage{
			Message: "apiserver: Server failed to start on port " + formattedPort,
			Level:   logrus.FatalLevel,
		})
	}
}
