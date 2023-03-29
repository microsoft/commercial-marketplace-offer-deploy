package routes

import (
	"net/http"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/handlers"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
)

func GetRoutes(databaseOptions *data.DatabaseOptions) hosting.Routes {
	return hosting.Routes{
		hosting.Route{
			Name:        "CreateDeployment",
			Method:      http.MethodPost,
			Path:        "/deployments",
			HandlerFunc: hosting.ToHandlerFunc(handlers.EventGridWebHook, databaseOptions),
		},
	}
}
