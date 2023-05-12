package sdk

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operation"
)

// Gets the status of a deployment operation, e,g. a dry run or a start deployment operation
//
//	id: the instance id of the deployment operation
func (client *Client) GetStatus(ctx context.Context, instanceId uuid.UUID) (*StatusResponse, error) {
	resp, err := client.internalClient.GetInvokedDeploymentOperation(ctx, instanceId.String(), nil)
	if err != nil {
		return nil, err
	}
	response := resp.InvokedDeploymentOperationResponse
	return &StatusResponse{
		Id:           uuid.MustParse(*response.ID),
		Name:         *response.Name,
		Status:       *response.Status,
		Result:       response.Result,
		DeploymentId: int(*response.DeploymentID),
	}, nil
}

// invoke a deployment operation with parameters
func (client *Client) invokeDeploymentOperation(ctx context.Context, wait bool, operationType operation.OperationType,
	deploymentId int, parameters map[string]interface{}, retries int) (*api.InvokedDeploymentOperationResponse, error) {

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

	return &response.InvokedDeploymentOperationResponse, nil
}
