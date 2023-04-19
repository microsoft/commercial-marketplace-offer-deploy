package data

import "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"

// Maps a command to create an event subscription to the EventSubscription data model
func FromCreateEventSubscription(from *api.CreateEventSubscriptionRequest) *EventSubscription {
	model := &EventSubscription{
		Name:     *from.Name,
		ApiKey:   *from.APIKey,
		Callback: *from.Callback,
	}
	return model
}
