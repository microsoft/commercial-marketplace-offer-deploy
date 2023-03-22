package handlers

import (
	"errors"
	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/models"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/persistence"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/persistence/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
)

func CreateDryRun(operation models.InvokeDeploymentOperation, d persistence.Database) (interface{}, error) {
	// call database to get the local template
	if d == nil {
		return nil, errors.New("Database is nil")
	}
	d.Instance().AutoMigrate(&model.Deployment{})
	retrieved := &model.Deployment{}
	d.Instance().First(retrieved, "name = ?", )
	template, err := utils.ReadJson(retrieved.Template.FilePath)
	if err != nil {
		return nil, err
	} 
	//var subscriptionId string
	//var resourceGroupName string
	//var location string
	azureDeployment := deployment.AzureDeployment {
		Template: template,
		Params: operation.Parameters,
	}
	res := deployment.DryRun(&azureDeployment)
	return res, nil
}