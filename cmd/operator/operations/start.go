package operations

import (
	"context"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"gorm.io/gorm"
)

type startDeployment struct {
	db *gorm.DB
}

func (exe *startDeployment) Execute(ctx context.Context, operation *data.InvokedOperation) error {
	deployment, err := exe.updateToRunning(ctx, operation)
	if err != nil {
		return err
	}

	azureDeployment := exe.mapAzureDeployment(deployment, operation)

	go func() {
		_, err := exe.deploy(ctx, azureDeployment)
		if err != nil {
			log.Error("Error calling deployment.Create: ", err)
			err = exe.updateToFailed(ctx, operation, err)

			if err != nil {
				log.Error("Error updating Deployment to failed: ", err)
			}
		}
	}()
	return nil
}

func (exe *startDeployment) updateToRunning(ctx context.Context, operation *data.InvokedOperation) (*data.Deployment, error) {
	db := exe.db

	deployment := &data.Deployment{}
	db.First(&deployment, operation.DeploymentId)
	deployment.Status = string(events.EventTypeRunning)
	db.Save(deployment)

	err := hook.Add(&events.EventHookMessage{
		Subject:   "/deployments/" + strconv.Itoa(int(deployment.ID)),
		EventType: events.EventTypeRunning.String(),
		Body: &events.EventHookDeploymentEventMessageBody{
			DeploymentId: int(deployment.ID),
			Status:       "Success",
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
	deployment.Status = string(events.EventTypeFailed)
	db.Save(deployment)

	err = hook.Add(&events.EventHookMessage{
		Subject:   "/deployments/" + strconv.Itoa(int(deployment.ID)),
		EventType: events.EventTypeFailed.String(),
		Body: &events.EventHookDeploymentEventMessageBody{
			DeploymentId: int(deployment.ID),
			Status:       "Failed",
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

//region factory

func NewStartDeploymentOperation(appConfig *config.AppConfig) Executor {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
	dryRunOperation := &startDeployment{
		db: db,
	}
	return dryRunOperation
}

//endregion factory
