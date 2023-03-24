package runtime

import (
	"log"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/routes"
)

type App struct {
	config *config.Configuration
	e      *echo.Echo
}

type AppBuilder struct {
	app *App
}

var appInstance *App

// Gets the App instance running
func GetApp() *App {
	return appInstance
}

func (app *App) GetConfiguration() *config.Configuration {
	return app.config
}

type ConfigureRoutesFunc func(c *config.Configuration) *routes.Routes
type ConfigureConfigurationFunc func(config *config.Configuration)

func NewAppBuilder() *AppBuilder {
	builder := &AppBuilder{app: &App{e: echo.New()}}
	return builder
}

func (b *AppBuilder) AddConfig(configure ConfigureConfigurationFunc) *AppBuilder {
	b.app.config = &config.Configuration{}
	if configure != nil {
		configure(b.app.config)
	}
	return b
}

func (b *AppBuilder) AddRoutes(configure ConfigureRoutesFunc) *AppBuilder {
	routes := configure(b.app.config)

	for _, route := range *routes {
		log.Printf("registering route: %s", route)
		router := b.app.e.Router()
		router.Add(route.Method, route.Path, route.HandlerFunc)
	}
	return b
}

func (b *AppBuilder) Build() *App {
	//add middleware
	b.app.e.Use(middleware.Logger())

	appInstance = b.app
	return appInstance
}

func (app *App) Start(port int) error {
	address := ":" + strconv.Itoa(port)
	log.Printf("Server starting on %s", address)

	return app.e.Start(address)
}
