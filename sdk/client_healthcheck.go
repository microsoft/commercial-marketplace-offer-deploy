package sdk

import (
	"context"
)

func (client *Client) HealthStatus(ctx context.Context) (*GetHealthResponse, error) {
	resp, err := client.internalClient.GetHealth(ctx, nil)
	if err != nil {
		return nil, err
	}
	healthResponse := resp.GetHealthResponse
	return &healthResponse, nil
}