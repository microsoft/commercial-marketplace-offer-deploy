package operations

import (
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
)

type deployeOperation struct {
	retryOperation operation.OperationFunc
}

// the operation to execute
func (op *deployeOperation) Do(context *operation.ExecutionContext) error {
	operation, err := op.getOperation(context)
	if err != nil {
		return err
	}
	return operation(context)
}

func (op *deployeOperation) getOperation(context *operation.ExecutionContext) (operation.OperationFunc, error) {
	do := op.do

	if context.Operation().IsRetry() { // this is a retry if so
		do = op.retryOperation
	}
	return do, nil
}

func (op *deployeOperation) do(context *operation.ExecutionContext) error {
	azureDeployment := op.mapAzureDeployment(context.Operation())
	deployer, err := op.newDeployer(azureDeployment.SubscriptionId)
	if err != nil {
		return err
	}

	beginResult, err := deployer.Begin(context.Context(), azureDeployment)
	if err != nil {
		return err
	}

	token := beginResult.ResumeToken

	context.Attribute(model.AttributeKeyResumeToken, token)
	context.Attribute(model.AttributeKeyCorrelationId, beginResult.CorrelationID)

	result, err := deployer.Wait(context.Context(), &token)
	context.Value(result)

	if err != nil {
		return err
	}

	return nil
}

func (op *deployeOperation) newDeployer(subscriptionId string) (deployment.Deployer, error) {
	return deployment.NewDeployer(deployment.DeploymentTypeARM, subscriptionId)
}

func (op *deployeOperation) mapAzureDeployment(invokedOperation *operation.Operation) deployment.AzureDeployment {
	d := invokedOperation.Deployment()

	return deployment.AzureDeployment{
		SubscriptionId:    d.SubscriptionId,
		ResourceGroupName: d.ResourceGroup,
		DeploymentName:    d.GetAzureDeploymentName(),
		Template:          d.Template,
		Params:            invokedOperation.Parameters,
		Tags: map[string]*string{
			string(deployment.LookupTagKeyDeploymentId): to.Ptr(strconv.Itoa(int(invokedOperation.DeploymentId))),
			string(deployment.LookupTagKeyOperationId):  to.Ptr(invokedOperation.ID.String()),
		},
	}
}

func NewDeploymentOperation() operation.OperationFunc {
	operation := &deployeOperation{
		retryOperation: NewRetryOperation(),
	}
	return operation.do
}
