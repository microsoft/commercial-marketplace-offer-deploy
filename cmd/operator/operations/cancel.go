package operations

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
)

type cancelOperation struct {
	cancel operation.OperationFunc
}

func (op *cancelOperation) Do(context operation.ExecutionContext) error {
	return op.do(context)
}

func (op *cancelOperation) do(context operation.ExecutionContext) error {
	azureCancelDeployment := op.mapToAzureCancel(context.Operation())
	ctx := context.Context()

	deployer, err := op.newDeployer(azureCancelDeployment.SubscriptionId)
	if err != nil {
		return err
	}

	result, err := deployer.Cancel(ctx, azureCancelDeployment)
	context.Value(result)

	if err != nil {
		return err
	}

	return nil
}

func (op *cancelOperation) newDeployer(subscriptionId string) (deployment.Deployer, error) {
	return deployment.NewDeployer(deployment.DeploymentTypeARM, subscriptionId)
}

func (op *cancelOperation) mapToAzureCancel(invokedOperation *operation.Operation) deployment.AzureCancelDeployment {
	d := invokedOperation.Deployment()
	return deployment.AzureCancelDeployment{
		SubscriptionId:    d.SubscriptionId,
		ResourceGroupName: d.ResourceGroup,
		DeploymentName:    d.GetAzureDeploymentName(),
	}
}

func NewCancelOperation() operation.OperationFunc {
	operation := &cancelOperation{}
	return operation.Do
}
