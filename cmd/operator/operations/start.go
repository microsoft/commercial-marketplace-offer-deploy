package operations

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"

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
	db *gorm.DB
}

func (exe *startDeployment) Execute(ctx context.Context, operation *data.InvokedOperation) error {
	deployment, err := exe.updateToRunning(ctx, operation)
	if err != nil {
		return err
	}

	azureDeployment := exe.mapAzureDeployment(deployment, operation)
	result, err := exe.deploy(ctx, azureDeployment)

	// if we waited this long, then we can assume we have the results, so we'll update the invoked operation results with it
	if err == nil {
		operation.Result = result
		exe.db.Save(operation)
	}
	return err
}

func (exe *startDeployment) updateToRunning(ctx context.Context, invokedOperation *data.InvokedOperation) (*data.Deployment, error) {
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
			Message:      "Deployment started successfully",
		},
	})
	if err != nil {
		return nil, err
	}
	return deployment, err
}

func (exe *startDeployment) updateToFailed(ctx context.Context, invokedOperation *data.InvokedOperation, err error) error {
	db := exe.db

	deployment := &data.Deployment{}
	db.First(&deployment, invokedOperation.DeploymentId)
	invokedOperation.Status = string(operation.StatusFailed)
	db.Save(invokedOperation)

	err = hook.Add(ctx, &events.EventHookMessage{
		Subject: "/deployments/" + strconv.Itoa(int(deployment.ID)),
		Status:  operation.StatusFailed.String(),
		Type:    string(events.EventTypeDeploymentCompleted),
		Data: &events.DeploymentEventData{
			DeploymentId: int(deployment.ID),
			OperationId:  to.Ptr(invokedOperation.ID.String()),
			Message:      fmt.Sprintf("Azure Deployment failed. Result: %s", err.Error()),
		},
	})
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
	executor := &startDeployment{
		db: db,
	}
	return executor
}

//endregion factory
