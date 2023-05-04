package operations

import (
	"context"
	"fmt"
	"strconv"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	deployments "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operation"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// deployment retry
// TODO: implement retry of an entire deployment and also a stage

type retryDeployment struct {
	db *gorm.DB
}

func (exe *retryDeployment) Execute(ctx context.Context, invokedOperation *data.InvokedOperation) error {
	deployment, err := exe.updateToRunning(ctx, invokedOperation)
	if err != nil {
		return err
	}
	azureRedeployment := exe.mapToAzureRedeployment(deployment, invokedOperation)
	results, err := deployments.Redeploy(ctx, azureRedeployment)
	if err != nil {
		log.WithError(err).Error("error redeploying deployment")
	}
	err = exe.updateWithResults(ctx, results, invokedOperation)

	if err != nil {
		return err
	}

	return nil
}

func (exe *retryDeployment) updateWithResults(ctx context.Context, results *deployments.AzureDeploymentResult, invokedOperation *data.InvokedOperation) error {
	db := exe.db

	invokedOperation.Result = results

	if results.Status == deployments.Failed {
		invokedOperation.Status = operation.StatusFailed.String()
	} else if results.Status == deployments.Succeeded {
		invokedOperation.Status = operation.StatusSuccess.String()
	}

	db.Save(invokedOperation)

	message := &events.EventHookMessage{
		Status: invokedOperation.Status,
		Data: &events.DeploymentEventData{
			DeploymentId: int(invokedOperation.DeploymentId),
			Message:      fmt.Sprintf("Retry deployment completed %s", invokedOperation.Status),
		},
	}
	message.SetSubject(invokedOperation.DeploymentId, nil)

	err := hook.Add(ctx, message)
	if err != nil {
		return err
	}
	return nil
}

func (exe *retryDeployment) updateToRunning(ctx context.Context, invokedOperation *data.InvokedOperation) (*data.Deployment, error) {
	db := exe.db

	deployment := &data.Deployment{}
	db.First(&deployment, invokedOperation.DeploymentId)
	invokedOperation.Status = operation.StatusRunning.String()
	db.Save(invokedOperation)

	err := hook.Add(ctx, &events.EventHookMessage{
		Subject: "/deployments/" + strconv.Itoa(int(deployment.ID)),
		Status:  invokedOperation.Status,
		Data: &events.DeploymentEventData{
			DeploymentId: int(deployment.ID),
			Message:      "Retry deployment started successfully",
		},
	})
	if err != nil {
		return nil, err
	}
	return deployment, err
}

func (exe *retryDeployment) mapToAzureRedeployment(dep *data.Deployment, operation *data.InvokedOperation) deployments.AzureRedeployment {
	azureRedeployment := deployments.AzureRedeployment{
		SubscriptionId:    dep.SubscriptionId,
		Location:          dep.Location,
		ResourceGroupName: dep.ResourceGroup,
		DeploymentName:    dep.GetAzureDeploymentName(),
	}
	log.WithField("azureRedeployment", azureRedeployment).Debug("AzureRedeployment object")
	return azureRedeployment
}

//region factory

func NewRetryDeploymentExecutor(appConfig *config.AppConfig) Executor {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
	return &retryDeployment{
		db: db,
	}
}

//endregion factory
