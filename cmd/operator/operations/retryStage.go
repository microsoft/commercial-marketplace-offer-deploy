package operations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type retryStage struct {
	db *gorm.DB
}

func (exe *retryStage) Execute(ctx context.Context, invokedOperation *data.InvokedOperation) error {
	db := exe.db

	dep := &data.Deployment{}
	db.First(&dep, invokedOperation.DeploymentId)

	db.Save(dep)
	invokedOperation.Status = string(sdk.StatusRunning)
	exe.save(invokedOperation)

	stageId, err := uuid.Parse(invokedOperation.Parameters["stageId"].(string))
	if err != nil {
		log.Errorf("error parsing stageId: %s", err)
		return err
	}

	foundStage := exe.findStage(dep, stageId)
	if foundStage == nil {
		errMsg := fmt.Sprintf("stage not found for deployment %v and stageId %v", dep.ID, stageId)
		return errors.New(errMsg)
	}

	redeployment := exe.mapToAzureRedeployment(dep, foundStage, invokedOperation)
	if redeployment == nil {
		return errors.New("unable to map to AzureRedeployment")
	}

	res, err := deployment.Redeploy(ctx, *redeployment)
	if err != nil {
		log.Errorf("error redeploying deployment: %s", err)
		return err
	}

	if res == nil {
		log.Debugf("Redeployment response is nil")
	}

	err = exe.sendHook(ctx, dep, invokedOperation, foundStage)
	if err != nil {
		log.Errorf("error adding hook: %s", err)
		return err
	}

	return nil
}

func (exe *retryStage) sendHook(ctx context.Context, deployment *data.Deployment, invokedOperation *data.InvokedOperation, stage *data.Stage) error {
	message := &sdk.EventHookMessage{
		Status: invokedOperation.Status,
		Data: &sdk.DeploymentEventData{
			DeploymentId: int(deployment.ID),
			StageId:      to.Ptr(stage.ID),
			Message:      "Retry stage started successfully",
		},
	}
	message.SetSubject(deployment.ID, &stage.ID)
	err := hook.Add(ctx, message)

	if err != nil {
		return err
	}

	return nil
}

func (exe *retryStage) save(operation *data.InvokedOperation) error {
	tx := exe.db.Begin()
	tx.Save(operation)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()

	return nil
}

func (exe *retryStage) mapToAzureRedeployment(dep *data.Deployment, stage *data.Stage, operation *data.InvokedOperation) *deployment.AzureRedeployment {
	b, err := json.MarshalIndent(dep, "", "  ")
	if err != nil {
		log.Error(err)
	}
	log.Debugf("retryStage.mapToAzureRedeployment: dep: %v, stage: %v, operation: %v", string(b), stage, operation)
	azureRedeployment := &deployment.AzureRedeployment{
		SubscriptionId:    dep.SubscriptionId,
		Location:          dep.Location,
		ResourceGroupName: dep.ResourceGroup,
		DeploymentName:    stage.DeploymentName,
	}

	return azureRedeployment
}

func (exe *retryStage) findStage(deployment *data.Deployment, stageId uuid.UUID) *data.Stage {
	for _, stage := range deployment.Stages {
		if stage.ID == stageId {
			return &stage
		}
	}
	return nil
}

func NewRetryStageExecutor(appConfig *config.AppConfig) Executor {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
	return &retryStage{
		db: db,
	}
}
