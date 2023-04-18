package main

import (
	operator "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/app"
)

var (
	configurationFilePath string = "."
)

func main() {
	app := operator.BuildApp(configurationFilePath)
	app.Start(nil)
}
