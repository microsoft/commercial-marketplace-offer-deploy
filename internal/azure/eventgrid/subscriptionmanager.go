package eventgrid

import (
	"github.com/google/uuid"
)

type EventGridSubscriptionCredential struct {
	ClientId     uuid.UUID
	ClientSecret *string
	TenantId     uuid.UUID
}

// Subscription manager.
//
//	notes: we're going to need to register with the credentials in order to secure the web hook's endpoint properly
//	see: https://learn.microsoft.com/en-us/azure/event-grid/secure-webhook-delivery#configure-the-event-subscription-by-using-an-azure-ad-application
type SubscriptionManager interface {
	Register(credential *EventGridSubscriptionCredential)
}
