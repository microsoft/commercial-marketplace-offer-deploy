package operations

import (
	"context"
	"fmt"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operations"
)

// Executor is the interface for the actual execution of a logically invoked operation from the API
// Requestor --> invoke this operation --> enqueue --> executor --> execute the operation
type Executor interface {
	Execute(ctx context.Context, operation *data.InvokedOperation) error
}

// this is so the dry run can be tested, detaching actual dry run implementation
type DryRunFunc func(azureDeployment *deployment.AzureDeployment) (*deployment.DryRunResponse, error)

type ExecutorFactory interface {
	Create(operationType operations.OperationType) (Executor, error)
}

func NewExecutorFactory(appConfig *config.AppConfig) ExecutorFactory {
	return &factory{
		appConfig: appConfig,
	}
}

type factory struct {
	appConfig *config.AppConfig
}

func (f *factory) Create(operationType operations.OperationType) (Executor, error) {
	var executor Executor

	switch operationType {
	case operations.OperationDryRun:
		executor = NewDryRunExecutor(f.appConfig)
	case operations.OperationStartDeployment:
		executor = NewStartDeploymentExecutor(f.appConfig)
	case operations.OperationRetryDeployment:
		executor = NewRetryDeploymentExecutor(f.appConfig)
	case operations.OperationRetryStage:
		executor = NewRetryStageExecutor(f.appConfig)
	}

	if executor == nil {
		return nil, fmt.Errorf("unknown operation type: %s", operationType)
	}
	return executor, nil
}

