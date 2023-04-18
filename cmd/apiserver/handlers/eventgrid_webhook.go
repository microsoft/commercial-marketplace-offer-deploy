package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/eventgrid/eventsfiltering"
	w "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/eventgrid/webhookevent"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"gorm.io/gorm"
)

//region handler

// filter event grid event messages by tags that have events=true
var matchAny deployment.LookupTags = deployment.LookupTags{
	deployment.LookupTagKeyEvents: to.Ptr("true"),
}

type eventGridWebHook struct {
	db             *gorm.DB
	messageFactory *w.WebHookEventMessageFactory
	sender         messaging.MessageSender
}

// HTTP handler is the webook endpoint that receives event grid events
// the validation middleware will handle validation requests first before this is reached
func (h *eventGridWebHook) Handle(c echo.Context) error {
	log.Print("Received event grid webhook request")

	ctx := c.Request().Context()

	events := []*eventgrid.Event{}
	err := c.Bind(&events)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	messages := h.messageFactory.Create(ctx, matchAny, events)

	if len(messages) == 0 {
		log.Print("No messages were created to process")
		return c.String(http.StatusOK, "OK")
	}

	err = h.enqueueResultForProcessing(ctx, messages)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "OK")
}

// send these event grid events through our message bus to be processed and published
// to the web hook endpoints that are subscribed to our MODM events
func (h *eventGridWebHook) enqueueResultForProcessing(ctx context.Context, messages []*events.WebHookEventMessage) error {
	errors := []string{}
	for _, message := range messages {
		sendResult, err := h.sender.Send(ctx, string(messaging.QueueNameEvents), message)

		if err != nil {
			errors = append(errors, err.Error())
		}

		if len(sendResult) == 1 && sendResult[0].Error != nil {
			errors = append(errors, sendResult[0].Error.Error())
		}
	}

	return utils.NewAggregateError(errors)
}

//endregion handler

//region factory

func NewEventGridWebHookHandler(appConfig *config.AppConfig, credential azcore.TokenCredential) echo.HandlerFunc {
	log.Printf("Creating event grid webhook handler")

	return func(c echo.Context) error {
		errors := []string{}

		db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
		sender, err := newMessageSender(appConfig, credential)
		if err != nil {
			errors = append(errors, err.Error())
		}

		messageFactory, err := newWebHookEventMessageFactory(appConfig.Azure.SubscriptionId, db, credential)
		if err != nil {
			errors = append(errors, err.Error())
		}

		handler := eventGridWebHook{
			db:             db,
			messageFactory: messageFactory,
			sender:         sender,
		}

		if len(errors) > 0 {
			err = utils.NewAggregateError(errors)
			log.Printf("Failed to create event grid webhook handler: %s", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		return handler.Handle(c)
	}
}

func newWebHookEventMessageFactory(subscriptionId string, db *gorm.DB, credential azcore.TokenCredential) (*w.WebHookEventMessageFactory, error) {
	filter, err := newEventsFilter(subscriptionId, credential)
	if err != nil {
		return nil, err
	}

	client, err := armresources.NewDeploymentsClient(subscriptionId, credential, nil)
	if err != nil {
		return nil, err
	}

	return w.NewWebHookEventMessageFactory(filter, client, db), nil
}

func newMessageSender(appConfig *config.AppConfig, credential azcore.TokenCredential) (messaging.MessageSender, error) {
	sender, err := messaging.NewServiceBusMessageSender(credential, messaging.MessageSenderOptions{
		SubscriptionId:          appConfig.Azure.SubscriptionId,
		Location:                appConfig.Azure.Location,
		ResourceGroupName:       appConfig.Azure.ResourceGroupName,
		FullyQualifiedNamespace: appConfig.Azure.GetFullQualifiedNamespace(),
	})

	if err != nil {
		return nil, err
	}

	return sender, nil
}

func newEventsFilter(subscriptionId string, credential azcore.TokenCredential) (eventsfiltering.EventGridEventFilter, error) {
	// TODO: probably should come from db as configurable at runtime
	includeKeys := []string{
		string(deployment.LookupTagKeyEvents),
		string(deployment.LookupTagKeyId),
		string(deployment.LookupTagKeyName),
		string(deployment.LookupTagKeyStageId),
	}
	resourceClient, err := eventsfiltering.NewAzureResourceClient(subscriptionId, credential)
	if err != nil {
		return nil, err
	}

	provider := eventsfiltering.NewEventGridEventResourceProvider(resourceClient)
	filter := eventsfiltering.NewTagsFilter(includeKeys, provider)
	return filter, nil
}

//endregion factory
