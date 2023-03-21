package sdk

import (
	"context"

	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk/internal/generated"
)

func (client *Client) ListDeployments(ctx context.Context) (generated.DeploymentManagementClientListDeploymentsResponse, error) {
	return client.internalClient.ListDeployments(ctx, nil)
}

func (client *Client) CreateDeployment(ctx context.Context, request generated.CreateDeployment) (string, error) {
	response, err := client.internalClient.CreateDeployment(ctx, request, nil)

	if err != nil {
		return "", err
	}
	return response.
}
