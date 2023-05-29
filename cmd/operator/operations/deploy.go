package operations

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
)

type deployeOperation struct {
	retryOperation        operation.OperationFunc
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
		do = exe.retryOperation
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

func NewDeploymentOperation() operation.OperationFunc {
	operation := &deployeOperation{
		retryOperation:        NewRetryOperation(),
		createAzureDeployment: deployment.Create,
	}
	return operation.do
}
