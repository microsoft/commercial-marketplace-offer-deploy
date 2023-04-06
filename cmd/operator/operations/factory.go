package operations

import (
	"fmt"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operations"
)

type DeploymentOperationFactory interface {
	Create(operationType operations.OperationType) (DeploymentOperation, error)
}

func NewDeploymentOperationFactory(appConfig *config.AppConfig) DeploymentOperationFactory {
	return &factory{
		appConfig: appConfig,
	}
}

type factory struct {
	appConfig *config.AppConfig
}

func (f *factory) Create(operationType operations.OperationType) (DeploymentOperation, error) {
	var operation DeploymentOperation

	switch operationType {
	case operations.OperationDryRun:
		operation = NewDryRunProcessor(f.appConfig)
	case operations.OperationStartDeployment:
		operation = NewStartDeploymentOperation(f.appConfig)
	}

	if operation == nil {
		return nil, fmt.Errorf("unknown operation type: %s", operationType)
	}
	return operation, nil
}
