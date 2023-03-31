package sdk

import (
	"context"

	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operations"
)

type DryRunResult struct {
	Id      string
	Results map[string]any
	Status  string
}

type StartDeploymentResult struct {
	Id      string
	Results map[string]any
	Status  string
}

// Performs a dry run of a deployment and returns the verification results
// returns: verification results
func (client *Client) DryRunDeployment(ctx context.Context, deploymentId int32, templateParameters map[string]interface{}) (*DryRunResult, error) {
	invokedOperation, err := client.invokeDeploymentOperation(ctx, true, operations.DryRunDeploymentOperation, deploymentId, templateParameters)

	if err != nil {
		return nil, err
	}

	return &DryRunResult{
		Id:      *invokedOperation.ID,
		Results: invokedOperation.Result.(map[string]any),
		Status:  *invokedOperation.Status,
	}, nil
}

func (client *Client) StartDeployment(ctx context.Context, deploymentId int32, templateParameters map[string]interface{}) (*StartDeploymentResult, error) {
	invokedOperation, err := client.invokeDeploymentOperation(ctx, false, operations.StartDeploymentOperation, deploymentId, templateParameters)

	if err != nil {
		return nil, err
	}

	return &StartDeploymentResult{
		Id:      *invokedOperation.ID,
		Results: invokedOperation.Result.(map[string]any),
		Status:  *invokedOperation.Status,
	}, nil
	// TODO: implement start deployment and return the invocation result
	//return invokedOperation.Name, nil
}

func (client *Client) CreateDeployment(ctx context.Context, request api.CreateDeployment) (*api.Deployment, error) {
	response, err := client.internalClient.CreateDeployment(ctx, request, nil)

	if err != nil {
		return nil, err
	}
	deployment := response.Deployment
	return &deployment, nil
}

func (client *Client) ListDeployments(ctx context.Context) (api.DeploymentManagementClientListDeploymentsResponse, error) {
	return client.internalClient.ListDeployments(ctx, nil)
}

// invoke a deployment operation with parameters
func (client *Client) invokeDeploymentOperation(ctx context.Context, wait bool, operationType operations.OperationType, deploymentId int32, parameters map[string]interface{}) (*api.InvokedOperation, error) {
	operationTypeName := operations.DryRunDeploymentOperation.String()
	command := &api.InvokeDeploymentOperation{
		Name:       &operationTypeName,
		Parameters: parameters,
		Wait:       &wait,
	}

	response, err := client.internalClient.InvokeDeploymentOperation(ctx, deploymentId, *command, nil)

	if err != nil {
		return nil, err
	}
	return &response.InvokedOperation, nil
}
