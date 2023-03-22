package handlers

import (
	"errors"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
)

func CreateDryRun(operation internal.InvokeDeploymentOperation, d data.Database) (interface{}, error) {
	// call database to get the local template
	if d == nil {
		return nil, errors.New("database is nil")
	}
	d.Instance().AutoMigrate(&model.Deployment{})
	retrieved := &model.Deployment{}
	d.Instance().First(retrieved, "name = ?", )
	template, err := utils.ReadJson(retrieved.Template.FilePath)
	if err != nil {
		return nil, err
	} 
	paramsMap := operation.Parameters.(map[string]interface{})
	deploymentParams := paramsMap["deploymentParams"]
	if deploymentParams == nil {
		return nil, errors.New("deploymentParams were not provided")
	}
	deploymentParamsMap := deploymentParams.(map[string]interface{})

	templateParams := paramsMap["templateParams"]
	if templateParams == nil {
		return nil, errors.New("templateParams were not provided")
	}

	subscriptionId := deploymentParamsMap["subscriptionId"]
	location := deploymentParamsMap["location"]
	resourceGroupName := deploymentParamsMap["resourceGroupName"]
	deploymentName := deploymentParamsMap["deploymentName"]
	deploymentType := deploymentParamsMap["deploymentType"]
	resumeToken :=  deploymentParamsMap["resumeToken"]

	azureDeployment := deployment.AzureDeployment {
		SubscriptionId: subscriptionId.(string),
		Location: location.(string),
		ResourceGroupName: resourceGroupName.(string),
		DeploymentName: deploymentName.(string),
		DeploymentType: deploymentType.(deployment.DeploymentType),
		ResumeToken: resumeToken.(string),
		Template: template,
		Params: templateParams.(map[string]interface{}),
	}
	res := deployment.DryRun(&azureDeployment)
	return res, nil
}