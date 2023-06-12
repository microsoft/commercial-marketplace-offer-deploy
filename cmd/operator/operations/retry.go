package operations

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	log "github.com/sirupsen/logrus"
)

type retryTask struct {
}

func (op *retryTask) Run(context operation.ExecutionContext) error {
	azureRedeployment := op.mapToAzureRedeployment(context)
	deployer, err := op.newDeployer(azureRedeployment.SubscriptionId)
	if err != nil {
		return err
	}
	result, err := deployer.Redeploy(context.Context(), azureRedeployment)
	context.Value(result)

	if err != nil {
		return err
	}
	return nil
}

func (task *retryTask) Continue(context operation.ExecutionContext) error {
	return nil
}

func (op *retryTask) newDeployer(subscriptionId string) (deployment.Deployer, error) {
	return deployment.NewDeployer(deployment.DeploymentTypeARM, subscriptionId)
}

func (op *retryTask) mapToAzureRedeployment(context operation.ExecutionContext) deployment.AzureRedeployment {
	dep := context.Operation().Deployment()

	azureRedeployment := deployment.AzureRedeployment{
		SubscriptionId:    dep.SubscriptionId,
		Location:          dep.Location,
		ResourceGroupName: dep.ResourceGroup,
		DeploymentName:    dep.GetAzureDeploymentName(),
	}
	log.WithField("azureRedeployment", azureRedeployment).Debug("AzureRedeployment object")
	return azureRedeployment
}

func NewRetryTask() operation.OperationTask {
	return &retryTask{}
}
