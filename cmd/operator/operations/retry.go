package operations

import (
	"context"
	"fmt"
	"strconv"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	deployments "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// deployment retry
// TODO: implement retry of an entire deployment and also a stage

type retryDeployment struct {
	db *gorm.DB
}

func (exe *retryDeployment) Execute(ctx context.Context, invokedOperation *model.InvokedOperation) error {
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

func (exe *retryDeployment) updateWithResults(ctx context.Context, results *deployments.AzureDeploymentResult, invokedOperation *model.InvokedOperation) error {
	db := exe.db

	if results != nil {
		invokedOperation.Result = results

		if results.Status == deployments.Failed {
			invokedOperation.Status = sdk.StatusFailed.String()
		} else if results.Status == deployments.Succeeded {
			invokedOperation.Status = sdk.StatusSuccess.String()
		}
	} else {
		invokedOperation.Status = sdk.StatusFailed.String()
	}

	db.Save(invokedOperation)

	message := &sdk.EventHookMessage{
		Status: invokedOperation.Status,
		Type:   string(sdk.EventTypeDeploymentRetried),
		Data: &sdk.DeploymentEventData{
			DeploymentId: int(invokedOperation.DeploymentId),
			OperationId:  invokedOperation.ID,
			Message:      fmt.Sprintf("Retry deployment %s", invokedOperation.Status),
		},
	}
	message.SetSubject(invokedOperation.DeploymentId, nil)

	// by sending this message, it will be caught (in operator/handlers/sdk.go) and the retry will get executed again
	// as long as the Attempts haven't exceeded the max set on Retries
	err := hook.Add(ctx, message)
	if err != nil {
		return err
	}
	return nil
}

func (exe *retryDeployment) updateToRunning(ctx context.Context, invokedOperation *model.InvokedOperation) (*model.Deployment, error) {
	db := exe.db

	deployment := &model.Deployment{}
	db.First(&deployment, invokedOperation.DeploymentId)
	invokedOperation.Status = sdk.StatusRunning.String()
	db.Save(invokedOperation)

	err := hook.Add(ctx, &sdk.EventHookMessage{
		Subject: "/deployments/" + strconv.Itoa(int(deployment.ID)),
		Status:  invokedOperation.Status,
		Data: &sdk.DeploymentEventData{
			DeploymentId: int(deployment.ID),
			Message:      "Retry deployment started successfully",
		},
	})
	if err != nil {
		return nil, err
	}
	return deployment, err
}

func (exe *retryDeployment) mapToAzureRedeployment(dep *model.Deployment, operation *model.InvokedOperation) deployments.AzureRedeployment {
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
