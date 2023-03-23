package deployment

import (
//	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	// "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	// "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	// "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armpolicy"
	"github.com/joho/godotenv"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
)

var subscriptionId string
var resourceGroupName string
var location string

func TestMain(m *testing.M) {
	log.Println("Test setup beginning")
	err := godotenv.Load(".env") 
	if err != nil {
		log.Println("Cannot load environment variables from .env")
	} 
	subscriptionId = os.Getenv("AZURE_SUBSCRIPTION_ID")
	resourceGroupName = "MODMTest"
	location = "eastus"

	utils.SetupResourceGroup(subscriptionId, resourceGroupName, location)
	utils.DeployPolicyDefinition(subscriptionId)
	utils.DeployPolicy(subscriptionId, resourceGroupName)

	exitVal := m.Run()
	log.Println("Cleaning up resources after the tests here")
	
	os.Exit(exitVal)
}


func TestCreateSuccess(t *testing.T) {
	log.Printf("Inside TestCreateSuccess")

	templatePath := "../../test/deployment/resourcename/success/mainTemplate.json" 
	template, err := readJson(templatePath)
	if err != nil {
		t.Errorf("TestDryRunPolicySuccess() could not read templateFile")
	} 

	parametersPath :=  "../../test/deployment/resourcename/success/parameters.json"
	parameters, err := readJson(parametersPath)
	if err != nil {
		t.Errorf("TestDryRunPolicySuccess() could not read parameters")
	} 

	deployment := AzureDeployment{
		SubscriptionId: subscriptionId,
		Location: location,
		ResourceGroupName: resourceGroupName,
		DeploymentName: "modmdeploy",
		Template: template,
		Params: parameters,
	}
	if got, err := Create(deployment); err != nil {
		t.Errorf("TestCreateSuccess() returned with an error %s", err)
	} else {
		gotdata, _ := json.MarshalIndent(got, "", "  ")
		jsonData := string(gotdata)
		log.Printf("TestDryRunPolicySuccess result - %s", jsonData)
	}
}

func TestDryRunPolicySuccess(t *testing.T) {
	log.Printf("Inside TestDryRunPolicySuccess")

	templatePath := "../../test/deployment/resourcename/success/mainTemplate.json" 
	template, err := readJson(templatePath)
	if err != nil {
		t.Errorf("TestDryRunPolicySuccess() could not read templateFile")
	} 

	parametersPath :=  "../../test/deployment/resourcename/success/parameters.json"
	parameters, err := readJson(parametersPath)
	if err != nil {
		t.Errorf("TestDryRunPolicySuccess() could not read parameters")
	} 

	deployment := &AzureDeployment{
		SubscriptionId: subscriptionId,
		Location: location,
		ResourceGroupName: resourceGroupName,
		DeploymentName: "modmdeploy",
		Template: template,
		Params: parameters,
	}
	if got := DryRun(deployment); got == nil {
		t.Errorf("TestDryRunPolicySuccess() return json with a lenth of 0")
	} else {
		gotdata, _ := json.MarshalIndent(got, "", "  ")
		jsonData := string(gotdata)
		log.Printf("TestDryRunPolicySuccess result - %s", jsonData)
	}
}

func TestDryRunNestedPolicyFailure(t *testing.T) {
	log.Printf("Inside TestDryRunNestedPolicyFailure")

	templatePath := "../../test/deployment/resourcename/nested/mainTemplate.json" 
	template, err := readJson(templatePath)
	if err != nil {
		t.Errorf("TestDryRunNestedPolicyFailure() could not read templateFile")
	} 

	parametersPath :=  "../../test/deployment/resourcename/nested/parameters.json"
	parameters, err := readJson(parametersPath)
	if err != nil {
		t.Errorf("TestDryRunNestedPolicyFailure() could not read parameters")
	} 

	deployment := &AzureDeployment{
		SubscriptionId: subscriptionId,
		Location: location,
		ResourceGroupName: resourceGroupName,
		DeploymentName: "modmdeploy",
		Template: template,
		Params: parameters,
	}
	if got := DryRun(deployment); got == nil {
		t.Errorf("TestDryRunNestedPolicyFailure() return json with a lenth of 0")
	} else {
		gotdata, _ := json.MarshalIndent(got, "", "  ")
		jsonData := string(gotdata)
		log.Printf("TestDryRunNestedPolicyFailure result - %s", jsonData)
	}
}

func TestDryRunPolicyFailure(t *testing.T) {
	log.Printf("Inside TestDryRunPolicyFailure")

	templatePath := "../../test/deployment/resourcename/failure/mainTemplate.json" 
	template, err := readJson(templatePath)
	if err != nil {
		t.Errorf("TestDryRunPolicyFailure() could not read templateFile")
	} 

	parametersPath :=  "../../test/deployment/resourcename/failure/parameters.json"
	parameters, err := readJson(parametersPath)
	if err != nil {
		t.Errorf("TestDryRunPolicyFailure() could not read parameters")
	} 

	deployment := &AzureDeployment{
		SubscriptionId: subscriptionId,
		Location: location,
		ResourceGroupName: resourceGroupName,
		DeploymentName: "modmdeploy",
		Template: template,
		Params: parameters,
	}
	if got := DryRun(deployment); got == nil {
		t.Errorf("TestDryRunPolicyFailure() return json with a lenth of 0")
	} else {
		gotdata, _ := json.MarshalIndent(got, "", "  ")
		jsonData := string(gotdata)
		log.Printf("TestDryRunPolicyFailure result - %s", jsonData)
	}
}

func TestDryRunQuotsFailure(t *testing.T) {
	log.Printf("Inside TestDryRunQuotaFailure")

	templatePath := "../../test/deployment/quota/failure/mainTemplate.json" 
	template, err := readJson(templatePath)
	if err != nil {
		t.Errorf("TestDryRunQuotsFailure() could not read templateFile")
	} 

	parametersPath :=  "../../test/deployment/quota/failure/parameters.json"
	parameters, err := readJson(parametersPath)
	if err != nil {
		t.Errorf("TestDryRunQuotsFailure() could not read parameters")
	} 

	deployment := &AzureDeployment{
		SubscriptionId: subscriptionId,
		Location: location,
		ResourceGroupName: resourceGroupName,
		DeploymentName: "modmdeploy",
		Template: template,
		Params: parameters,
	}
	if got := DryRun(deployment); got == nil {
		t.Errorf("DryRun() return json with a lenth of 0")
	} else {
		gotdata, _ := json.MarshalIndent(got, "", "  ")
		jsonData := string(gotdata)
		log.Printf("TestDryRunQuotsFailure result - %s", jsonData)
	}
}

func TestConvertMapParams(t *testing.T) {
	// var deploymentParamsObject interface{}
	// deploymentParams := map[string]interface{} {
	// 	"subscriptionId": "sub1",
	// 	"location": "westus",
	// }
	// deploymentParamsObject = deploymentParams

	var templateParamsObject interface{}
	templateParams := map[string]interface{}{
		"param1": "val1",
		"param2": 5,
	}
	templateParamsObject = templateParams

	var templateParamsConvert map[string]interface{}
	templateParamsConvert = templateParamsObject.(map[string]interface{})
	param1Interface := templateParamsConvert["param1"]
	param1Str := fmt.Sprintf("%v", param1Interface)
	if param1Str != "val1" {
		t.Errorf("param1Str does not equal 'val1")
	}
}