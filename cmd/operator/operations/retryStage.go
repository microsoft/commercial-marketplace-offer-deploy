package operations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operation"
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
	invokedOperation.Status = string(operation.StatusRunning)
	exe.save(invokedOperation)

	stageId := uuid.MustParse(invokedOperation.Parameters["stageId"].(string))
	foundStage := exe.findStage(dep, stageId)
	if foundStage == nil {
		errMsg := fmt.Sprintf("stage not found for deployment %v and stageId %v", dep.ID, stageId)
		return errors.New(errMsg)
	}

	redeployment := exe.mapToAzureRedeployment(dep, foundStage, invokedOperation)
	if redeployment == nil {
		return errors.New("unable to map to AzureRedeployment")
	}

	res, err := deployment.Redeploy(*redeployment)
	if err != nil {
		log.Errorf("error redeploying deployment: %s", err)
		return err
	}

	if res == nil {
		log.Debugf("Redeployment response is nil")
	}

	err = exe.sendHook(dep, foundStage)
	if err != nil {
		log.Errorf("error adding hook: %s", err)
		return err
	}

	return nil
}

func (exe *retryStage) sendHook(deployment *data.Deployment, stage *data.Stage) error {
	subject := fmt.Sprintf("/deployments/%v/%v", strconv.Itoa(int(deployment.ID)), stage.Name)
	log.Debugf("retryStage.updateToRunning: subject: %s", subject)

	stageName := stage.ID.String()
	log.Debugf("retryStage.updateToRunning: stageName: %s", stageName)

	err := hook.Add(&events.EventHookMessage{
		Subject: subject,
		Status:  deployment.Status,
		Data: &events.DeploymentEventData{
			DeploymentId: int(deployment.ID),
			StageId:      &stageName,
			Message:      "Deployment started successfully",
		},
	})

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

// func (exe *retryStage) updateToRunning(ctx context.Context, operation *data.InvokedOperation, deployment *data.Deployment, stage *data.Stage) error {
// 	log.Debugf("retryStage.updateToRunning: operation: %v, deployment: %v, stage: %v", operation, deployment, stage)
// 	db := exe.db

// 	if operation == nil {
// 		log.Error("retryStage.updateToRunning: operation is nil")
// 		return errors.New("operation is nil")
// 	}

// 	operation.Status = events.StatusRunning.String()
// 	db.Save(operation)
// 	log.Debug("retryStage.updateToRunning: operation saved")

// 	subject := fmt.Sprintf("/deployments/%v/%v", strconv.Itoa(int(deployment.ID)), stage.Name)
// 	log.Debugf("retryStage.updateToRunning: subject: %s", subject)

// 	stageName := stage.ID.String()
// 	log.Debugf("retryStage.updateToRunning: stageName: %s", stageName)

// 	err := hook.Add(&events.EventHookMessage{
// 		Subject: subject,
// 		Status:  deployment.Status,
// 		Data: &events.DeploymentEventData{
// 			DeploymentId: int(deployment.ID),
// 			StageId: &stageName,
// 			Message:      "Deployment started successfully",
// 		},
// 	})

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

//region factory

func NewRetryStageExecutor(appConfig *config.AppConfig) Executor {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
	return &retryStage{
		db: db,
	}
}

//endregion factory
