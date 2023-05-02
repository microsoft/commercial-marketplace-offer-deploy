package operations

import (
	"context"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
)

// Executor is the interface for the actual execution of a logically invoked operation from the API
// Requestor --> invoke this operation --> enqueue --> executor --> execute the operation
type Executor interface {
	Execute(ctx context.Context, operation *data.InvokedOperation) error
}

// this is so the dry run can be tested, detaching actual dry run implementation
type DryRunFunc func(azureDeployment *deployment.AzureDeployment) (*deployment.DryRunResponse, error)
