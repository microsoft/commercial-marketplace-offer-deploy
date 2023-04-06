package operations

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
)

type DeploymentOperationHandlerFunc func(operation *data.InvokedOperation) error

type DeploymentOperationHandler interface {
	Handle(operation *data.InvokedOperation) error
}

type DryRunProcessorFunc func(azureDeployment *deployment.AzureDeployment) *deployment.DryRunResponse
