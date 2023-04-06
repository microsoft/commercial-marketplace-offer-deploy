package operations

import (
	"context"
	"log"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"gorm.io/gorm"
)

type startDeploymentOperation struct {
	db *gorm.DB
}

func (p *startDeploymentOperation) Invoke(operation *data.InvokedOperation) error {
	tx := p.db.Begin()

	deployment := &data.Deployment{}
	tx.First(&deployment, operation.DeploymentId)
	deployment.Status = events.DeploymentStartingEventType.String()
	tx.Save(deployment)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()

	azureDeployment := p.mapAzureDeployment(deployment, operation)

	go func() {
		_, err := p.deploy(context.TODO(), azureDeployment)
		if err != nil {
			log.Println("Error calling deployment.Create: ", err)
		}
	}()

	deployment.Status = events.DeploymentStartedEventType.String()
	tx.Save(&deployment)

	tx.Commit()

	return nil
}

func (p *startDeploymentOperation) mapAzureDeployment(d *data.Deployment, io *data.InvokedOperation) *deployment.AzureDeployment {
	return &deployment.AzureDeployment{
		SubscriptionId:    d.SubscriptionId,
		ResourceGroupName: d.ResourceGroup,
		DeploymentName:    d.GetAzureDeploymentName(),
		Template:          d.Template,
		Params:            io.Parameters,
	}
}

func (p *startDeploymentOperation) deploy(ctx context.Context, azureDeployment *deployment.AzureDeployment) (*deployment.AzureDeploymentResult, error) {
	return deployment.Create(*azureDeployment)
}

//region factory

func NewStartDeploymentOperation(appConfig *config.AppConfig) DeploymentOperation {
	db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
	dryRunOperation := &startDeploymentOperation{
		db: db,
	}
	return dryRunOperation
}

//endregion factory
