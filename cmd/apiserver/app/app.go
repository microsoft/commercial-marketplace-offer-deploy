package app

import (
	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/middleware"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/routes"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/tasks"
)

func BuildApp(configurationFilePath string) *hosting.App {
	builder := hosting.NewAppBuilder()

	appConfig := &config.AppConfig{}
	config.LoadConfiguration(configurationFilePath, nil, appConfig)
	builder.AddConfig(appConfig)

	builder.AddRoutes(func(options *hosting.RouteOptions) {
		routes := routes.GetRoutes(appConfig)

		*options.Routes = routes
	})

	app := builder.Build(func(e *echo.Echo) {
		e.Use(middleware.EventGridWebHookSubscriptionValidation())
	})

	return app
}

func StartApp(app *hosting.App, options *hosting.AppStartOptions) error {
	go app.Start(options)

	runner := tasks.NewTaskRunner()
	runner.Add(getEventGridRegistrationTask())

	go runner.Start()

	select {}
}

func getEventGridRegistrationTask() tasks.Task {
	appConfig := config.GetAppConfig()

	taskOptions := EventGridRegistrationTaskOptions{
		Credential:      hosting.GetAzureCredential(),
		ResourceGroupId: appConfig.Azure.GetResourceGroupId(),
		EndpointUrl:     appConfig.Http.GetBaseUrl(true) + "/eventgrid",
	}
	task := NewEventGridRegistrationTask(taskOptions)

	return task
}
