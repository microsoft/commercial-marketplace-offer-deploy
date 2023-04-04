package handlers

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/eventgrid/eventsfiltering"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
)

// HTTP handler is the webook endpoint that receives event grid events
// the validation middleware will handle validation requests first before this is reached
func EventGridWebHook(c echo.Context, filter eventsfiltering.EventGridEventFilter, sender messaging.MessageSender) error {
	ctx := c.Request().Context()

	messages := []*eventgrid.Event{}
	err := c.Bind(&messages)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	matchAny := eventsfiltering.FilterTags{
		string(eventsfiltering.FilterTagKeyEvents): to.Ptr("true"),
	}

	result := filter.Filter(ctx, matchAny, messages)

	if len(result.Items) == 0 {
		return c.String(http.StatusOK, "No resources to process")
	}

	err = enqueueForPublishing(ctx, sender, &result)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "OK")
}

// send these event grid events through our message bus to be processed and published
// to the web hook endpoints that are subscribed to our MODM events
func enqueueForPublishing(ctx context.Context, sender messaging.MessageSender, result *eventsfiltering.FilterResult) error {
	results, err := sender.Send(ctx, messaging.EventsQueue, result.Items)
	if err != nil {
		return err
	}

	if len(results) > 0 {
		return utils.NewAggregateError(getErrorMessages(results))
	}
	return nil
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
