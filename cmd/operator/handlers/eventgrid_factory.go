package handlers

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/operator/eventgrid/eventsfiltering"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
)

func NewEventGridWebHookHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		credential, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return nil
		}

		filter := newEventsFilter(credential)
		sender, err := newMessageSender(credential)

		if err != nil {
			return nil
		}

		databaseOptions := config.GetAppConfig().GetDatabaseOptions()
		db := data.NewDatabase(databaseOptions).Instance()

		return EventGridWebHook(c, db, filter, sender)
	}
}

func newMessageSender(credential azcore.TokenCredential) (messaging.MessageSender, error) {
	appConfig := config.GetAppConfig()

	sender, err := messaging.NewServiceBusMessageSender(messaging.MessageSenderOptions{
		SubscriptionId:      appConfig.Azure.SubscriptionId,
		Location:            appConfig.Azure.Location,
		ResourceGroupName:   appConfig.Azure.ResourceGroupName,
		ServiceBusNamespace: appConfig.Azure.ServiceBusNamespace,
	}, credential)

	if err != nil {
		return nil, err
	}

	return sender, nil
}

func newEventsFilter(credential azcore.TokenCredential) eventsfiltering.EventGridEventFilter {
	// TODO: probably should come from db as configurable at runtime
	includeKeys := []string{
		string(eventsfiltering.FilterTagKeyEvents),
		string(eventsfiltering.FilterTagKeyId),
		string(eventsfiltering.FilterTagKeyName),
		string(eventsfiltering.FilterTagKeyStageId),
	}
	filter := eventsfiltering.NewTagsFilter(includeKeys, credential)
	return filter
}
