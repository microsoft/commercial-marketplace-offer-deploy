package data

import "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"

// Maps a command to create an event subscription to the hook data model
func FromCreateEventHook(from *api.CreateEventHookRequest) *EventHook {
	model := &EventHook{
		Name:     *from.Name,
		ApiKey:   *from.APIKey,
		Callback: *from.Callback,
	}
	return model
}
