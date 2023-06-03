package operations

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	log "github.com/sirupsen/logrus"
)

type retryStageOperation struct{}

func (op *retryStageOperation) Do(context *operation.ExecutionContext) error {
	stageId, err := op.getStageId(context)
	if err != nil {
		return err
	}

	parent := context.Operation().Deployment()
	stage := op.findStage(parent, stageId)

	if stage == nil {
		return fmt.Errorf("stage not found for deployment %v and stageId %v", parent.ID, stageId)
	}

	redeployment := op.mapToAzureRedeployment(parent, stage)
	deployer, err := op.newDeployer(redeployment.SubscriptionId)
	if err != nil {
		return err
	}

	result, err := deployer.Redeploy(context.Context(), redeployment)
	if err != nil {
		return err
	}

	context.Value(result)
	return nil
}

func (op *retryStageOperation) getStageId(context *operation.ExecutionContext) (uuid.UUID, error) {
	stageId, err := uuid.Parse(context.Operation().Parameters["stageId"].(string))

	if err != nil {
		log.Errorf("error parsing stageId: %s", err)
		return uuid.Nil, err
	}
	return stageId, nil
}

func (op *retryStageOperation) newDeployer(subscriptionId string) (deployment.Deployer, error) {
	return deployment.NewDeployer(deployment.DeploymentTypeARM, subscriptionId)
}

func (op *retryStageOperation) mapToAzureRedeployment(dep *model.Deployment, stage *model.Stage) deployment.AzureRedeployment {
	azureRedeployment := deployment.AzureRedeployment{
		SubscriptionId:    dep.SubscriptionId,
		Location:          dep.Location,
		ResourceGroupName: dep.ResourceGroup,
		DeploymentName:    stage.AzureDeploymentName,
	}
	return azureRedeployment
}

func (op *retryStageOperation) findStage(deployment *model.Deployment, stageId uuid.UUID) *model.Stage {
	for _, stage := range deployment.Stages {
		if stage.ID == stageId {
			return &stage
		}
	}
	return nil
}

func NewRetryStageOperation() operation.OperationFunc {
	operation := &retryStageOperation{}
	return operation.Do
}
