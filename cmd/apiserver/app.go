package main

import (
	"log"
	"strconv"

	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/routes"
)

type App struct {
	config *config.Configuration
	e      *echo.Echo
}

type ConfigureConfigurationFunc func(config *config.Configuration)

func NewApp() *App {
	app := &App{e: echo.New()}
	return app
}

func (app *App) AddConfig(configure ConfigureConfigurationFunc) *App {
	app.config = &config.Configuration{}
	if configure != nil {
		configure(app.config)
	}
	return app
}

func (app *App) AddRoutes(routes *routes.Routes) *App {
	for _, route := range *routes {
		log.Print(route)
		router := app.e.Router()
		router.Add(route.Method, route.Path, route.HandlerFunc)
	}
	return app
}

func (app *App) Start(port int) {
	address := ":" + strconv.Itoa(port)
	log.Printf("Server started on %s", address)
	app.e.Logger.Fatal(app.e.Start(address))
}
