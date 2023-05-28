package operations

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	log "github.com/sirupsen/logrus"
)

type retryStageOperation struct{}

func (op *retryStageOperation) Do(context *operation.ExecutionContext) error {
	stageId, err := op.getStageId(context)
	if err != nil {
		return err
	}

	parent := context.InvokedOperation().Deployment()
	stage := op.findStage(parent, stageId)

	if stage == nil {
		return fmt.Errorf("stage not found for deployment %v and stageId %v", parent.ID, stageId)
	}

	redeployment := op.mapToAzureRedeployment(parent, stage)
	if redeployment == nil {
		return errors.New("unable to map to AzureRedeployment")
	}

	result, err := deployment.Redeploy(context.Context(), *redeployment)
	if err != nil {
		return err
	}

	context.Value(result)
	return nil
}

func (op *retryStageOperation) getStageId(context *operation.ExecutionContext) (uuid.UUID, error) {
	stageId, err := uuid.Parse(context.InvokedOperation().Parameters["stageId"].(string))

	if err != nil {
		log.Errorf("error parsing stageId: %s", err)
		return uuid.Nil, err
	}
	return stageId, nil
}

func (exe *retryStageOperation) mapToAzureRedeployment(dep *model.Deployment, stage *model.Stage) *deployment.AzureRedeployment {
	azureRedeployment := &deployment.AzureRedeployment{
		SubscriptionId:    dep.SubscriptionId,
		Location:          dep.Location,
		ResourceGroupName: dep.ResourceGroup,
		DeploymentName:    stage.DeploymentName,
	}
	return azureRedeployment
}

func (exe *retryStageOperation) findStage(deployment *model.Deployment, stageId uuid.UUID) *model.Stage {
	for _, stage := range deployment.Stages {
		if stage.ID == stageId {
			return &stage
		}
	}
	return nil
}

func NewRetryStageExecutor(appConfig *config.AppConfig) operation.OperationFunc {
	operation := &retryStageOperation{}
	return operation.Do
}
