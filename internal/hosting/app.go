package hosting

import (
	"log"
	"strconv"
	"sync"
	"time"
	"github.com/labstack/echo"
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

type RouteOptions struct {
	AppConfig any
	Routes    *Routes
}

var mutex sync.Mutex
var appInstance *App
var serverStarted = make(chan bool)

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
	go app.startServer(options)
	go app.startServices()
	go app.startTasks()

	select {}
}

func (app *App) startServer(options *AppStartOptions) {
	if options != nil && options.WebServer {
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
		app.waitForReadiness()
	} else {
		serverStarted <- true
	}
}

// run until server is started so we know we can execute other tasks that depend on the server
func (app *App) waitForReadiness() {
	time.Sleep(1 * time.Second)
}

func (app *App) startTasks() {
	<-serverStarted

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
		log.Printf("Starting service: %s", service.GetName())
		go service.Start()
	}
}
