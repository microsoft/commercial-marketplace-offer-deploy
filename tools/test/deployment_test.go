package test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type deploymentSuite struct {
	suite.Suite
	subscriptionId    		string
	resourceGroupName 		string
	parentDeploymentName 	string
	childDeploymentName 	string
	location          		string
	endpoint          		string

}

func TestDryRunSuite(t *testing.T) {
	suite.Run(t, &deploymentSuite{})
}

func (s *deploymentSuite) SetupSuite() {
	s.subscriptionId = "31e9f9a0-9fd2-4294-a0a3-0101246d9700"
	s.resourceGroupName = "demo2"
	s.parentDeploymentName = "modm.1-TaggedDeployment"
	s.childDeploymentName = "storageAccounts"
	s.location = "eastus"
	s.endpoint = "http://localhost:8080"
}

func (s *deploymentSuite) TestExportNestedDeployment() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Print(err)
	}
	
	deploymentsClient, err := armresources.NewDeploymentsClient(s.subscriptionId, cred, nil)
	if err != nil {
		assert.Fail(s.T(), "Failed to get deployments client")
	}

	ctx := context.Background()

	deployment, err := deploymentsClient.Get(ctx, s.resourceGroupName, s.childDeploymentName, nil)
	if err != nil {
		assert.Fail(s.T(), "Failed to get deployment")
	}

	if deployment.Properties == nil || deployment.Properties.TemplateLink == nil {
		fmt.Printf("Deployment template not found: %v\n", err)
	}

	params := (*deployment.DeploymentExtended.Properties).Parameters
	if params == nil {
		fmt.Printf("Deployment parameters not found: %v\n", err)
	}

	castParams := params.(map[string]interface{})
	if castParams == nil {
		fmt.Printf("Deployment parameters not found: %v\n", err)
	}

	template, err := deploymentsClient.ExportTemplate(ctx, s.resourceGroupName, s.childDeploymentName, nil)
	if err != nil {
		assert.Fail(s.T(), "Failed to export template")
	}

	if template.Template == nil {
		fmt.Printf("Deployment template not found: %v\n", err)
	}

	t := template.Template
	if t == nil {
		assert.Fail(s.T(), "Failed to get template")
	}

	paramValuesMap := getParamsMapFromTemplate(template.Template.(map[string]interface{}), castParams)
	
	deploymentPollerResp, err := deploymentsClient.BeginCreateOrUpdate(
		ctx,
		s.resourceGroupName,
		s.childDeploymentName,
		armresources.Deployment{
			Properties: &armresources.DeploymentProperties{
				Template:   template.Template.(map[string]interface{}),
				Parameters: paramValuesMap,
				Mode:       to.Ptr(armresources.DeploymentModeIncremental),
			},
		},
		nil)
	
	if err != nil {
		assert.Fail(s.T(), "Failed to create deployment")
	}

	resp, err := deploymentPollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		assert.Fail(s.T(), "Failed to poll until done")
	}
	if resp.ID == nil {
		assert.Fail(s.T(), "Failed to get response")
	}
}

func ReadJson(path string) (map[string]interface{}, error) {
	templateFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	template := make(map[string]interface{})
	if err := json.Unmarshal(templateFile, &template); err != nil {
		return nil, err
	}
	return template, nil
}

func getParamsMapFromTemplate(template map[string]interface{}, params map[string]interface{}) map[string]interface{} { 
	paramValues := make(map[string]interface{})
	
	templateParams := template["parameters"].(map[string]interface{})
	for k := range templateParams {
		valueMap := make(map[string]interface{})
		templateValueMap := params[k].(map[string]interface{})

		valueMap["value"] = templateValueMap["value"]
		paramValues[k] = valueMap
	}

	return paramValues
}


