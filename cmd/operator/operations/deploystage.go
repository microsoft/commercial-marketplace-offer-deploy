package operations

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/stage"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

type nameFinderFactory func(context operation.ExecutionContext) (*operation.AzureDeploymentNameFinder, error)

type deployStageOperation struct {
	pollerFactory     *stage.DeployStagePollerFactory
	nameFinderFactory nameFinderFactory
}

func (op *deployStageOperation) Do(context operation.ExecutionContext) error {
	finder, err := op.nameFinderFactory(context)
	if err != nil {
		return err
	}

	azureDeploymentName, err := finder.Find(context.Context())
	if err != nil {
		return err
	}

	// save the deployment name to the operation so we can fetch it later
	context.Operation().Attribute(model.AttributeKeyAzureDeploymentName, azureDeploymentName)
	context.SaveChanges()

	isFirstAttempt := context.Operation().IsFirstAttempt()
	if isFirstAttempt {
		err := op.wait(context, azureDeploymentName)
		if err != nil {
			return err
		}
	} else { // retry the stage
		retryStage := NewRetryStageOperation()
		err := retryStage(context)
		if err != nil {
			return err
		}
	}
	return nil
}

func (op *deployStageOperation) wait(context operation.ExecutionContext, azureDeploymentName string) error {
	poller, err := op.pollerFactory.Create(context.Operation(), azureDeploymentName, nil)
	if err != nil {
		return err
	}
	response, err := poller.PollUntilDone(context.Context())
	if err != nil {
		return err
	}

	context.Value(response)

	if response.Status == sdk.StatusFailed {
		return operation.NewError(context.Operation())
	}

	return nil
}

func NewDeployStageOperation(appConfig *config.AppConfig) operation.OperationFunc {
	pollerFactory := stage.NewDeployStagePollerFactory()

	operation := &deployStageOperation{
		pollerFactory: pollerFactory,
		nameFinderFactory: func(context operation.ExecutionContext) (*operation.AzureDeploymentNameFinder, error) {
			return operation.NewAzureDeploymentNameFinder(context.Operation())
		},
	}
	return operation.Do
}
