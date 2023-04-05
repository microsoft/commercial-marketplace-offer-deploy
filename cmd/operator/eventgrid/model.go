package eventgrid

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/eventgrid/eventgrid"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
)

// ResourceEventData is the data structure for the event grid event
// use only for unmarshalling in order to map to resource
// this is the .Data field of the event grid event when the event type is Microsoft.Resources.ResourceWriteSuccess,
// for example.
type ResourceEventData struct {
	CorrelationID    string `json:"correlationId"`
	ResourceProvider string `json:"resourceProvider"`
	ResourceURI      string `json:"resourceUri"`
	OperationName    string `json:"operationName"`
	Status           string `json:"status"`
	SubscriptionID   string `json:"subscriptionId"`
	TenantID         string `json:"tenantId"`
}

// list of event grid event resources
type EventGridEventResources []*EventGridEventResource

// combines event grid event message and the resource instance it's related to
type EventGridEventResource struct {
	// the event grid event message
	Message *eventgrid.Event

	// the azure resource instance
	Resource *armresources.GenericResource
	// the lookup tags for the resource
	Tags deployment.LookupTags
}
