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
		receivedBody := string(reqBody)
		if len(receivedBody) > 0 {
			log.WithFields(log.Fields{
				"method": c.Request().Method,
				"url":    c.Request().URL,
				"body":   string(receivedBody),
			}).Print("Received")
		}

		responseBody := string(resBody)
		if len(responseBody) > 0 {
			log.WithFields(log.Fields{
				"body": string(responseBody),
			}).Print("Returned")
		}
	}))

	if err := e.Start(formattedPort); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
