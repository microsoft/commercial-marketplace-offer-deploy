package webhookevent

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	"github.com/google/uuid"
	eg "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/eventgrid"
	filtering "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/eventgrid/eventsfiltering"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	d "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

// this factory is intented to create a list of WebHookEventMessages from a list of EventGridEventResources
// so the messages can be relayed via queue to be published to MODM consumer webhook subscription
type WebHookEventMessageFactory struct {
	client *armresources.DeploymentsClient
	filter filtering.EventGridEventFilter
	db     *gorm.DB
}

func NewWebHookEventMessageFactory(filter filtering.EventGridEventFilter, client *armresources.DeploymentsClient, db *gorm.DB) *WebHookEventMessageFactory {
	return &WebHookEventMessageFactory{
		client: client,
		filter: filter,
		db:     db,
	}
}

// Creates a list of WebHookEventMessage from a list of EventGridEventResource
func (f *WebHookEventMessageFactory) Create(ctx context.Context, matchAny d.LookupTags, eventGridEvents []*eventgrid.Event) []*events.WebHookEventMessage {
	result := f.filter.Filter(ctx, matchAny, eventGridEvents)
	messages := []*events.WebHookEventMessage{}

	log.Printf("factory received %d EventGridEvents but filtered to %d messages", len(eventGridEvents), len(result))

	for _, item := range result {
		message, err := f.convert(item)
		if err != nil {
			continue
		}

		messages = append(messages, message)
	}

	return messages
}

//region private methods

func (f *WebHookEventMessageFactory) convert(item *eg.EventGridEventResource) (*events.WebHookEventMessage, error) {
	deployment, err := f.getRelatedDeployment(item)
	if err != nil {
		return nil, err
	}

	messageId, _ := uuid.Parse(*item.Message.ID)

	eventData := eg.ResourceEventData{}
	mapstructure.Decode(item.Message.Data, &eventData)

	message := &events.WebHookEventMessage{
		Id:             messageId,
		SubscriptionId: [16]byte{},
		EventType:      *item.Message.EventType,
		Body: events.WebHookDeploymentEventMessageBody{
			ResourceId: *item.Resource.ID,
			Status:     eventData.Status,
			Message:    eventData.OperationName,
		},
	}

	var stageId uuid.UUID
	if item.Tags[d.LookupTagKeyStageId] != nil {
		value := *item.Tags[d.LookupTagKeyStageId]
		stageId, _ = uuid.Parse(value)
	}
	message.SetSubject(int(deployment.ID), &stageId, item.Resource.Name)

	return message, nil
}

func (f *WebHookEventMessageFactory) getRelatedDeployment(item *eg.EventGridEventResource) (*data.Deployment, error) {
	eventData := eg.ResourceEventData{}
	mapstructure.Decode(item.Message.Data, &eventData)

	correlationId := eventData.CorrelationID
	resourceId, err := arm.ParseResourceID(*item.Resource.ID)
	if err != nil {
		return nil, err
	}

	pager := f.client.NewListByResourceGroupPager(resourceId.ResourceGroupName, nil)
	deploymentId, err := f.lookupDeploymentId(context.Background(), correlationId, pager)
	if err != nil {
		return nil, err
	}

	deployment := &data.Deployment{}
	tx := f.db.First(deployment, deploymentId)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return deployment, nil
}

func (f *WebHookEventMessageFactory) lookupDeploymentId(ctx context.Context, correlationId string, pager *runtime.Pager[armresources.DeploymentsClientListByResourceGroupResponse]) (*int, error) {
	deployment := &data.Deployment{}

	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		if nextResult.DeploymentListResult.Value != nil {
			for _, item := range nextResult.DeploymentListResult.Value {
				correlationIdMatches := strings.EqualFold(*item.Properties.CorrelationID, correlationId)
				if correlationIdMatches {
					id, err := deployment.ParseAzureDeploymentName(deployment.Name)
					if err != nil {
						// the name didn't match our pattern so we're not interested in this azure deployment, keep searching for a match
						// until we find 1=1 for our deployment (the top level "main deployment")
						continue
					} else {
						return id, nil
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("deployment not found for correlationId: %s", correlationId)
}

//endregion private methods
