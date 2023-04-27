package main

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/tools/testharness/app"
)

var (
	port int = 8280
)

func main() {
	formattedPort := ":" + strconv.Itoa(port)
	log.Printf("Server starting on %s", formattedPort)

	e := echo.New()
	app.AddRoutes(e)

	if err := e.Start(formattedPort); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
