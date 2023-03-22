package sdk

import (
	"context"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operations"
)

// Performs a dry run of a deployment and returns the verification results
// returns: verification results
<<<<<<< HEAD
func (client *Client) DryRunDeployment(ctx context.Context, deploymentId int64, templateParams map[string]interface{}) error {
	wait := true

	operation := internal.InvokeDeploymentOperation{
		Name:       operations.DryRunDeployment.String(),
		Parameters: templateParams,
=======
func (client *Client) DryRunDeployment(ctx context.Context, deploymentId int64, deploymentParams map[string]interface{}, templateParams map[string]interface{}) error {
	wait := true
	parameters := map[string]interface{}{
		"deploymentParams": deploymentParams,
		"templateParams": templateParams,
	}
	operation := internal.InvokeDeploymentOperation{
		Name:       operations.DryRunDeployment.String(),
		Parameters: parameters,
>>>>>>> main
		Wait:       &wait,
	}

	_, nil := client.internalClient.InvokeDeploymentOperation(ctx, deploymentId, operation, nil)

	return nil
}

func (client *Client) StartDeployment(ctx context.Context, deploymentId int64) (string, error) {
	wait := false
	operation := internal.InvokeDeploymentOperation{
		Name:       operations.StartDeployment.String(),
		Parameters: nil,
		Wait:       &wait,
	}

	// TODO: implement

	_, nil := client.internalClient.InvokeDeploymentOperation(ctx, deploymentId, operation, nil)

	return "", nil
}

func (client *Client) CreateDeployment(ctx context.Context, request internal.CreateDeployment) (*internal.Deployment, error) {
	response, err := client.internalClient.CreateDeployment(ctx, request, nil)

	if err != nil {
		return nil, err
	}
	return &response.Deployment, nil
}

func (client *Client) ListDeployments(ctx context.Context) (internal.DeploymentManagementClientListDeploymentsResponse, error) {
	return client.internalClient.ListDeployments(ctx, nil)
}
