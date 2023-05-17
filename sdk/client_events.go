package sdk

import (
	"context"
)

// lists the event topics available to register a web hook
func (client *Client) ListEventTypes(ctx context.Context) ([]*EventType, error) {
	response, err := client.internalClient.GetEventTypes(ctx, nil)

	if err != nil {
		return nil, err
	}
	types := response.EventTypeArray
	return types, nil
}

func (client *Client) CreateEventHook(ctx context.Context, request CreateEventHookRequest) (*CreateEventHookResponse, error) {
	response, err := client.internalClient.CreateEvenHook(ctx, request, nil)

	if err != nil {
		return nil, err
	}
	return &response.CreateEventHookResponse, nil
}
