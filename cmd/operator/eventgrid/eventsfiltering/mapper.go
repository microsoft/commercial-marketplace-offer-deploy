package eventsfiltering

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	eg "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/eventgrid"
	"github.com/mitchellh/mapstructure"
)

// maps event grid events to resources for the filter to filter out events
//
//	returns: map[resourceId]resource
type eventGridEventMapper interface {
	// Map maps event grid events to resources
	Map(ctx context.Context, events []*eventgrid.Event) eg.EventGridEventResources
}

type mapper struct {
	credential azcore.TokenCredential
}

func newEventGridEventMapper(credential azcore.TokenCredential) eventGridEventMapper {
	return &mapper{credential: credential}
}

// Map implements EventGridEventMapper
func (m *mapper) Map(ctx context.Context, events []*eventgrid.Event) eg.EventGridEventResources {
	result := eg.EventGridEventResources{}
	for _, event := range events {
		resourceId, err := m.getResourceId(event)

		if err != nil {
			continue
		}
		resource, err := m.getResource(ctx, resourceId)

		if err != nil {
			continue
		}
		result = append(result, &eg.EventGridEventResource{
			Message:  event,
			Resource: resource,
		})
	}
	return result
}

func (m *mapper) getResourceId(event *eventgrid.Event) (*arm.ResourceID, error) {
	data := eg.ResourceEventData{}
	mapstructure.Decode(event.Data, &data)

	resourceId, err := arm.ParseResourceID(data.ResourceURI)

	if err != nil {
		log.Printf("failed to parse ResourceURI: [%s], err: %v", data.ResourceURI, err)
		return nil, err
	}

	return resourceId, nil
}

func (m *mapper) getResource(ctx context.Context, resourceId *arm.ResourceID) (*armresources.GenericResource, error) {
	client, err := armresources.NewClient(resourceId.SubscriptionID, m.credential, nil)
	if err != nil {
		return nil, err
	}

	apiVersion, err := m.resolveApiVersion(ctx, resourceId)
	if err != nil {
		log.Printf("err: %v", err)
		return nil, err
	}

	response, err := client.GetByID(ctx, resourceId.String(), apiVersion, nil)
	if err != nil {
		log.Printf("failed to get associated resource: %s, err: %v", resourceId.String(), err)
		return nil, err
	}

	return &response.GenericResource, nil
}

// port from Python code from Azure CLI
// reference: https://github.com/Azure/azure-cli/blob/dev/src/azure-cli/azure/cli/command_modules/resource/custom.py

func (m *mapper) resolveApiVersion(ctx context.Context, resourceId *arm.ResourceID) (string, error) {
	defaultApiVersion := "2021-04-01"

	providerClient, err := armresources.NewProvidersClient(resourceId.SubscriptionID, m.credential, nil)
	if err != nil {
		return defaultApiVersion, err
	}

	response, err := providerClient.Get(ctx, resourceId.ResourceType.Namespace, nil)
	if err != nil {
		return defaultApiVersion, err
	}
	for _, resourceType := range response.ResourceTypes {
		isResourceTypeMatch := strings.EqualFold(*resourceType.ResourceType, resourceId.ResourceType.Type)
		if isResourceTypeMatch {
			if len(resourceType.APIVersions) > 0 {
				apiVersion := *resourceType.APIVersions[0]
				log.Printf("resolved api version: %s for resource: %s", apiVersion, resourceId.String())
				return apiVersion, nil
			}
		}
	}
	return defaultApiVersion, fmt.Errorf("failed to resolve api version for resource: %s", resourceId.String())
}
