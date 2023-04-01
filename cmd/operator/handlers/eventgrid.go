package handlers

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/eventgrid/eventsfiltering"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"gorm.io/gorm"
)

// HTTP handler is the webook endpoint that receives event grid events
// the validation middleware will handle validation requests first before this is reached
func EventGridWebHook(c echo.Context, db *gorm.DB) error {
	ctx := c.Request().Context()
	messages := []*eventgrid.Event{}
	err := c.Bind(&messages)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil
	}

	filter := getEventsFilter(credential)
	filteredEvents := filter.Filter(ctx, messages)

	if len(filteredEvents) == 0 {
		return c.String(http.StatusOK, "No resources to process")
	}

	enqueueForPublishing(credential, filteredEvents, ctx)
	return c.String(http.StatusOK, "OK")
}

// send these event grid events through our message bus to be processed and published
// to the web hook endpoints that are subscribed to our MODM events
func enqueueForPublishing(credential *azidentity.DefaultAzureCredential, events []*eventgrid.Event, ctx context.Context) {
	sender := getMessageSender(credential)
	sender.Send(ctx, events)
}

func getMessageSender(credential azcore.TokenCredential) messaging.MessageSender {
	return messaging.NewMessageSender(credential)
}

func getEventsFilter(credential azcore.TokenCredential) eventsfiltering.EventGridEventFilter {
	// TODO: this should probably come from the db
	filterTags := eventsfiltering.FilterTags{
		string(eventsfiltering.FilterTagKeyEvents): to.Ptr("true"),
	}
	filter := eventsfiltering.NewTagsFilter(filterTags, credential)

	return filter
}
