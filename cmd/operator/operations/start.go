package operations

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/avast/retry-go"
	log "github.com/sirupsen/logrus"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
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

const defaultRetryCount = 3

type startDeployment struct {
	db *gorm.DB
}

func (exe *startDeployment) Execute(ctx context.Context, operation *data.InvokedOperation) error {
	deployment, err := exe.updateToRunning(ctx, operation)
	if err != nil {
		return err
	}

	// TODO: need to update the invoked operation state

	azureDeployment := exe.mapAzureDeployment(deployment, operation)
	exe.executeAsync(ctx, operation, azureDeployment)
	return nil
}

func (exe *startDeployment) executeAsync(ctx context.Context, operation *data.InvokedOperation, azureDeployment *deployment.AzureDeployment) {
	go func() {
		retry.Do(func() error {
			_, err := exe.deploy(ctx, azureDeployment)
			if err != nil {
				log.Error("Error calling deployment.Create: ", err)
				err = exe.updateToFailed(ctx, operation, err)

				if err != nil {
					log.Error("Error updating Deployment to failed: ", err)
				}

			}
			return &RetriableError{Err: err, RetryAfter: 10 * time.Second}
		},
			retry.Attempts(defaultRetryCount),
			retry.DelayType(backOffRetryHandler))
	}()
}

func (exe *startDeployment) updateToRunning(ctx context.Context, operation *data.InvokedOperation) (*data.Deployment, error) {
	db := exe.db

	deployment := &data.Deployment{}
	db.First(&deployment, operation.DeploymentId)
	deployment.Status = events.StatusRunning.String()
	db.Save(deployment)

	err := hook.Add(&events.EventHookMessage{
		Subject: "/deployments/" + strconv.Itoa(int(deployment.ID)),
		Status:  deployment.Status,
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

func (exe *startDeployment) updateToFailed(ctx context.Context, operation *data.InvokedOperation, err error) error {
	db := exe.db

	deployment := &data.Deployment{}
	db.First(&deployment, operation.DeploymentId)
	deployment.Status = string(events.StatusFailed)
	db.Save(deployment)

	err = hook.Add(&events.EventHookMessage{
		Subject: "/deployments/" + strconv.Itoa(int(deployment.ID)),
		Status:  events.StatusFailed.String(),
		Type:    string(events.EventTypeDeploymentCompleted),
		Data: &events.DeploymentEventData{
			DeploymentId: int(deployment.ID),
			OperationId:  to.Ptr(operation.ID.String()),
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
	return deployment.Create(*azureDeployment)
}

func backOffRetryHandler(n uint, err error, config *retry.Config) time.Duration {
	fmt.Println("Deployment failed with: " + err.Error())
	if retriable, ok := err.(*RetriableError); ok {
		fmt.Printf("Retry after %v\n", retriable.RetryAfter)
		return retriable.RetryAfter
	}
	return retry.BackOffDelay(n, err, config)
}

//region factory

func NewStartDeploymentOperation(appConfig *config.AppConfig) Executor {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
	dryRunOperation := &startDeployment{
		db: db,
	}
	return dryRunOperation
}

//endregion factory
