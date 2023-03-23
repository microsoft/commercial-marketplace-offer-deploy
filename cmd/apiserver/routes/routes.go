package routes

import (
	"net/http"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/handlers"
)

func GetRoutes() Routes {
	return routes
}

var routes = Routes{
	Route{
		"GET",
		"Index",
		"/",
		Index,
	},

	Route{
		http.MethodPost,
		"CreateDeployment",
		"/deployments",
		handlers.ToHandlerFunc(handlers.CreateDeploymentHandler),
	},

	// Route{
	// 	"GetDeployment",
	// 	strings.ToUpper("Get"),
	// 	"/deployments/{deploymentId}",
	// 	handlers.GetDeployment,
	// },

	// Route{
	// 	"InvokeOperation",
	// 	strings.ToUpper("Post"),
	// 	"/deployment/{deploymentId}/operation",
	// 	handlers.WithDatabase(handlers.InvokeOperation),
	// },

	// Route{
	// 	"ListDeployments",
	// 	strings.ToUpper("Get"),
	// 	"/deployments",
	// 	handlers.ListDeployments,
	// },

	// Route{
	// 	"UpdateDeployment",
	// 	strings.ToUpper("Put"),
	// 	"/deployments",
	// 	handlers.UpdateDeployment,
	// },

	// Route{
	// 	"CreatEventSubscription",
	// 	strings.ToUpper("Post"),
	// 	"/events/{topic}/subscriptions",
	// 	handlers.CreatEventSubscription,
	// },

	// Route{
	// 	"DeleteEventSubscription",
	// 	strings.ToUpper("Delete"),
	// 	"/events/subscriptions/{subscriptionId}",
	// 	handlers.DeleteEventSubscription,
	// },

	// Route{
	// 	"GetEventSubscription",
	// 	strings.ToUpper("Get"),
	// 	"/events/subscriptions/{subscriptionId}",
	// 	handlers.GetEventSubscription,
	// },

	// Route{
	// 	"ListEventSubscriptions",
	// 	strings.ToUpper("Get"),
	// 	"/events/{topic}/subscriptions",
	// 	handlers.ListEventSubscriptions,
	// },

	// Route{
	// 	"GetEvents",
	// 	strings.ToUpper("Get"),
	// 	"/events",
	// 	handlers.GetEvents,
	// },

	// Route{
	// 	"GetDeploymentOperation",
	// 	strings.ToUpper("Get"),
	// 	"/operations/{operationId}",
	// 	handlers.GetDeploymentOperation,
	// },

	// Route{
	// 	"ListOperations",
	// 	strings.ToUpper("Get"),
	// 	"/operations",
	// 	handlers.ListOperations,
	// },
}
