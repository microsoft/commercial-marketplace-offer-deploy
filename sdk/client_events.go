package sdk

import (
	"context"

	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
)

// lists the event topics available to register a web hook
func (client *Client) ListEventTypes(ctx context.Context) ([]*api.EventType, error) {
	response, err := client.internalClient.GetEventTypes(ctx, nil)

	if err != nil {
		return nil, err
	}
	types := response.EventTypeArray
	return types, nil
}
