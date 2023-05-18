package deployment_test

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/stretchr/testify/assert"
)

func TestStartDeployment(t *testing.T) {
	fullPath := filepath.Join("../../test/testdata/nameviolation/nestedfailure", "mainTemplate.json")
	template, err := utils.ReadJson(fullPath)
	assert.NoError(t, err)
	assert.NotNil(t, template)

	resources := deployment.FindResourcesByType(template, "Microsoft.Resources/deployments")
	assert.Greater(t, len(resources), 0)
}

func TestGetDeploymentParams(t *testing.T) {
	templatePath := filepath.Join("../../test/testdata/missingparam", "mainTemplate.json")
	template, err := utils.ReadJson(templatePath)
	assert.NoError(t, err)
	assert.NotNil(t, template)

	paramsPath := filepath.Join("../../test/testdata/missingparam", "parameters.json")
	params, err := utils.ReadJson(paramsPath)
	assert.NoError(t, err)
	assert.NotNil(t, params)

	azureDeployment := &deployment.AzureDeployment{
		SubscriptionId: uuid.NewString(),
		Location:       "eastus",
		ResourceGroupName: "TestResourceGroup",
		DeploymentName: "TestDeployment",
		DeploymentType: deployment.AzureResourceManager,
		Template:       template,
		Params:         params,
	}

	result := azureDeployment.GetParams()
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	assert.NotNil(t, result["aapName"])
	assert.NotNil(t, result["testName"])
}