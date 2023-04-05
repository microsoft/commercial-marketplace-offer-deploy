package webhookevent

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	"github.com/google/uuid"
	eg "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/eventgrid"
	filtering "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/eventgrid/eventsfiltering"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"gorm.io/gorm"
)

type WebHookEventMessageFactory struct {
	filter filtering.EventGridEventFilter
	db     *gorm.DB
}

func NewWebHookEventMessageFactory(filter filtering.EventGridEventFilter, db *gorm.DB) *WebHookEventMessageFactory {
	return &WebHookEventMessageFactory{
		filter: filter,
		db:     db,
	}
}

// Creates a list of WebHookEventMessage from a list of EventGridEventResource
func (f *WebHookEventMessageFactory) Create(ctx context.Context, matchAny deployment.LookupTags, eventGridEvents []*eventgrid.Event) []*events.WebHookEventMessage {
	result := f.filter.Filter(ctx, matchAny, eventGridEvents)
	messages := []*events.WebHookEventMessage{}

	for _, item := range result {
		message, err := f.convert(item)
		if err != nil {
			continue
		}

		messages = append(messages, message)
	}

	return messages
}

func (f *WebHookEventMessageFactory) convert(item *eg.EventGridEventResource) (*events.WebHookEventMessage, error) {
	deployment := f.getRelatedDeployment(item)

	messageId, _ := uuid.Parse(*item.Message.ID)
	eventData := item.Message.Data.(eg.ResourceEventData)

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
	// if stage != nil {
	// 	stageId = stage.ID
	// }
	message.SetSubject(int(deployment.ID), &stageId)

	return message, nil
}

func (f *WebHookEventMessageFactory) getRelatedDeployment(item *eg.EventGridEventResource) *data.Deployment {
	panic("unimplemented")
}
