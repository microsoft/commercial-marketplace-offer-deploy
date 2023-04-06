package data

import "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"

func FromCreateDeployment(from *api.CreateDeployment) *Deployment {
	//TODO: parse out template into the stages
	template := from.Template

	deployment := &Deployment{
		Name:           *from.Name,
		Status:         "New",
		SubscriptionId: *from.SubscriptionID,
		ResourceGroup:  *from.ResourceGroup,
		Location:       *from.Location,
		Template:       template.(map[string]interface{}),
	}
	return deployment
}

// Maps a command to create an event subscription to the EventSubscription data model
func FromCreateEventSubscription(eventType string, from *api.CreateEventSubscriptionRequest) *EventSubscription {
	model := &EventSubscription{
		Name:      *from.Name,
		EventType: eventType,
		ApiKey:    *from.APIKey,
		Callback:  *from.Callback,
	}
	return model
}
