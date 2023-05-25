package mapper

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

type InvokedDeploymentOperationResponseMapper struct {
}

func (m *InvokedDeploymentOperationResponseMapper) Map(invokedOperation *model.InvokedOperation) sdk.InvokedDeploymentOperationResponse {

	response := sdk.InvokedDeploymentOperationResponse{
		InvokedOperation: &sdk.InvokedOperation{
			DeploymentID: to.Ptr(int32(invokedOperation.DeploymentId)),
			ID:           to.Ptr(invokedOperation.ID.String()),
			InvokedOn:    to.Ptr(invokedOperation.CreatedAt),
			Name:         &invokedOperation.Name,
			Parameters:   &invokedOperation.Parameters,
			Result:       to.Ptr(invokedOperation.LatestResult()),
			Status:       &invokedOperation.Status,
		},
	}
	return response
}

type InvokedDeploymentResponseMapper struct {
}

func (m *InvokedDeploymentResponseMapper) MapList(items []model.InvokedOperation) *sdk.ListInvokedOperationResponse {
	response := &sdk.ListInvokedOperationResponse{
		Items: []*sdk.InvokedOperation{},
	}

	for _, item := range items {
		response.Items = append(response.Items, m.Map(item).InvokedOperation)
	}
	return response
}

func (m *InvokedDeploymentResponseMapper) Map(invokedOperation model.InvokedOperation) *sdk.GetInvokedOperationResponse {
	return &sdk.GetInvokedOperationResponse{
		InvokedOperation: &sdk.InvokedOperation{
			DeploymentID: to.Ptr(int32(invokedOperation.DeploymentId)),
			ID:           to.Ptr(invokedOperation.ID.String()),
			InvokedOn:    to.Ptr(invokedOperation.CreatedAt),
			Name:         &invokedOperation.Name,
			Parameters:   &invokedOperation.Parameters,
			Result:       to.Ptr(invokedOperation.LatestResult()),
			Status:       &invokedOperation.Status,
		},
	}
}
