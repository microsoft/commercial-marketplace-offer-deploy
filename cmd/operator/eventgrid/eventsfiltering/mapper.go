package eventsfiltering

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
)

// maps event grid events to resources for the filter to filter out events
//
//	returns: map[resourceId]resource
type eventGridEventMapper interface {
	// Map maps event grid events to resources
	Map(ctx context.Context, events []*eventgrid.Event) eventGridEventResources
}

type mapper struct {
	credential azcore.TokenCredential
}

func newEventGridEventMapper(credential azcore.TokenCredential) eventGridEventMapper {
	return &mapper{credential: credential}
}

// Map implements EventGridEventMapper
func (m *mapper) Map(ctx context.Context, events []*eventgrid.Event) eventGridEventResources {
	result := eventGridEventResources{}

	for _, event := range events {
		resourceId, err := m.getResourceId(event)

		if err != nil {
			continue
		}
		resource, err := m.getResource(ctx, resourceId)

		if err != nil {
			continue
		}
		result = append(result, &eventGridEventResource{
			event:    event,
			resource: resource,
		})
	}
	return result
}

func (m *mapper) getResourceId(event *eventgrid.Event) (*arm.ResourceID, error) {
	data := event.Data.(resourceEventData)
	resourceId, err := arm.ParseResourceID(data.ResourceURI)

	if err != nil {
		return nil, err
	}

	return resourceId, nil
}

func (m *mapper) getResource(ctx context.Context, resourceId *arm.ResourceID) (*armresources.GenericResource, error) {
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
