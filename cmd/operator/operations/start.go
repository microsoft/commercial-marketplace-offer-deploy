package operations

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operation"
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
	db      *gorm.DB
	factory ExecutorFactory
}

func (exe *startDeployment) Execute(ctx context.Context, invokedOperation *data.InvokedOperation) error {
	invokedOperation.Attempts = invokedOperation.Attempts + 1

	if invokedOperation.Attempts > invokedOperation.Retries {
		return nil
	}

	execute := exe.start

	if invokedOperation.Attempts > 1 { // this is a retry if so
		executor, err := exe.factory.Create(operation.TypeRetryDeployment)
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

func (exe *startDeployment) start(ctx context.Context, invokedOperation *data.InvokedOperation) error {
	deployment, err := exe.updateToRunning(ctx, invokedOperation)
	if err != nil {
		return err
	}

	azureDeployment := exe.mapAzureDeployment(deployment, invokedOperation)
	result, err := exe.deploy(ctx, azureDeployment)

	// if we waited this long, then we can assume we have the results, so we'll update the invoked operation results with it
	if err == nil {
		invokedOperation.Result = result
		exe.db.Save(&invokedOperation)
	}
	return err
}

func (exe *startDeployment) updateToRunning(ctx context.Context, invokedOperation *data.InvokedOperation) (*data.Deployment, error) {
	db := exe.db

	deployment := &data.Deployment{}
	db.First(&deployment, invokedOperation.DeploymentId)
	invokedOperation.Status = operation.StatusRunning.String()
	db.Save(&invokedOperation)

	message := "Deployment started successfully"
	if invokedOperation.Attempts > 1 {
		message = "Deployment retry started successfully"
	}
	err := hook.Add(ctx, &events.EventHookMessage{
		Subject: "/deployments/" + strconv.Itoa(int(deployment.ID)),
		Status:  invokedOperation.Status,
		Data: &events.DeploymentEventData{
			DeploymentId: int(deployment.ID),
			OperationId:  invokedOperation.ID,
			Message:      message,
		},
	})
	if err != nil {
		return nil, err
	}
	return deployment, err
}

func (exe *startDeployment) updateToFailed(ctx context.Context, invokedOperation *data.InvokedOperation, err error) error {
	db := exe.db

	invokedOperation.Status = string(operation.StatusFailed)
	db.Save(&invokedOperation)

	eventType := string(events.EventTypeDeploymentCompleted)
	if invokedOperation.Attempts > 1 {
		eventType = string(events.EventTypeDeploymentRetried)
	}

	message := &events.EventHookMessage{
		Subject: "/deployments/" + strconv.Itoa(int(invokedOperation.DeploymentId)),
		Status:  operation.StatusFailed.String(),
		Type:    eventType,
		Data: &events.DeploymentEventData{
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

func (p *startDeployment) mapAzureDeployment(d *data.Deployment, io *data.InvokedOperation) *deployment.AzureDeployment {
	return &deployment.AzureDeployment{
		SubscriptionId:    d.SubscriptionId,
		ResourceGroupName: d.ResourceGroup,
		DeploymentName:    d.GetAzureDeploymentName(),
		Template:          d.Template,
		Params:            io.Parameters,
	}
}

func (p *startDeployment) deploy(ctx context.Context, azureDeployment *deployment.AzureDeployment) (*deployment.AzureDeploymentResult, error) {
	return deployment.Create(ctx, *azureDeployment)
}

//region factory

func NewStartDeploymentExecutor(appConfig *config.AppConfig) Executor {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
	factory := NewExecutorFactory(appConfig)
	executor := &startDeployment{
		db:      db,
		factory: factory,
	}
	return executor
}

//endregion factory
