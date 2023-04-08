package eventsfiltering

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	eg "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/eventgrid"
	"github.com/mitchellh/mapstructure"
)

// maps event grid events to resources for the filter to filter out events
//
//	returns: map[resourceId]resource
type EventGridEventResourceProvider interface {
	// Map maps event grid events to resources
	Get(ctx context.Context, events []*eventgrid.Event) eg.EventGridEventResources
}

type provider struct {
	resourceClient AzureResourceClient
}

func NewEventGridEventResourceProvider(resourceClient AzureResourceClient) EventGridEventResourceProvider {
	return &provider{resourceClient: resourceClient}
}

// Map implements EventGridEventMapper
func (m *provider) Get(ctx context.Context, events []*eventgrid.Event) eg.EventGridEventResources {
	result := eg.EventGridEventResources{}
	for _, event := range events {
		resourceId, err := m.getResourceId(event)

		if err != nil {
			continue
		}
		resource, err := m.resourceClient.Get(ctx, resourceId)

		if err != nil {
			log.Printf("error: %v", err)
			continue
		}
		result = append(result, &eg.EventGridEventResource{
			Message:  event,
			Resource: resource,
		})
	}
	return result
}

func (m *provider) getResourceId(event *eventgrid.Event) (*arm.ResourceID, error) {
	data := eg.ResourceEventData{}
	mapstructure.Decode(event.Data, &data)

	resourceId, err := arm.ParseResourceID(data.ResourceURI)

	if err != nil {
		log.Printf("failed to parse ResourceURI: [%s], err: %v", data.ResourceURI, err)
		return nil, err
	}

	return resourceId, nil
}
