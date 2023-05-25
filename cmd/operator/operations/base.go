package operations

import (
	"context"
	"fmt"
	"time"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
)

// RetriableError is a custom error that contains a positive duration for the next retry
type RetriableError struct {
	Err        error
	RetryAfter time.Duration
}

// Error returns error message and a Retry-After duration
func (e *RetriableError) Error() string {
	return fmt.Sprintf("%s (retry after %v)", e.Err.Error(), e.RetryAfter)
}

// Executor is the interface for the actual execution of a logically invoked operation from the API
// Requestor --> invoke this operation --> enqueue --> executor --> execute the operation
type Executor interface {
	Execute(ctx context.Context, operation *model.InvokedOperation) error
}

type Execute func(ctx context.Context, operation *model.InvokedOperation) error

// this is so the dry run can be tested, detaching actual dry run implementation
type DryRunFunc func(context.Context, *deployment.AzureDeployment) (*sdk.DryRunResult, error)

type ExecutorFactory interface {
	Create(operationType sdk.OperationType) (Executor, error)
}

func NewExecutorFactory(appConfig *config.AppConfig) ExecutorFactory {
	return &factory{
		appConfig: appConfig,
	}
}

type factory struct {
	appConfig *config.AppConfig
}

func (f *factory) Create(operationType sdk.OperationType) (Executor, error) {
	var executor Executor
	log.Debugf("Creating executor for operation type: %s", string(operationType))

	switch operationType {
	case sdk.OperationDryRun:
		executor = NewDryRunExecutor(f.appConfig)
	case sdk.OperationDeploy:
		executor = NewStartDeploymentExecutor(f.appConfig)
	case sdk.OperationRetry:
		executor = NewRetryDeploymentExecutor(f.appConfig)
	case sdk.OperationRetryStage:
		executor = NewRetryStageExecutor(f.appConfig)
	}

	if executor == nil {
		return nil, fmt.Errorf("unknown operation type: %s", operationType)
	}
	return executor, nil
}

func Trace(execute Execute) Execute {
	return func(ctx context.Context, invokedOperation *model.InvokedOperation) error {
		logger := log.WithFields(
			log.Fields{
				"name": invokedOperation.Name,
				"id":   invokedOperation.ID.String(),
			})
		logger.Debug("execution started")
		err := execute(ctx, invokedOperation)
		logger.Debug("execution done")

		if err != nil {
			log.WithError(err).Error("execution failed")
		}
		return err
	}
}
