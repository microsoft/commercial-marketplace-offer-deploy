package routes

import (
	"strings"

	. "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/handlers"
)

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"CreateDeployment",
		strings.ToUpper("Post"),
		"/deployments",
		CreateDeployment,
	},

	Route{
		"GetDeployment",
		strings.ToUpper("Get"),
		"/deployments/{deploymentId}",
		GetDeployment,
	},

	Route{
		"InvokeOperation",
		strings.ToUpper("Post"),
		"/deployment/{deploymentId}/operation",
		InvokeOperation,
	},

	Route{
		"ListDeployments",
		strings.ToUpper("Get"),
		"/deployments",
		ListDeployments,
	},

	Route{
		"UpdateDeployment",
		strings.ToUpper("Put"),
		"/deployments",
		UpdateDeployment,
	},

	Route{
		"CreatEventSubscription",
		strings.ToUpper("Post"),
		"/events/{topic}/subscriptions",
		CreatEventSubscription,
	},

	Route{
		"DeleteEventSubscription",
		strings.ToUpper("Delete"),
		"/events/subscriptions/{subscriptionId}",
		DeleteEventSubscription,
	},

	Route{
		"GetEventSubscription",
		strings.ToUpper("Get"),
		"/events/subscriptions/{subscriptionId}",
		GetEventSubscription,
	},

	Route{
		"ListEventSubscriptions",
		strings.ToUpper("Get"),
		"/events/{topic}/subscriptions",
		ListEventSubscriptions,
	},

	Route{
		"GetEvents",
		strings.ToUpper("Get"),
		"/events",
		GetEvents,
	},

	Route{
		"GetDeploymentOperation",
		strings.ToUpper("Get"),
		"/operations/{operationId}",
		GetDeploymentOperation,
	},

	Route{
		"ListOperations",
		strings.ToUpper("Get"),
		"/operations",
		ListOperations,
	},
}
