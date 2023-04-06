package app

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
)

func GetRoutes(appConfig *config.AppConfig) hosting.Routes {

	return hosting.Routes{
		hosting.Route{
			Name:        "WebHookResponse",
			Method:      http.MethodGet,
			Path:        "/deploymentevent",
			HandlerFunc: GetDeployment,
		},
	}
}

func GetDeployment(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

func BuildApp(configurationFilePath string) *hosting.App {
	builder := hosting.NewAppBuilder()

	appConfig := &config.AppConfig{}
	hosting.LoadConfiguration(configurationFilePath, nil, appConfig)
	builder.AddConfig(appConfig)

	builder.AddRoutes(func(options *hosting.RouteOptions) {
		routes := GetRoutes(appConfig)
		*options.Routes = routes
	})

	app := builder.Build(nil)
	return app
}