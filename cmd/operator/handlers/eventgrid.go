package handlers

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	internal "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/eventgrid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/eventgrid/eventsfiltering"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"gorm.io/gorm"
)

// HTTP handler is the webook endpoint that receives event grid events
// the validation middleware will handle validation requests first before this is reached
func EventGridWebHook(c echo.Context, db *gorm.DB, filter eventsfiltering.EventGridEventFilter, sender messaging.MessageSender) error {
	ctx := c.Request().Context()

	messages := []*eventgrid.Event{}
	err := c.Bind(&messages)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	result := filterEvents(ctx, filter, messages)
	if len(result.Items) == 0 {
		return c.String(http.StatusOK, "No resources to process")
	}

	err = enqueueResultForProcessing(ctx, sender, &result)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "OK")
}

func filterEvents(ctx context.Context, filter eventsfiltering.EventGridEventFilter, messages []*eventgrid.Event) eventsfiltering.FilterResult {
	result := filter.Filter(ctx, eventsfiltering.FilterTags{
		eventsfiltering.FilterTagKeyEvents: to.Ptr("true"),
	}, messages)
	return result
}

// send these event grid events through our message bus to be processed and published
// to the web hook endpoints that are subscribed to our MODM events
func enqueueResultForProcessing(ctx context.Context, sender messaging.MessageSender, result *eventsfiltering.FilterResult) error {
	if len(result.Items) == 0 {
		return nil
	}

	messages := mapFilterResultToMessages(result)
	sendResults, err := sender.Send(ctx, "events", messages)

	if err != nil {
		return err
	}

	errors := getErrorMessages(sendResults)
	return utils.NewAggregateError(errors)
}

func getErrorMessages(sendResults []messaging.SendMessageResult) *[]string {
	errors := []string{}
	for _, result := range sendResults {
		if result.Error != nil {
			errors = append(errors, result.Error.Error())
		}
	}

	return &errors
}

func mapFilterResultToMessages(result *eventsfiltering.FilterResult) []*events.WebHookEventMessage {
	messages := []*events.WebHookEventMessage{}

	for _, item := range result.Items {
		message, err := mapTo(&item)
		if err != nil {
			continue
		}

		messages = append(messages, message)
	}

	return messages
}

func mapTo(item *eventsfiltering.FilterResultItem) (*events.WebHookEventMessage, error) {
	deployment := getRelatedDeployment(item)
	stage := getRelatedStage(item)

	messageId, _ := uuid.Parse(*item.EventGridEvent.ID)
	eventData := item.EventGridEvent.Data.(internal.ResourceEventData)

	message := &events.WebHookEventMessage{
		Id:             messageId,
		SubscriptionId: [16]byte{},
		EventType:      *item.EventGridEvent.EventType,
		Body: events.WebHookDeploymentEventMessageBody{
			ResourceId: *item.Resource.ID,
			Status:     eventData.Status,
			Message:    eventData.OperationName,
		},
	}

	var stageId uuid.UUID
	if stage != nil {
		stageId = stage.ID
	}
	message.SetSubject(int(deployment.ID), &stageId)

	return message, nil
}

func getRelatedStage(item *eventsfiltering.FilterResultItem) *data.Stage {
	//correlationId := item.EventGridEvent.Data.(internal.ResourceEventData).CorrelationID

	return nil
}

func getRelatedDeployment(item *eventsfiltering.FilterResultItem) *data.Deployment {
	// correlationId := item.EventGridEvent.Data.(internal.ResourceEventData).CorrelationIDs

	return nil
}
