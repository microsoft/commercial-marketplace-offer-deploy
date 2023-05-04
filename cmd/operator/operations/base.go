package operations

import (
	"context"
	"fmt"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operation"
	log "github.com/sirupsen/logrus"
)

// Executor is the interface for the actual execution of a logically invoked operation from the API
// Requestor --> invoke this operation --> enqueue --> executor --> execute the operation
type Executor interface {
	Execute(ctx context.Context, operation *data.InvokedOperation) error
}

type Execute func(ctx context.Context, operation *data.InvokedOperation) error

// this is so the dry run can be tested, detaching actual dry run implementation
type DryRunFunc func(context.Context, *deployment.AzureDeployment) (*deployment.DryRunResponse, error)

type ExecutorFactory interface {
	Create(operationType operation.OperationType) (Executor, error)
}

func NewExecutorFactory(appConfig *config.AppConfig) ExecutorFactory {
	return &factory{
		appConfig: appConfig,
	}
}

type factory struct {
	appConfig *config.AppConfig
}

func (f *factory) Create(operationType operation.OperationType) (Executor, error) {
	var executor Executor
	log.Debugf("Creating executor for operation type: %s", string(operationType))

	switch operationType {
	case operation.TypeDryRun:
		executor = NewDryRunExecutor(f.appConfig)
	case operation.TypeStartDeployment:
		executor = NewStartDeploymentExecutor(f.appConfig)
	case operation.TypeRetryDeployment:
		executor = NewRetryDeploymentExecutor(f.appConfig)
	case operation.TypeRetryStage:
		executor = NewRetryStageExecutor(f.appConfig)
	}

	if executor == nil {
		return nil, fmt.Errorf("unknown operation type: %s", operationType)
	}
	return executor, nil
}

func Trace(execute Execute) Execute {
	return func(ctx context.Context, invokedOperation *data.InvokedOperation) error {
		log.Debugf("executing '%s' [%s]", invokedOperation.Name, invokedOperation.ID.String())

		err := execute(ctx, invokedOperation)

		log.Debugf("execution done '%v' [%v]", invokedOperation.Name, invokedOperation.ID.String())

		if err != nil {
			log.Errorf("error executing %v [%v]: %v", invokedOperation.Name, invokedOperation.ID.String(), err)
		}
		return err
	}
}
