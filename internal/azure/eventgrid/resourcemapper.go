package eventgrid

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
)

type EventGridEventMapper interface {
	// Map maps event grid events to resources
	Map(ctx context.Context, events []*eventgrid.Event) []*armresources.GenericResource
}

type eventGridEventMapper struct {
	credential azcore.TokenCredential
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

func NewEventGridEventMapper(credential azcore.TokenCredential) EventGridEventMapper {
	return &eventGridEventMapper{credential: credential}
}

// Map implements EventGridEventMapper
func (m *eventGridEventMapper) Map(ctx context.Context, events []*eventgrid.Event) []*armresources.GenericResource {
	resources := []*armresources.GenericResource{}

	for _, event := range events {
		resourceId, err := m.getResourceId(event)

		if err != nil {
			continue
		}
		resource, err := m.getResource(ctx, resourceId)

		if err != nil {
			continue
		}
		resources = append(resources, resource)

	}
	return resources
}

func (m *eventGridEventMapper) getResourceId(event *eventgrid.Event) (*arm.ResourceID, error) {
	data := event.Data.(resourceEventData)
	resourceId, err := arm.ParseResourceID(data.ResourceURI)

	if err != nil {
		return nil, err
	}

	return resourceId, nil
}

func (m *eventGridEventMapper) getResource(ctx context.Context, resourceId *arm.ResourceID) (*armresources.GenericResource, error) {
	client, err := armresources.NewClient(resourceId.SubscriptionID, m.credential, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.GetByID(ctx, resourceId.String(), "2021-04-01", nil)
	if err != nil {
		return nil, err
	}

	return &response.GenericResource, nil
}
