package mapper

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
)

type InvokedDeploymentOperationResponseMapper struct {
}

func (m *InvokedDeploymentOperationResponseMapper) Map(invokedOperation *data.InvokedOperation) api.InvokedDeploymentOperationResponse {
	result := api.InvokedDeploymentOperationResponse{
		InvokedOperation: &api.InvokedOperation{
			DeploymentID: to.Ptr(int32(invokedOperation.DeploymentId)),
			ID:           to.Ptr(invokedOperation.ID.String()),
			InvokedOn:    to.Ptr(invokedOperation.CreatedAt),
			Name:         &invokedOperation.Name,
			Parameters:   &invokedOperation.Parameters,
			Result:       &invokedOperation.Result,
			Status:       &invokedOperation.Status,
		},
	}
	return result
}

type InvokedDeploymentResponseMapper struct {
}

func (m *InvokedDeploymentResponseMapper) MapList(items []data.InvokedOperation) *api.ListInvokedOperationResponse {
	response := &api.ListInvokedOperationResponse{
		Items: []*api.InvokedOperation{},
	}

	for _, item := range items {
		response.Items = append(response.Items, m.Map(item).InvokedOperation)
	}
	return response
}

func (m *InvokedDeploymentResponseMapper) Map(invokedOperation data.InvokedOperation) *api.GetInvokedOperationResponse {
	return &api.GetInvokedOperationResponse{
		InvokedOperation: &api.InvokedOperation{
			DeploymentID: to.Ptr(int32(invokedOperation.DeploymentId)),
			ID:           to.Ptr(invokedOperation.ID.String()),
			InvokedOn:    to.Ptr(invokedOperation.CreatedAt),
			Name:         &invokedOperation.Name,
			Parameters:   &invokedOperation.Parameters,
			Result:       &invokedOperation.Result,
			Status:       &invokedOperation.Status,
		},
	}
}
