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
			Name:        "EventGridWebHook",
			Method:      http.MethodPost,
			Path:        "/eventgrid",
			HandlerFunc: handlers.NewEventGridWebHookHandler(appConfig, hosting.GetAzureCredential()),
		},

		hosting.Route{
			Name:        "CreateDeployment",
			Method:      http.MethodPost,
			Path:        "/deployments",
			HandlerFunc: middleware.AddJwtBearer(handlers.NewCreateDeploymentHandler(appConfig, hosting.GetAzureCredential()), appConfig),
		},

		hosting.Route{
			Name:        "GetDeployment",
			Method:      http.MethodGet,
			Path:        "/deployments/:deploymentId",
			HandlerFunc: middleware.AddJwtBearer(handlers.GetDeployment, appConfig),
		},

		hosting.Route{
			Name:        "InvokeDeploymentOperation",
			Method:      http.MethodPost,
			Path:        "/deployments/:deploymentId/operation",
			HandlerFunc: middleware.AddJwtBearer(handlers.NewInvokeDeploymentOperationHandler(appConfig, hosting.GetAzureCredential()), appConfig),
		},

		hosting.Route{
			Name:        "ListDeployments",
			Method:      http.MethodGet,
			Path:        "/deployments",
			HandlerFunc: middleware.AddJwtBearer(handlers.NewListDeploymentsHandler(appConfig), appConfig),
		},

		hosting.Route{
			Name:        "UpdateDeployment",
			Method:      strings.ToUpper("Put"),
			Path:        "/deployments",
			HandlerFunc: middleware.AddJwtBearer(handlers.UpdateDeployment, appConfig),
		},

		hosting.Route{
			Name:        "CreatEventSubscription",
			Method:      http.MethodPost,
			Path:        "/events/subscriptions",
			HandlerFunc: middleware.AddJwtBearer(hosting.ToHandlerFunc(handlers.CreateEventHook, databaseOptions), appConfig),
		},

		hosting.Route{
			Name:        "DeleteEventSubscription",
			Method:      strings.ToUpper("Delete"),
			Path:        "/events/subscriptions/:subscriptionId",
			HandlerFunc: middleware.AddJwtBearer(handlers.DeleteEventSubscription, appConfig),
		},

		hosting.Route{
			Name:        "GetEventSubscription",
			Method:      http.MethodGet,
			Path:        "/events/subscriptions/:subscriptionId",
			HandlerFunc: middleware.AddJwtBearer(handlers.GetEventSubscription, appConfig),
		},

		hosting.Route{
			Name:        "ListEventSubscriptions",
			Method:      http.MethodGet,
			Path:        "/events/subscriptions",
			HandlerFunc: middleware.AddJwtBearer(handlers.ListEventSubscriptions, appConfig),
		},

		hosting.Route{
			Name:        "GetEvents",
			Method:      http.MethodGet,
			Path:        "/events",
			HandlerFunc: middleware.AddJwtBearer(handlers.GetEvents, appConfig),
		},

		hosting.Route{
			Name:        "GetDeploymentOperation",
			Method:      http.MethodGet,
			Path:        "/deployments/operations/:operationId",
			HandlerFunc: middleware.AddJwtBearer(handlers.GetDeploymentOperation, appConfig),
		},

		hosting.Route{
			Name:        "ListOperations",
			Method:      http.MethodGet,
			Path:        "/operations",
			HandlerFunc: middleware.AddJwtBearer(handlers.ListOperations, appConfig),
		},
	}

}
