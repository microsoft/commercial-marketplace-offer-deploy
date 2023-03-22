package handlers

import (
	"errors"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
)

func CreateDryRun(deploymentId int, operation internal.InvokeDeploymentOperation, d data.Database) (interface{}, error) {
	// call database to get the local template
	if d == nil {
		return nil, errors.New("database is nil")
	}

	d.Instance().AutoMigrate(&data.Deployment{})
	retrieved := &data.Deployment{}
	d.Instance().First(retrieved, "id = ?", deploymentId)
	
	templateParams := operation.Parameters
	if templateParams == nil {
		return nil, errors.New("templateParams were not provided")
	}

	azureDeployment := deployment.AzureDeployment {
		SubscriptionId: retrieved.SubscriptionId,
		Location: retrieved.Location,
		ResourceGroupName: retrieved.ResourceGroup,
		DeploymentName: retrieved.Name,
		Template: retrieved.Template,
		Params: templateParams.(map[string]interface{}),
	}
	res := deployment.DryRun(&azureDeployment)
	return res, nil
}