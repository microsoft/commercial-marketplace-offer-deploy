package sdk

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operation"
)

// Gets the status of an invoked operation by deployment and operation name
//
//	deploymentId: the id of the deployment
//	operationName: the name of the operation
func (client *Client) GetDeploymentOperationStatus(ctx context.Context, deploymentId int, operationName operation.OperationType) (*StatusResponse, error) {
	items, err := client.ListInvokedOperations(ctx)

	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if *item.Name == operationName.String() && int(*item.DeploymentID) == deploymentId {
			return &StatusResponse{
				Id:           uuid.MustParse(*item.ID),
				Name:         *item.Name,
				Status:       *item.Status,
				Result:       item.Result,
				DeploymentId: int(*item.DeploymentID),
			}, nil
		}
	}
	return nil, errors.New("no operation was invoked with that deployment id and operation name")
}

// List all invoked operations
func (client *Client) ListInvokedOperations(ctx context.Context) ([]*api.InvokedOperation, error) {
	resp, err := client.internalClient.ListInvokedOperations(ctx, nil)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// Gets the status of an invoked operation, including adry run or a start deployment operation
//
//	id: the instance id of invoked operation
//
// returns: the status of the invoked operation
func (client *Client) GetStatus(ctx context.Context, id uuid.UUID) (*StatusResponse, error) {
	resp, err := client.internalClient.GetInvokedDeploymentOperation(ctx, id.String(), nil)
	if err != nil {
		return nil, err
	}
	invokedOperation := resp.InvokedOperation
	return &StatusResponse{
		Id:           uuid.MustParse(*invokedOperation.ID),
		Name:         *invokedOperation.Name,
		Status:       *invokedOperation.Status,
		Result:       invokedOperation.Result,
		DeploymentId: int(*invokedOperation.DeploymentID),
	}, nil
}

// invoke a deployment operation with parameters
func (client *Client) invokeDeploymentOperation(ctx context.Context, wait bool, operationType operation.OperationType,
	deploymentId int, parameters map[string]interface{}, retries int) (*api.InvokedOperation, error) {

	request := api.InvokeDeploymentOperationRequest{
		Name:       to.Ptr(operationType.String()),
		Parameters: parameters,
		Retries:    to.Ptr(int32(retries)),
		Wait:       &wait,
	}

	response, err := client.internalClient.InvokeDeploymentOperation(ctx, int32(deploymentId), request, nil)
	if err != nil {
		return nil, err
	}

	return response.InvokedDeploymentOperationResponse.InvokedOperation, nil
}
