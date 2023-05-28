package operations

import (
	"fmt"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
)

func NewExecutorFactory(appConfig *config.AppConfig) operation.ExecutorFactory {
	return &factory{
		appConfig: appConfig,
	}
}

type factory struct {
	appConfig *config.AppConfig
}

func (f *factory) Create(operationType sdk.OperationType) (operation.Executor, error) {
	var operationFunc operation.OperationFunc
	log.Debugf("Creating executor for operation type: %s", string(operationType))

	switch operationType {
	case sdk.OperationDryRun:
		operationFunc = NewdryRunOperation(f.appConfig)
	case sdk.OperationDeploy:
		operationFunc = NewDeploymentOperation(f.appConfig)
	case sdk.OperationRetry: //explicit retry
		operationFunc = NewRetryDeploymentExecutor(f.appConfig)
	case sdk.OperationRetryStage:
		operationFunc = NewRetryStageExecutor(f.appConfig)
	}

	if operationFunc == nil {
		return nil, fmt.Errorf("unknown operation. Unable to execute: %s", operationType)
	}
	return operation.NewExecutor(operationFunc), nil
}
