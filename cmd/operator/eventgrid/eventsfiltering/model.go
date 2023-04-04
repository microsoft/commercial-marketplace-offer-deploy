package eventsfiltering

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
)

type FilterTagKey string

const (
	// filter tag key for events

	// The unique id for modm to identify something
	FilterTagKeyId FilterTagKey = "modm.id"

	// whether or not to send events, if this is not set to true, then the event will not be sent
	FilterTagKeyEvents FilterTagKey = "modm.events"

	// the friendly name of the resource used for logging
	FilterTagKeyName FilterTagKey = "modm.name"

	// the stage id reference. Use is on a resource that's a child of a 1-level parent deployment
	FilterTagKeyStageId FilterTagKey = "modm.stage.id"
)

type FilterResult struct {
	Items []FilterResultItem
}

type FilterResultItem struct {
	EventGridEvent eventgrid.Event
	Resource       armresources.GenericResource

	// tags that matched the filter
	MatchedTags map[FilterTagKey]string
}

type FilterTags map[FilterTagKey]*string

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
