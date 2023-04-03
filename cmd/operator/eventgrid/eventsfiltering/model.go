package eventsfiltering

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
)

type FilterTagKey string

const (
	// filter tag key for events
	FilterTagKeyEvents FilterTagKey = "modm.events"
)

type FilterTags map[string]*string

// private types

// map with key as resource id and value as event grid event resource
type eventGridEventResources []*eventGridEventResource

type eventGridEventResource struct {
	event    *eventgrid.Event
	resource *armresources.GenericResource
}

// ResourceEventData is the data structure for the event grid event
// use only for unmarshalling in order to map to resource
type resourceEventData struct {
	Authorization    any    `json:"authorization"`
	Claims           any    `json:"claims"`
	CorrelationID    string `json:"correlationId"`
	ResourceProvider string `json:"resourceProvider"`
	ResourceURI      string `json:"resourceUri"`
	OperationName    string `json:"operationName"`
	Status           string `json:"status"`
	SubscriptionID   string `json:"subscriptionId"`
	TenantID         string `json:"tenantId"`
}
