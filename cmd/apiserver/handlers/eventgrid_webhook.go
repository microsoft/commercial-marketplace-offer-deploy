package handlers

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	"github.com/labstack/echo/v4"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/eventgrid/eventhook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/eventgrid/eventsfiltering"
	filter "github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/eventgrid/eventsfiltering"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"gorm.io/gorm"
)

//region handler

// filter event grid event messages by tags that have events=true
var matchAny deployment.LookupTags = deployment.LookupTags{
	deployment.LookupTagKeyEvents: to.Ptr("true"),
}

type eventGridWebHook struct {
	messageFactory *eventhook.EventHookMessageFactory
	filter         filter.EventGridEventFilter
}

// HTTP handler is the webook endpoint that receives event grid events
// the validation middleware will handle validation requests first before this is reached
func (h *eventGridWebHook) Handle(c echo.Context) error {
	log.Debug("Received event grid webhook")

	ctx := c.Request().Context()

	events := []*eventgrid.Event{}
	err := c.Bind(&events)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resources := h.filter.Filter(ctx, matchAny, events)

	if len(resources) == 0 {
		return c.String(http.StatusOK, "OK")
	}

	//

	messages := h.messageFactory.Create(ctx, resources)
	log.Debugf("Event Hook messages total: %d", len(messages))

	if len(messages) == 0 {
		return c.String(http.StatusOK, "OK")
	}

	err = h.add(ctx, messages)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "OK")
}

// send these event grid events through our message bus to be processed and published
// to the web hook endpoints that are subscribed to our MODM events
func (h *eventGridWebHook) add(ctx context.Context, messages []*sdk.EventHookMessage) error {
	errors := []string{}
	for _, message := range messages {
		log.Debugf("Adding event hook message: %+v", message)
		err := hook.Add(ctx, message)

		if err != nil {
			errors = append(errors, err.Error())
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
		messageFactory, err := newWebHookEventMessageFactory(appConfig.Azure.SubscriptionId, db, credential)
		if err != nil {
			errors = append(errors, err.Error())
		}

		eventsFilter, err := newEventsFilter(appConfig.Azure.SubscriptionId, credential)
		if err != nil {
			errors = append(errors, err.Error())
		}

		handler := eventGridWebHook{
			messageFactory: messageFactory,
			filter:         eventsFilter,
		}

		if len(errors) > 0 {
			err = utils.NewAggregateError(errors)
			log.Errorf("Failed to create event grid webhook handler: %s", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		return handler.Handle(c)
	}
}

func newWebHookEventMessageFactory(subscriptionId string, db *gorm.DB, credential azcore.TokenCredential) (*eventhook.EventHookMessageFactory, error) {
	client, err := armresources.NewDeploymentsClient(subscriptionId, credential, nil)
	if err != nil {
		return nil, err
	}

	return eventhook.NewEventHookMessageFactory(client, db), nil
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
