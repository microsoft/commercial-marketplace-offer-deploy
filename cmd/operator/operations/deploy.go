package operations

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"gorm.io/gorm"
)

// RetriableError is a custom error that contains a positive duration for the next retry
type RetriableError struct {
	Err        error
	RetryAfter time.Duration
}

// Error returns error message and a Retry-After duration
func (e *RetriableError) Error() string {
	return fmt.Sprintf("%s (retry after %v)", e.Err.Error(), e.RetryAfter)
}

type startDeployment struct {
	db                    *gorm.DB
	factory               ExecutorFactory
	createAzureDeployment deployment.CreateDeployment
}

func (exe *startDeployment) Execute(ctx context.Context, invokedOperation *model.InvokedOperation) error {
	invokedOperation.Attempts = invokedOperation.Attempts + 1

	if invokedOperation.Attempts > invokedOperation.Retries {
		return nil
	}

	execute := exe.start

	if invokedOperation.Attempts > 1 { // this is a retry if so
		executor, err := exe.factory.Create(sdk.OperationRetryDeployment)
		if err != nil {
			exe.updateToFailed(ctx, invokedOperation, err)
			return err
		}
		execute = executor.Execute
	}

	err := execute(ctx, invokedOperation)
	if err != nil {
		exe.updateToFailed(ctx, invokedOperation, err)
	}
	return nil
}

func (exe *startDeployment) start(ctx context.Context, invokedOperation *model.InvokedOperation) error {
	err := exe.updateToRunning(ctx, invokedOperation)
	if err != nil {
		return err
	}

	deployment := &model.Deployment{}
	exe.db.First(deployment, invokedOperation.DeploymentId)

	azureDeployment := exe.mapAzureDeployment(deployment, invokedOperation)
	result, err := exe.createAzureDeployment(ctx, azureDeployment)

	// if we waited this long, then we can assume we have the results, so we'll update the invoked operation results with it
	if err == nil {
		invokedOperation.Result = result
		exe.db.Save(invokedOperation)
	}
	return err
}

func (exe *startDeployment) updateToRunning(ctx context.Context, invokedOperation *model.InvokedOperation) error {
	db := exe.db

	invokedOperation.Status = sdk.StatusRunning.String()
	db.Save(&invokedOperation)

	message := "Deployment started successfully"
	if invokedOperation.Attempts > 1 {
		message = "Deployment retry started successfully"
	}
	err := hook.Add(ctx, &sdk.EventHookMessage{
		Subject: "/deployments/" + strconv.Itoa(int(invokedOperation.DeploymentId)),
		Status:  invokedOperation.Status,
		Type:    string(sdk.EventTypeDeploymentStarted),
		Data: &sdk.DeploymentEventData{
			DeploymentId: int(invokedOperation.DeploymentId),
			OperationId:  invokedOperation.ID,
			Message:      message,
		},
	})
	return err
}

func (exe *startDeployment) updateToFailed(ctx context.Context, invokedOperation *model.InvokedOperation, err error) error {
	db := exe.db

	invokedOperation.Status = string(sdk.StatusFailed)
	db.Save(&invokedOperation)

	eventType := string(sdk.EventTypeDeploymentCompleted)
	if invokedOperation.Attempts > 1 {
		eventType = string(sdk.EventTypeDeploymentRetried)
	}

	message := &sdk.EventHookMessage{
		Subject: "/deployments/" + strconv.Itoa(int(invokedOperation.DeploymentId)),
		Status:  sdk.StatusFailed.String(),
		Type:    eventType,
		Data: &sdk.DeploymentEventData{
			DeploymentId: int(invokedOperation.DeploymentId),
			OperationId:  invokedOperation.ID,
			Attempts:     invokedOperation.Attempts,
			Message:      fmt.Sprintf("Azure Deployment failed on the attempt %d. Error: %s", invokedOperation.Attempts, err.Error()),
		},
	}

	err = hook.Add(ctx, message)

	if err != nil {
		return err
	}
	return nil
}

func (p *startDeployment) mapAzureDeployment(d *model.Deployment, io *model.InvokedOperation) deployment.AzureDeployment {
	return deployment.AzureDeployment{
		SubscriptionId:    d.SubscriptionId,
		ResourceGroupName: d.ResourceGroup,
		DeploymentName:    d.GetAzureDeploymentName(),
		Template:          d.Template,
		Params:            io.Parameters,
	}
}

//region factory

func NewStartDeploymentExecutor(appConfig *config.AppConfig) Executor {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
	factory := NewExecutorFactory(appConfig)
	executor := &startDeployment{
		db:                    db,
		factory:               factory,
		createAzureDeployment: deployment.Create,
	}
	return executor
}

//endregion factory
