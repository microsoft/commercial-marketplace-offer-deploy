package azureevents

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/structure"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	log "github.com/sirupsen/logrus"
)

// maps event grid events to resources for the filter to filter out events
//
//	returns: map[resourceId]resource
type ResourceEventSubjectFactory interface {
	// Map maps event grid events to resources
	Create(ctx context.Context, events []*eventgrid.Event) []*ResourceEventSubject
}

type factory struct {
	resourceClient      AzureResourceClient
	operationRepository operation.Repository
}

func NewResourceEventSubjectFactory(resourceClient AzureResourceClient, operationRepository operation.Repository) ResourceEventSubjectFactory {
	return &factory{
		resourceClient:      resourceClient,
		operationRepository: operationRepository,
	}
}

// Map implements EventGridEventMapper
func (m *factory) Create(ctx context.Context, events []*eventgrid.Event) []*ResourceEventSubject {
	result := []*ResourceEventSubject{}
	for _, event := range events {
		subject, err := m.newSubject(ctx, event)
		if err != nil {
			log.Warnf("failed to initialize target subject: %v", err)
			continue
		}

		if subject.IsAzureDeployment() {
			azureDeployment, err := m.resourceClient.GetDeployment(ctx, subject.ResourceID())
			if err != nil {
				log.Warnf("error: %v", err)
			} else {
				subject.azureDeployment = azureDeployment
			}

			operationId, err := m.getOperationId(subject)
			if err != nil {
				log.Warnf("error: %v", err)
				continue
			}
			operation, err := m.operationRepository.First(operationId)
			if err != nil {
				log.Warnf("error: %v", err)
				continue
			}
			subject.operation = operation
		}

		result = append(result, subject)
	}
	return result
}

func (m *factory) newSubject(ctx context.Context, event *eventgrid.Event) (*ResourceEventSubject, error) {
	eventData, err := m.getEventData(event)
	if err != nil {
		return nil, err
	}

	// by parsing the resource URI we're guaranteed to have a way to get the azure resource
	// this event is for
	resourceId, err := m.getResourceId(eventData.ResourceURI)

	if err != nil {
		return nil, err
	}

	azureResource, err := m.resourceClient.Get(ctx, resourceId)

	if err != nil {
		log.Warnf("failed to get Azure resource: %v", err)
		return nil, err
	}

	subject, err := NewResourceEventSubject(eventData, event, azureResource)
	return subject, err
}

func (m *factory) getOperationId(subject *ResourceEventSubject) (uuid.UUID, error) {
	value, ok := subject.azureResource.Tags[string(deployment.LookupTagKeyOperationId)]
	if ok && value != nil {
		operationId, err := uuid.Parse(*value)
		if err != nil {
			return uuid.Nil, err
		}
		return operationId, nil
	}

	return uuid.Nil, errors.New("operationId not found")
}

func (m *factory) getEventData(event *eventgrid.Event) (*ResourceEventData, error) {
	data := ResourceEventData{}
	err := structure.Decode(event.Data, &data)
	if err != nil {
		log.Warnf("failed to decode event data: %v", err)
		return nil, err
	}
	return &data, nil
}

func (m *factory) getResourceId(resourceURI string) (*arm.ResourceID, error) {
	resourceId, err := arm.ParseResourceID(resourceURI)
	if err != nil {
		log.Warnf("failed to parse ResourceURI: [%s], err: %v", resourceURI, err)
		return nil, err
	}
	return resourceId, nil
}
