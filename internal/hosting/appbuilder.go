package hosting

import (
	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	logger "github.com/microsoft/commercial-marketplace-offer-deploy/internal/log"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/tasks"
)

type AppBuilder struct {
	app *App
}

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
			ready:    make(chan bool),
		}
	}

	builder := &AppBuilder{app: appInstance}
	return builder
}

func (b *AppBuilder) AddConfig(config *config.AppConfig) *AppBuilder {
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

	loggingConfig := b.app.config.GetLoggingOptions()
	logger.ConfigureLogging(loggingConfig)

	if configure != nil {
		configure(b.app.server)
	}
	appInstance = b.app
	return appInstance
}
