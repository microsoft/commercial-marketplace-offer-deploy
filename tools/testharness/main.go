package main

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/microsoft/commercial-marketplace-offer-deploy/tools/testharness/app"
)

var (
	port int = 8280
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceQuote:             false,
		DisableQuote:           true,
		DisableLevelTruncation: true,
	})

	formattedPort := ":" + strconv.Itoa(port)
	log.Printf("Server starting on %s", formattedPort)

	e := echo.New()
	app.AddRoutes(e)

	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		log.WithFields(log.Fields{
			"method": c.Request().Method,
			"url":    c.Request().URL,
			"body":   string(reqBody),
		}).Debug("Received")
	}))

	if err := e.Start(formattedPort); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
