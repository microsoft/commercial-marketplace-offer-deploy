package handlers

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"gorm.io/gorm"
)

func CreateDryRun(deploymentId int, operation api.InvokeDeploymentOperation, db *gorm.DB) (interface{}, error) {
	log.Printf("Inisde CreateDryRun deploymentId: %d", deploymentId)

	retrieved := &data.Deployment{}
	db.First(&retrieved, deploymentId)

	log.Printf("Inisde CreateDryRun deploymentId: %v", retrieved)

	templateParams := operation.Parameters
	if templateParams == nil {
		return nil, errors.New("templateParams were not provided")
	}

	azureDeployment := deployment.AzureDeployment{
		SubscriptionId:    retrieved.SubscriptionId,
		Location:          retrieved.Location,
		ResourceGroupName: retrieved.ResourceGroup,
		DeploymentName:    retrieved.Name,
		Template:          retrieved.Template,
		Params:            templateParams.(map[string]interface{}),
	}

	res := deployment.DryRun(&azureDeployment)
	uuid := uuid.New().String()
	timestamp := time.Now().UTC()
	status := "OK"
	returnedResult := api.InvokedOperation{
		ID:        &uuid,
		InvokedOn: &timestamp,
		Result:    *res,
		Status:    &status,
	}
	return returnedResult, nil
}
