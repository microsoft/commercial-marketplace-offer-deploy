package operations

import (
	"fmt"

	"github.com/labstack/gommon/log"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

type OperationFuncProvider struct {
}

func (p *OperationFuncProvider) Get(operationType sdk.OperationType) (operation.OperationFunc, error) {
	return GetOperation(operationType)
}

func GetOperation(operationType sdk.OperationType) (operation.OperationFunc, error) {
	var operationFunc operation.OperationFunc
	log.Debugf("Creating executor for operation type: %s", string(operationType))

	switch operationType {
	case sdk.OperationDryRun:
		operationFunc = NewDryRunOperation()
	case sdk.OperationDeploy:
		operationFunc = NewDeploymentOperation()
	case sdk.OperationRetry: //explicit retry
		operationFunc = NewRetryOperation()
	case sdk.OperationRetryStage:
		operationFunc = NewRetryStageOperation()
	case sdk.OperationCancel:
		operationFunc = NewCancelOperation()
	}

	if operationFunc == nil {
		return nil, fmt.Errorf("unknown operation. Unable to execute: %s", operationType)
	}
	return operationFunc, nil
}
