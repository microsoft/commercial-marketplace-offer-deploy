package hosting

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/diagnostics"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	logger "github.com/microsoft/commercial-marketplace-offer-deploy/internal/log"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/tasks"
)

type AppBuilder struct {
	app *App
}

type ConfigureRoutesFunc func(options *RouteOptions)
type ConfigureAppConfigFunc func(config any)
type ConfigureEchoFunc func(e *echo.Echo)

func NewAppBuilder(name string) *AppBuilder {
	mutex.Lock()
	defer mutex.Unlock()

	if appInstance == nil {
		appInstance = &App{
			name:               name,
			server:             echo.New(),
			services:           []BackgroundService{},
			tasks:              []tasks.Task{},
			healthCheckResults: []diagnostics.HealthCheckResult{},
			healthCheckService: diagnostics.NewHealthCheckService(),
		}
	}

	builder := &AppBuilder{app: appInstance}
	return builder
}

func (b *AppBuilder) AddConfig(config *config.AppConfig) *AppBuilder {
	b.app.config = config
	return b
}

// adds readiness checks to the app instance when it starts. The order that the checks are added will be the order in which they are executed.
func (b *AppBuilder) AddReadinessCheck(check diagnostics.HealthCheck) *AppBuilder {
	b.app.healthCheckService.AddHealthCheck(check)
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
		log.Debug("registering route: { %s %s %s }", route.Name, route.Method, route.Path)
		router.Add(route.Method, route.Path, route.HandlerFunc)
	}

	return b
}

func (b *AppBuilder) Build(configure ConfigureEchoFunc) *App {
	//add middleware
	b.app.server.Use(middleware.Logger())
	b.app.server.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		var body any
		err := json.Unmarshal(reqBody, &body)
		if err != nil {
			return
		}
		log.WithField("request", body).Debug("Body")
	}))

	loggingConfig := b.app.config.GetLoggingOptions(b.app.name)
	logger.ConfigureLogging(loggingConfig)

	// configure event hook subsystem
	hook.Configure(b.app.config)

	if configure != nil {
		configure(b.app.server)
	}
	appInstance = b.app
	return appInstance
}
