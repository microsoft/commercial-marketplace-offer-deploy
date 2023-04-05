package routes

import (
	"net/http"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/handlers"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
)

func GetRoutes(appConfig *config.AppConfig) hosting.Routes {
	eventGridWebHookHandler := handlers.NewEventGridWebHookHandler(appConfig, hosting.GetAzureCredential())

	return hosting.Routes{
		hosting.Route{
			Name:        "EventGridWebHook",
			Method:      http.MethodPost,
			Path:        "/eventgrid",
			HandlerFunc: eventGridWebHookHandler,
		},
	}
}
