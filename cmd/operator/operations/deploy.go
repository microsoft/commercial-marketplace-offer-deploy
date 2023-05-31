package operations

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
)

type deployeOperation struct {
	retryOperation operation.OperationFunc
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
	azureDeployment := exe.mapAzureDeployment(context.Operation())
	deployer, err := exe.newDeployer(azureDeployment.SubscriptionId)
	if err != nil {
		return err
	}

	token, err := deployer.Begin(context.Context(), azureDeployment)
	if err != nil {
		return err
	}

	context.Attribute(model.AttributeKeyResumeToken, token)

	result, err := deployer.Wait(context.Context(), token)
	context.Value(result)

	if err != nil {
		return err
	}

	return nil
}

func (p *deployeOperation) newDeployer(subscriptionId string) (deployment.Deployer, error) {
	return deployment.NewDeployer(deployment.DeploymentTypeARM, subscriptionId)
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
		retryOperation: NewRetryOperation(),
	}
	return operation.do
}
