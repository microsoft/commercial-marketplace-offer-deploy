package operations

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
)

type retryStageOperation struct{}

func (op *retryStageOperation) Do(context operation.ExecutionContext) error {
	redeployment, err := op.mapToAzureRedeployment(context.Operation())
	if err != nil {
		return err
	}

	deployer, err := op.newDeployer(redeployment.SubscriptionId)
	if err != nil {
		return err
	}

	result, err := deployer.Redeploy(context.Context(), redeployment)

	context.Value(result)

	if err != nil {
		return err
	}

	return nil
}

func (op *retryStageOperation) newDeployer(subscriptionId string) (deployment.Deployer, error) {
	return deployment.NewDeployer(deployment.DeploymentTypeARM, subscriptionId)
}

func (op *retryStageOperation) mapToAzureRedeployment(operation *operation.Operation) (deployment.AzureRedeployment, error) {
	dep := operation.Deployment()

	azureDeploymentName, ok := operation.AttributeValue(model.AttributeKeyAzureDeploymentName)
	if !ok {
		return deployment.AzureRedeployment{}, fmt.Errorf("error getting azure deployment name from operation %v", operation.ID)
	}

	azureRedeployment := deployment.AzureRedeployment{
		SubscriptionId:    dep.SubscriptionId,
		Location:          dep.Location,
		ResourceGroupName: dep.ResourceGroup,
		DeploymentName:    azureDeploymentName.(string),
		OperationId:       operation.ID,
		Tags: map[string]*string{
			string(deployment.LookupTagKeyOperationId): to.Ptr(operation.ID.String()),
		},
	}
	return azureRedeployment, nil
}

func NewRetryStageOperation() operation.OperationFunc {
	operation := &retryStageOperation{}
	return operation.Do
}
