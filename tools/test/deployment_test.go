package test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	//	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"

	//"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2021-04-01/resources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	//	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/stretchr/testify/assert"
	//	"github.com/stretchr/testify/require"
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



	// convertedTemplate, err := convertForRedeploy(template.Template)
	// if err != nil {
	// 	assert.Fail(s.T(), "Failed to convert template")
	// }

	// convertedParams, err := convertForRedeploy(params)
	// if err != nil {
	// 	assert.Fail(s.T(), "Failed to convert params")
	// }

	paramsFromFile := getParmsAsMap(s, "/Users/bobjacobs/work/src/github.com/microsoft/commercial-marketplace-offer-deploy/test/testdata/nameviolation/nestedfailure/parameters.json")
	if paramsFromFile == nil {
		assert.Fail(s.T(), "Failed to get params from file")
	}

	var handCraftedParams map[string]interface{}
	handCraftedParams = make(map[string]interface{})

	kindMap := make(map[string]interface{})
	kindMap["value"] = "StorageV2"

	handCraftedParams["kind"] = kindMap

	locationMap := make(map[string]interface{})
	locationMap["value"] = "eastus"

	handCraftedParams["location"] = locationMap

	nameMap := make(map[string]interface{})
	nameMap["value"] = "bobjacbicepsa"

	handCraftedParams["name"] = nameMap

	skuNameMap := make(map[string]interface{})
	skuNameMap["value"] = "Standard_LRS"

	handCraftedParams["sku_name"] = skuNameMap

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

func getParmsAsMap(s *deploymentSuite, paramsFile string) map[string]interface{} {
	params, err := ReadJson(paramsFile)
	if err != nil {
		assert.Fail(s.T(), "Failed to read params file")
	}
	return params
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
		//keys = append(keys, k)
		valueMap := make(map[string]interface{})
		templateValueMap := params[k].(map[string]interface{})

		valueMap["value"] = templateValueMap["value"]
		paramValues[k] = valueMap
	}

	// for _, k := range keys {
	// 	valueMap := make(map[string]interface{})
	// 	valueMap["value"] = template[k]["value"]
	// 	paramValues[k] = valueMap
	// }

	return paramValues
}



func convertForRedeploy(in interface{}) (map[string]interface{}, error) {
	m, ok := in.(map[string]interface{})
    if !ok {
        return nil, errors.New("input is not a map")
    }
    return m, nil
}
