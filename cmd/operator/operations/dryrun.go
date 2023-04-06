package operations

import (
	"time"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"gorm.io/gorm"
)

type dryRunOperation struct {
	db      *gorm.DB
	process DryRunProcessorFunc
}

func (h *dryRunOperation) Handle(operation *data.InvokedOperation) error {
	azureDeployment := h.getAzureDeployment(operation)
	response := deployment.DryRun(azureDeployment)

	operation.Status = *response.Status
	operation.Result = response.DryRunResult
	operation.UpdatedAt = time.Now().UTC()

	err := h.save(operation)
	if err != nil {
		return err
	}
	return nil
}

func (h *dryRunOperation) getAzureDeployment(operation *data.InvokedOperation) *deployment.AzureDeployment {
	retrieved := &data.Deployment{}
	h.db.First(&retrieved, operation.DeploymentId)

	return &deployment.AzureDeployment{
		SubscriptionId:    retrieved.SubscriptionId,
		Location:          retrieved.Location,
		ResourceGroupName: retrieved.ResourceGroup,
		DeploymentName:    retrieved.Name,
		Template:          retrieved.Template,
		Params:            operation.Parameters,
	}
}

func (h *dryRunOperation) save(operation *data.InvokedOperation) error {
	tx := h.db.Begin()
	tx.Save(&operation)

	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()

	return nil
}

//region factory

func NewDryRunHandler(appConfig *config.AppConfig) DeploymentOperationHandlerFunc {
	return func(operation *data.InvokedOperation) error {
		db := data.NewDatabase(appConfig.GetDatabaseOptions()).Instance()
		dryRunOperation := &dryRunOperation{
			db:      db,
			process: deployment.DryRun,
		}
		return dryRunOperation.Handle(operation)
	}
}

//endregion factory
