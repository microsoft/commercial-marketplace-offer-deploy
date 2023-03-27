package routes

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/handlers"
)

type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc echo.HandlerFunc
}

type Routes []Route

func GetRoutes(c *config.Configuration) *Routes {
	configuration = c
	return &routes
}

var configuration *config.Configuration

var routes = Routes{
	Route{
		"Index",
		http.MethodGet,
		"/",
		handlers.Index,
	},

	Route{
		"CreateDeployment",
		http.MethodPost,
		"/deployments",
		handlers.ToHandlerFunc(handlers.CreateDeployment, configuration),
	},

	Route{
		"GetDeployment",
		http.MethodGet,
		"/deployments/:deploymentId",
		handlers.GetDeployment,
	},

	Route{
		"InvokeOperation",
		http.MethodPost,
		"/deployment/:deploymentId/operation",
		handlers.ToHandlerFunc(handlers.InvokeOperation, configuration),
	},

	Route{
		"ListDeployments",
		http.MethodGet,
		"/deployments",
		handlers.ListDeployments,
	},

	Route{
		"UpdateDeployment",
		strings.ToUpper("Put"),
		"/deployments",
		handlers.UpdateDeployment,
	},

	Route{
		"CreatEventSubscription",
		http.MethodPost,
		"/events/:eventType/subscriptions",
		handlers.CreatEventSubscription,
	},

	Route{
		"DeleteEventSubscription",
		strings.ToUpper("Delete"),
		"/events/subscriptions/:subscriptionId",
		handlers.DeleteEventSubscription,
	},

	Route{
		"GetEventSubscription",
		http.MethodGet,
		"/events/subscriptions/:subscriptionId",
		handlers.GetEventSubscription,
	},

	Route{
		"ListEventSubscriptions",
		http.MethodGet,
		"/events/:eventType/subscriptions",
		handlers.ListEventSubscriptions,
	},

	Route{
		"GetEvents",
		http.MethodGet,
		"/events",
		handlers.GetEvents,
	},

	Route{
		"GetDeploymentOperation",
		http.MethodGet,
		"/operations/:operationId",
		handlers.GetDeploymentOperation,
	},

	Route{
		"ListOperations",
		http.MethodGet,
		"/operations",
		handlers.ListOperations,
	},
}
