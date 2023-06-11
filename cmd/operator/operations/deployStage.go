package operations

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/stage"
)

type deployStageOperation struct {
	pollerFactory *stage.DeployStagePollerFactory
}

func (op *deployStageOperation) Do(context operation.ExecutionContext) error {

	poller, err := op.pollerFactory.Create(context.Operation(), nil)
	if err != nil {
		return err
	}
	response, err := poller.PollUntilDone(context.Context())
	if err != nil {
		return err
	}

	if response.ProvisioningState != "Succeeded" {
		return operation.NewError(context.Operation())
	}

	return nil
}

func NewDeployStageOperation(appConfig *config.AppConfig) operation.OperationFunc {
	pollerFactory := stage.NewDeployStagePollerFactory()

	operation := &deployStageOperation{
		pollerFactory: pollerFactory,
	}
	return operation.Do
}
