package hosting

import (
	"log"
	"strconv"
	"sync"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/tasks"
)

type App struct {
	config   any
	server   *echo.Echo
	services []BackgroundService
	tasks    []tasks.Task
}

type AppStartOptions struct {
	Port               *int
	ConfigureWebServer ConfigureEchoFunc
	WebServer          bool
}

type AppBuilder struct {
	app *App
}

type RouteOptions struct {
	AppConfig any
	Routes    *Routes
}

var mutex sync.Mutex
var appInstance *App

// Gets the App instance running
func GetApp() *App {
	return appInstance
}

// Gets strongly typed the App configuration
func GetAppConfig[T any]() T {
	return GetApp().GetConfig().(T)
}

// GetConfig gets the app configuration
func (app *App) GetConfig() any {
	return app.config
}

// Start starts the server
// port: the port to listen on
// configure: (optional) a function to configure the echo server
func (app *App) Start(options *AppStartOptions) error {
	if options != nil {
		if options.WebServer {
			port := 8080

			if options.Port == nil {
				port = *options.Port
			}
			address := ":" + strconv.Itoa(port)
			log.Printf("Server starting on local port %s", address)

			if options.ConfigureWebServer != nil {
				options.ConfigureWebServer(app.server)
			}
			go app.server.Start(address)
		}
	}

	for _, service := range app.services {
		log.Printf("Starting service: %s", service.GetName())
		go service.Start()
	}

	if len(app.tasks) > 0 {
		runner := tasks.NewTaskRunner()
		for _, task := range app.tasks {
			runner.Add(task)
		}
		go runner.Start()
	}

	select {}
}

// App Builder

type ConfigureRoutesFunc func(options *RouteOptions)
type ConfigureAppConfigFunc func(config any)
type ConfigureEchoFunc func(e *echo.Echo)

func NewAppBuilder() *AppBuilder {
	mutex.Lock()
	defer mutex.Unlock()

	if appInstance == nil {
		appInstance = &App{
			server:   echo.New(),
			services: []BackgroundService{},
		}
	}

	builder := &AppBuilder{app: appInstance}
	return builder
}

func (b *AppBuilder) AddConfig(config any) *AppBuilder {
	b.app.config = config
	return b
}

func (b *AppBuilder) AddService(service BackgroundService) *AppBuilder {
	b.app.services = append(b.app.services, service)
	return b
}

func (b *AppBuilder) AddTask(task tasks.Task) *AppBuilder {
	b.app.tasks = append(b.app.tasks, task)
	return b
}

func (b *AppBuilder) AddRoutes(configure ConfigureRoutesFunc) *AppBuilder {
	router := b.app.server.Router()
	options := RouteOptions{Routes: &Routes{}, AppConfig: b.app.config}
	configure(&options)

	for _, route := range *options.Routes {
		log.Printf("registering route: { %s %s %s }", route.Name, route.Method, route.Path)
		router.Add(route.Method, route.Path, route.HandlerFunc)
	}

	return b
}

func (b *AppBuilder) Build(configure ConfigureEchoFunc) *App {
	//add middleware
	b.app.server.Use(middleware.Logger())
	b.app.server.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		log.Printf("Body:\n %v\n", string(reqBody))
	}))

	if configure != nil {
		configure(b.app.server)
	}
	appInstance = b.app
	return appInstance
}
