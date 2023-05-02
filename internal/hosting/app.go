package hosting

import (
	"context"
	"os"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/diagnostics"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/tasks"
)

type App struct {
	name               string
	config             *config.AppConfig
	server             *echo.Echo
	services           []BackgroundService
	tasks              []tasks.Task
	healthCheckService diagnostics.HealthCheckService
	healthCheckResults []diagnostics.HealthCheckResult
}

type AppStartOptions struct {
	Port               *int
	ConfigureWebServer ConfigureEchoFunc
	WebServer          bool
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
func GetAppConfig() *config.AppConfig {
	return GetApp().GetConfig()
}

// whether the app is ready
func (app *App) GetHealthCheckResults() []diagnostics.HealthCheckResult {
	return app.healthCheckResults
}

// whether the app is ready
func (app *App) IsReady() bool {
	if _, err := os.Stat(app.config.GetReadinessFilePath()); err == nil {
		return true
	}
	return false
}

// GetConfig gets the app configuration
func (app *App) GetConfig() *config.AppConfig {
	return app.config
}

func (app *App) Name() string {
	return app.name
}

// Start starts the server
// port: the port to listen on
// configure: (optional) a function to configure the echo server
func (app *App) Start(options *AppStartOptions) error {
	go app.startServer(options)

	app.checkReadiness()

	go app.startServices()
	go app.startTasks()
	select {}
}

func (app *App) checkReadiness() {
	ctx := context.Background()
	results := app.healthCheckService.CheckHealth(ctx)
	app.healthCheckResults = results

	log.Infof("App Health check results: %v", results)
}

func (app *App) startServer(options *AppStartOptions) {
	if options != nil && options.WebServer {
		port := 8080

		if options.Port == nil {
			port = *options.Port
		}
		address := ":" + strconv.Itoa(port)
		log.Debugf("Server starting on local port %s", address)
		log.Debugf("Public domain: %s", app.config.Http.DomainName)

		if options.ConfigureWebServer != nil {
			options.ConfigureWebServer(app.server)
		}

		go app.server.Start(address)
	}
}

func (app *App) startTasks() {
	if len(app.tasks) > 0 {
		runner := tasks.NewTaskRunner()
		for _, task := range app.tasks {
			runner.Add(task)
		}
		go runner.Start()
	}
}

func (app *App) startServices() {
	for _, service := range app.services {
		log.Debug("Starting service: %s", service.GetName())
		go service.Start()
	}
}
