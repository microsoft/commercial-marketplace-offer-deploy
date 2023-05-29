package operations

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

type deployeOperation struct {
	factory               operation.ExecutorFactory
	createAzureDeployment deployment.CreateDeployment
}

// the operation to execute
func (exe *deployeOperation) Do(context *operation.ExecutionContext) error {
	operation, err := exe.getOperation(context)
	if err != nil {
		return err
	}
	return operation(context)
}

func (exe *deployeOperation) getOperation(context *operation.ExecutionContext) (operation.OperationFunc, error) {
	do := exe.do

	if context.Operation().IsRetry() { // this is a retry if so
		executor, err := exe.factory.Create(sdk.OperationRetry)
		if err != nil {
			return nil, err
		}
		do = executor.Execute
	}
	return do, nil
}

func (exe *deployeOperation) do(context *operation.ExecutionContext) error {
	invokedOperation := context.Operation()

	azureDeployment := exe.mapAzureDeployment(invokedOperation)
	result, err := exe.createAzureDeployment(context.Context(), azureDeployment)
	invokedOperation.Value(result)

	if err != nil {
		return err
	}

	return nil
}

func (p *deployeOperation) mapAzureDeployment(invokedOperation *operation.Operation) deployment.AzureDeployment {
	d := invokedOperation.Deployment()

	return deployment.AzureDeployment{
		SubscriptionId:    d.SubscriptionId,
		ResourceGroupName: d.ResourceGroup,
		DeploymentName:    d.GetAzureDeploymentName(),
		Template:          d.Template,
		Params:            invokedOperation.Parameters,
	}
}

//region factory

func NewDeploymentOperation(appConfig *config.AppConfig) operation.OperationFunc {
	factory := NewExecutorFactory(appConfig)

	operation := &deployeOperation{
		factory:               factory,
		createAzureDeployment: deployment.Create,
	}
	return operation.do
}

//endregion factory
