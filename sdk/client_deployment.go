package sdk

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operations"
)

type DryRunResult struct {
	Id string
	//Results map[string]any
	Results any
	Status  string
}

type StartDeploymentResult struct {
	Id     string
	Result string
	Status string
}

// Performs a dry run of a deployment and returns the verification results
// returns: verification results
func (client *Client) DryRunDeployment(ctx context.Context, deploymentId int32, templateParameters map[string]interface{}) (*DryRunResult, error) {
	invokedOperation, err := client.invokeDeploymentOperation(ctx, true, operations.OperationDryRun, deploymentId, templateParameters)
	if err != nil {
		return nil, err
	}
	if invokedOperation == nil {
		return nil, errors.New("invokedOperation is nil")
	}

	return &DryRunResult{
		Id:      *invokedOperation.ID,
		Results: invokedOperation.Result,
		Status:  *invokedOperation.Status,
	}, nil
}

func (client *Client) StartDeployment(ctx context.Context, deploymentId int32, templateParameters map[string]interface{}) (*StartDeploymentResult, error) {
	invokedOperation, err := client.invokeDeploymentOperation(ctx, false, operations.OperationStartDeployment, deploymentId, templateParameters)

	if err != nil {
		return nil, err
	}

	return &StartDeploymentResult{
		Id:     *invokedOperation.ID,
		Result: invokedOperation.Result.(string),
		Status: *invokedOperation.Status,
	}, nil
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
func (client *Client) invokeDeploymentOperation(ctx context.Context, wait bool, operationType operations.OperationType, deploymentId int32, parameters map[string]interface{}) (*api.InvokedDeploymentOperationResponse, error) {
	command := api.InvokeDeploymentOperationRequest{
		Name:       to.Ptr(operationType.String()),
		Parameters: parameters,
		Wait:       &wait,
	}

	response, err := client.internalClient.InvokeDeploymentOperation(ctx, deploymentId, command, nil)

	if err != nil {
		return nil, err
	}
	return &response.InvokedDeploymentOperationResponse, nil
}
