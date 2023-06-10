package operations

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
)

type deployStageOperation struct {
}

func (op *deployStageOperation) Do(context operation.ExecutionContext) error {
	return nil
}

func NewDeployStageOperation(appConfig *config.AppConfig) operation.OperationFunc {

	operation := &deployStageOperation{}
	return operation.Do
}
