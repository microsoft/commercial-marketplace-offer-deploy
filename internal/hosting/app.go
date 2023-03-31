package hosting

import (
	"log"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type App struct {
	config any
	e      *echo.Echo
}

type AppBuilder struct {
	app *App
}

type RouteOptions struct {
	AppConfig any
	Routes    *Routes
}

var appInstance *App

// Gets the App instance running
func GetApp() *App {
	return appInstance
}

func (app *App) GetConfig() any {
	return app.config
}

type ConfigureRoutesFunc func(options *RouteOptions)
type ConfigureAppConfigFunc func(config *any)
type ConfigureEchoFunc func(e *echo.Echo)

func NewAppBuilder() *AppBuilder {
	builder := &AppBuilder{app: &App{e: echo.New()}}
	return builder
}

func (b *AppBuilder) AddConfig(configure ConfigureAppConfigFunc) *AppBuilder {
	b.app.config = new(any)
	if configure != nil {
		configure(&b.app.config)
	}
	return b
}

func (b *AppBuilder) AddRoutes(configure ConfigureRoutesFunc) *AppBuilder {
	router := b.app.e.Router()
	options := RouteOptions{Routes: &Routes{}, AppConfig: b.app.config}
	configure(&options)

	for _, route := range *options.Routes {
		log.Printf("registering route: { %s %s %s }", route.Name, route.Method, route.Path)
		router.Add(route.Method, route.Path, route.HandlerFunc)
	}
	return b
}

func (b *AppBuilder) Build() *App {
	//add middleware
	b.app.e.Use(middleware.Logger())
	b.app.e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		log.Printf("Request:\n %v", string(reqBody))
	}))

	appInstance = b.app
	return appInstance
}

// Start starts the server
// port: the port to listen on
// configure: (optional) a function to configure the echo server
func (app *App) Start(port int, configure ConfigureEchoFunc) error {
	address := ":" + strconv.Itoa(port)
	log.Printf("Server starting on %s", address)

	if configure != nil {
		configure(app.e)
	}
	return app.e.Start(address)
}
