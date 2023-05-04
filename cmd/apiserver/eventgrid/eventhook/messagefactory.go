package eventhook

import (
	"context"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	"github.com/google/uuid"
	eg "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/eventgrid"
	filtering "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/eventgrid/eventsfiltering"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	d "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

// this factory is intented to create a list of WebHookEventMessages from a list of EventGridEventResources
// so the messages can be relayed via queue to be published to MODM consumer webhook subscription
type EventHookMessageFactory struct {
	client *armresources.DeploymentsClient
	filter filtering.EventGridEventFilter
	db     *gorm.DB
}

func NewEventHookMessageFactory(filter filtering.EventGridEventFilter, client *armresources.DeploymentsClient, db *gorm.DB) *EventHookMessageFactory {
	return &EventHookMessageFactory{
		client: client,
		filter: filter,
		db:     db,
	}
}

// Creates a list of messages from a list of EventGridEventResource
func (f *EventHookMessageFactory) Create(ctx context.Context, matchAny d.LookupTags, eventGridEvents []*eventgrid.Event) []*events.EventHookMessage {
	result := f.filter.Filter(ctx, matchAny, eventGridEvents)
	messages := []*events.EventHookMessage{}

	log.Debugf("factory received %d EventGridEvents, filtered to %d messages", len(eventGridEvents), len(result))

	for _, item := range result {
		message, err := f.convert(item)
		if err != nil {
			log.Error("failed to convert EventGridEventResource to WebHookEventMessage: %s", err.Error())
		}

		messages = append(messages, message)
	}

	return messages
}

//region private methods

func (f *EventHookMessageFactory) convert(item *eg.EventGridEventResource) (*events.EventHookMessage, error) {
	deployment, err := f.getRelatedDeployment(item)
	if err != nil {
		return nil, err
	}

	messageId, _ := uuid.Parse(*item.Message.ID)
	eventData := eg.ResourceEventData{}
	mapstructure.Decode(item.Message.Data, &eventData)

	message := &events.EventHookMessage{
		Id:     messageId,
		Status: strcase.ToLowerCamel(eventData.Status),
		Type:   string(events.EventTypeDeploymentAzureEventReceived),
		Data:   eventData,
	}

	var stageId uuid.UUID
	if item.Tags[d.LookupTagKeyStageId] != nil {
		value := *item.Tags[d.LookupTagKeyStageId]
		stageId, _ = uuid.Parse(value)
	}
	message.SetSubject(int(deployment.ID), &stageId)

	return message, nil
}

func (f *EventHookMessageFactory) getRelatedDeployment(item *eg.EventGridEventResource) (*data.Deployment, error) {
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

func (f *EventHookMessageFactory) lookupDeploymentId(ctx context.Context, correlationId string, pager *runtime.Pager[armresources.DeploymentsClientListByResourceGroupResponse]) (*int, error) {
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
					id, err := deployment.ParseAzureDeploymentName(*item.Name)
					if err != nil {
						fmt.Printf("correlation: resource/%s/correlationId/%s", *item.Name, *item.Properties.CorrelationID)

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
