package routes

import (
	"net/http"
	"strings"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/handlers"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/middleware"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hosting"
)

func GetRoutes(appConfig *config.AppConfig) hosting.Routes {
	databaseOptions := appConfig.GetDatabaseOptions()

	return hosting.Routes{
		hosting.Route{
			Name:        "Index",
			Method:      http.MethodGet,
			Path:        "/",
			HandlerFunc: handlers.Index,
		},

		hosting.Route{
			Name:        "CreateDeployment",
			Method:      http.MethodPost,
			Path:        "/deployments",
			HandlerFunc: middleware.AddJwtBearer(hosting.ToHandlerFunc(handlers.CreateDeployment, databaseOptions), appConfig),
		},

		hosting.Route{
			Name:        "GetDeployment",
			Method:      http.MethodGet,
			Path:        "/deployments/:deploymentId",
			HandlerFunc: handlers.GetDeployment,
		},

		hosting.Route{
			Name:        "InvokeDeploymentOperation",
			Method:      http.MethodPost,
			Path:        "/deployments/:deploymentId/operation",
			HandlerFunc: handlers.NewInvokeDeploymentOperationHandler(appConfig, hosting.GetAzureCredential()),
		},

		hosting.Route{
			Name:        "ListDeployments",
			Method:      http.MethodGet,
			Path:        "/deployments",
			HandlerFunc: handlers.ListDeployments,
		},

		hosting.Route{
			Name:        "UpdateDeployment",
			Method:      strings.ToUpper("Put"),
			Path:        "/deployments",
			HandlerFunc: handlers.UpdateDeployment,
		},

		hosting.Route{
			Name:        "CreatEventSubscription",
			Method:      http.MethodPost,
			Path:        "/events/subscriptions",
			HandlerFunc: hosting.ToHandlerFunc(handlers.CreateEventSubscription, databaseOptions),
		},

		hosting.Route{
			Name:        "DeleteEventSubscription",
			Method:      strings.ToUpper("Delete"),
			Path:        "/events/subscriptions/:subscriptionId",
			HandlerFunc: handlers.DeleteEventSubscription,
		},

		hosting.Route{
			Name:        "GetEventSubscription",
			Method:      http.MethodGet,
			Path:        "/events/subscriptions/:subscriptionId",
			HandlerFunc: handlers.GetEventSubscription,
		},

		hosting.Route{
			Name:        "ListEventSubscriptions",
			Method:      http.MethodGet,
			Path:        "/events/subscriptions",
			HandlerFunc: handlers.ListEventSubscriptions,
		},

		hosting.Route{
			Name:        "GetEvents",
			Method:      http.MethodGet,
			Path:        "/events",
			HandlerFunc: handlers.GetEvents,
		},

		hosting.Route{
			Name:        "GetDeploymentOperation",
			Method:      http.MethodGet,
			Path:        "/deployments/operations/:operationId",
			HandlerFunc: handlers.GetDeploymentOperation,
		},

		hosting.Route{
			Name:        "ListOperations",
			Method:      http.MethodGet,
			Path:        "/operations",
			HandlerFunc: handlers.ListOperations,
		},
	}

}
