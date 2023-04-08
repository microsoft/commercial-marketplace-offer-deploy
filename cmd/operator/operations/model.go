package operations

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
)

type DeploymentOperationFunc func(operation *data.InvokedOperation) error

type DeploymentOperation interface {
	Invoke(operation *data.InvokedOperation) error
}

type DryRunProcessorFunc func(azureDeployment *deployment.AzureDeployment) *deployment.DryRunResponse
