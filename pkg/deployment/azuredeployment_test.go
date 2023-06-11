package deployment

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestStartDeployment(t *testing.T) {
	fullPath := filepath.Join("../../test/testdata/nameviolation/nestedfailure", "mainTemplate.json")
	template, err := utils.ReadJson(fullPath)
	assert.NoError(t, err)
	assert.NotNil(t, template)

	resources := findResourcesByType(template, "Microsoft.Resources/deployments")
	assert.Greater(t, len(resources), 0)
}

func TestGetDeploymentParamsNested(t *testing.T) {
	templatePath := filepath.Join("../../test/testdata/missingparam", "mainTemplate.json")
	template, err := utils.ReadJson(templatePath)
	assert.NoError(t, err)
	assert.NotNil(t, template)

	paramsPath := filepath.Join("../../test/testdata/missingparam", "parameters.json")
	params, err := utils.ReadJson(paramsPath)
	assert.NoError(t, err)
	assert.NotNil(t, params)

	azureDeployment := &AzureDeployment{
		SubscriptionId:    uuid.NewString(),
		Location:          "eastus",
		ResourceGroupName: "TestResourceGroup",
		DeploymentName:    "TestDeployment",
		DeploymentType:    DeploymentTypeARM,
		Template:          template,
		Params:            params,
	}

	result := azureDeployment.GetParameters()
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	assert.NotNil(t, result["aapName"])
	assert.NotNil(t, result["testName"])
}

func TestGetDeploymentParamsUnNested(t *testing.T) {
	templatePath := filepath.Join("../../test/testdata/missingparam", "mainTemplate.json")
	template, err := utils.ReadJson(templatePath)
	assert.NoError(t, err)
	assert.NotNil(t, template)

	paramsMap := make(map[string]interface{})

	aapNameMap := make(map[string]interface{})
	aapNameMap["value"] = "test"

	testNameMap := make(map[string]interface{})
	testNameMap["value"] = "test2"

	paramsMap["aapName"] = aapNameMap
	paramsMap["testName"] = testNameMap

	azureDeployment := &AzureDeployment{
		SubscriptionId:    uuid.NewString(),
		Location:          "eastus",
		ResourceGroupName: "TestResourceGroup",
		DeploymentName:    "TestDeployment",
		DeploymentType:    DeploymentTypeARM,
		Template:          template,
		Params:            paramsMap,
	}

	result := azureDeployment.GetParameters()
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	assert.NotNil(t, result["aapName"])
	assert.NotNil(t, result["testName"])

	aapNameValue := result["aapName"].(map[string]interface{})["value"]
	assert.Equal(t, "test", aapNameValue)

	testNameValue := result["testName"].(map[string]interface{})["value"]
	assert.Equal(t, "test2", testNameValue)
}

func TestGetDeploymentTemplateParamsNested(t *testing.T) {
	templatePath := filepath.Join("../../test/testdata/missingparam", "mainTemplate.json")
	template, err := utils.ReadJson(templatePath)
	assert.NoError(t, err)
	assert.NotNil(t, template)

	paramsPath := filepath.Join("../../test/testdata/missingparam", "parameters.json")
	params, err := utils.ReadJson(paramsPath)
	assert.NoError(t, err)
	assert.NotNil(t, params)

	azureDeployment := &AzureDeployment{
		SubscriptionId:    uuid.NewString(),
		Location:          "eastus",
		ResourceGroupName: "TestResourceGroup",
		DeploymentName:    "TestDeployment",
		DeploymentType:    DeploymentTypeARM,
		Template:          template,
		Params:            params,
	}

	result := azureDeployment.GetParametersFromTemplate()
	assert.NotNil(t, result)
	assert.Equal(t, 3, len(result))
	assert.NotNil(t, result["aapName"])
	assert.NotNil(t, result["testName"])
	assert.NotNil(t, result["testName2"])
}

//region helpers

func findResourcesByType(template AzureTemplate, resourceType string) []string {
	deploymentResources := []string{}
	if template != nil && template["resources"] != nil {
		resources := template["resources"].([]interface{})
		for _, resource := range resources {
			resourceMap := resource.(map[string]interface{})
			if resourceMap["type"] != nil && resourceMap["type"].(string) == resourceType {
				deploymentResources = append(deploymentResources, resourceMap["name"].(string))
			}
		}
	}
	return deploymentResources
}

//endregion helpers
