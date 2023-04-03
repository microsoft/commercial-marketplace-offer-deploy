package handlers

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/eventgrid/eventsfiltering"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"gorm.io/gorm"
)

// global app settings
var appSettings config.AppSettings = config.GetAppSettings()

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

	err = enqueueForPublishing(credential, filteredEvents, ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "OK")
}

// send these event grid events through our message bus to be processed and published
// to the web hook endpoints that are subscribed to our MODM events
func enqueueForPublishing(credential *azidentity.DefaultAzureCredential, events []*eventgrid.Event, ctx context.Context) error {
	sender, err := getMessageSender(credential)
	if err != nil {
		return err
	}
	results, err := sender.Send(ctx, messaging.EventsQueue, events)
	if err != nil {
		return err
	}

	if len(results) > 0 {
		return utils.NewAggregateError(getErrorMessages(results))
	}
	return nil
}

func getMessageSender(credential azcore.TokenCredential) (messaging.MessageSender, error) {
	sender, err := messaging.NewServiceBusMessageSender(messaging.MessageSenderOptions{
		SubscriptionId:      appSettings.Azure.SubscriptionId,
		Location:            appSettings.Azure.Location,
		ResourceGroupName:   appSettings.Azure.ResourceGroupName,
		ServiceBusNamespace: appSettings.Azure.ServiceBusNamespace,
	}, credential)
	if err != nil {
		return nil, err
	}

	return sender, nil
}

func getEventsFilter(credential azcore.TokenCredential) eventsfiltering.EventGridEventFilter {
	// TODO: this should probably come from the db
	filterTags := eventsfiltering.FilterTags{
		string(eventsfiltering.FilterTagKeyEvents): to.Ptr("true"),
	}
	filter := eventsfiltering.NewTagsFilter(filterTags, credential)

	return filter
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
