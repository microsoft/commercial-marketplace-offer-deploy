package operations

import (
	"context"
	"errors"
	"fmt"
	"strings"
	//	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	//"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	deployment "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	//	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"gorm.io/gorm"
)

type retryStage struct {
	db *gorm.DB
}     

func (exe *retryStage) Execute(ctx context.Context, operation *data.InvokedOperation) error {
	db := exe.db

	dep := &data.Deployment{}
	db.First(&dep, operation.DeploymentId)

	stageName := operation.Parameters["stageName"].(string)
	foundStage := exe.findStage(dep, stageName)

	if foundStage == nil {
		errMsg := fmt.Sprintf("stage not found for deployment %v and stageId %v", dep.ID, stageName)
		return errors.New(errMsg)
	}

	redeployment := exe.mapToAzureRedeployment(dep, foundStage, operation)
	if redeployment == nil {
		return errors.New("unable to map to AzureRedeployment")
	}

	res, err := deployment.Redeploy(*redeployment)
	if err != nil {
		return err
	}

	if res == nil {
		log.Debugf("Redeployment response is nil")
	}

	// err := exe.updateToRunning(ctx, operation, dep)
	// if err != nil {
	// 	log.Errorf("error updating deployment to running: %s", err)
	// 	return err
	// }
	return nil
}

func (exe *retryStage) mapToAzureRedeployment(dep *data.Deployment, stage *data.Stage, operation *data.InvokedOperation) *deployment.AzureRedeployment {
	azureRedeployment := &deployment.AzureRedeployment {
		SubscriptionId: dep.SubscriptionId,
		Location: dep.Location,
		ResourceGroupName: dep.ResourceGroup,
		DeploymentName: stage.DeploymentName,
	}

	return azureRedeployment
}

func (exe *retryStage) findStage(deployment *data.Deployment, stageName string) *data.Stage {
	for _, stage := range deployment.Stages {
		if strings.EqualFold(stage.Name, stageName) {
			return &stage
		}
	}
	return nil
}

func (exe *retryStage) updateToRunning(ctx context.Context, operation *data.InvokedOperation, deployment *data.Deployment) error {
	// db := exe.db

	// operation.Status = events.StatusRunning.String()
	// db.Save(operation)

	// subject := fmt.Sprintf("/deployments")
	// err := hook.Add(&events.EventHookMessage{
	// 	Subject: "/deployments/" + strconv.Itoa(int(deployment.ID)),
	// 	Status:  deployment.Status,
	// 	Data: &events.DeploymentEventData{
	// 		DeploymentId: int(deployment.ID),
	// 		Message:      "Deployment started successfully",
	// 	},
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// return deployment, err
	return nil
}

//region factory

func NewRetryStageExecutor(appConfig *config.AppConfig) Executor {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
	return &retryStage{
		db: db,
	}
}

//endregion factory
