package sdk

import (
	"context"

	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/generated"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operations"
)

type DryRunResult struct {
	Id      string
	Results map[string]any
	Status  string
}

// Performs a dry run of a deployment and returns the verification results
// returns: verification results
func (client *Client) DryRunDeployment(ctx context.Context, deploymentId int32, templateParams map[string]interface{}) (*DryRunResult, error) {
	wait := true
	operation := &generated.InvokeDeploymentOperation{
		Name:       operations.DryRunDeployment.String(),
		Parameters: templateParams,
		Wait:       &wait,
	}

	response, err := client.internalClient.InvokeDeploymentOperation(ctx, deploymentId, *operation, nil)

	if err != nil {
		return nil, err
	}
	invokedOperation := response.InvokedOperation

	result := &DryRunResult{
		Id:      *invokedOperation.ID,
		Results: invokedOperation.Result.(map[string]any),
		Status:  *invokedOperation.Status,
	}
	return result, nil
}

func (client *Client) StartDeployment(ctx context.Context, deploymentId int32) (string, error) {
	wait := false
	operation := generated.InvokeDeploymentOperation{
		Name:       operations.StartDeployment.String(),
		Parameters: nil,
		Wait:       &wait,
	}

	// TODO: implement

	_, nil := client.internalClient.InvokeDeploymentOperation(ctx, deploymentId, operation, nil)

	return "", nil
}

func (client *Client) CreateDeployment(ctx context.Context, request generated.CreateDeployment) (*generated.Deployment, error) {
	response, err := client.internalClient.CreateDeployment(ctx, request, nil)

	if err != nil {
		return nil, err
	}
	deployment := response.Deployment
	return &deployment, nil
}

func (client *Client) ListDeployments(ctx context.Context) (generated.DeploymentManagementClientListDeploymentsResponse, error) {
	return client.internalClient.ListDeployments(ctx, nil)
}
