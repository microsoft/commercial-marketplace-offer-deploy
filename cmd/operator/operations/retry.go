package operations

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/operation"
	deployments "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	log "github.com/sirupsen/logrus"
)

type retryOperation struct {
}

func (op *retryOperation) Do(context *operation.ExecutionContext) error {
	azureRedeployment := op.mapToAzureRedeployment(context)

	result, err := deployments.Redeploy(context.Context(), azureRedeployment)
	context.Value(result)

	if err != nil {
		return err
	}
	context.Value(result)
	return nil
}

func (op *retryOperation) mapToAzureRedeployment(context *operation.ExecutionContext) deployments.AzureRedeployment {
	dep := context.Operation().Deployment()

	azureRedeployment := deployments.AzureRedeployment{
		SubscriptionId:    dep.SubscriptionId,
		Location:          dep.Location,
		ResourceGroupName: dep.ResourceGroup,
		DeploymentName:    dep.GetAzureDeploymentName(),
	}
	log.WithField("azureRedeployment", azureRedeployment).Debug("AzureRedeployment object")
	return azureRedeployment
}

func NewRetryDeploymentExecutor(appConfig *config.AppConfig) operation.OperationFunc {
	operation := &retryOperation{}
	return operation.Do
}
