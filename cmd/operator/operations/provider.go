package operations

import (
	"fmt"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
)

type OperationTaskProvider struct {
	appConfig *config.AppConfig
}

func NewOperationTaskProvider(appConfig *config.AppConfig) *OperationTaskProvider {
	return &OperationTaskProvider{
		appConfig: appConfig,
	}
}

func (p *OperationTaskProvider) Get(operationType sdk.OperationType) (operation.OperationTask, error) {
	return GetOperationTask(operationType, p.appConfig)
}

func GetOperationTask(operationType sdk.OperationType, appConfig *config.AppConfig) (operation.OperationTask, error) {
	log.Tracef("Creating operation task for type: %s", string(operationType))

	var task operation.OperationTask

	switch operationType {
	case sdk.OperationDryRun:
		task = NewDryRunTask()
	case sdk.OperationDeploy:
		task = NewDeployTask(appConfig)
	case sdk.OperationDeployStage:
		task = NewDeployStageOperation(appConfig)
	case sdk.OperationRetry: //explicit retry
		task = NewRetryTask()
	case sdk.OperationRetryStage:
		task = NewRetryStageTask()
	case sdk.OperationCancel:
		task = NewCancelTask()
	}

	if task == nil {
		return nil, fmt.Errorf("unknown operation. Unable to execute: %s", operationType)
	}
	return task, nil
}
