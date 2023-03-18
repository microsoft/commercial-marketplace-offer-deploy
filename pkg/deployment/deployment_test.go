package deployment

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armpolicy"
	"github.com/joho/godotenv"
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

	setupResourceGroup()
	deployPolicyDefinition()
	deployPolicy()

	exitVal := m.Run()
	log.Println("Cleaning up resources after the tests here")
	
	os.Exit(exitVal)
}

func setupResourceGroup() {
	resp, err := createResourceGroup(subscriptionId, resourceGroupName, location)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("%s was created", *resp.Name)
	}
}

func deployPolicyDefinition() {
	log.Printf("Inside deployPolicyDefinition()")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armpolicy.NewDefinitionsClient(subscriptionId, cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.CreateOrUpdate(ctx,
		"ResourceNaming",
		armpolicy.Definition{
			Properties: &armpolicy.DefinitionProperties{
				Description: to.Ptr("Force resource names to begin with given 'prefix' and/or end with given 'suffix'"),
				DisplayName: to.Ptr("Enforce resource naming convention"),
				Metadata: map[string]interface{}{
					"category": "Naming",
				},
				Mode: to.Ptr("All"),
				// Parameters: map[string]*armpolicy.ParameterDefinitionsValue{
				// 	"prefix": {
				// 		Type: to.Ptr(armpolicy.ParameterTypeString),
				// 		Metadata: &armpolicy.ParameterDefinitionsValueMetadata{
				// 			Description: to.Ptr("Resource name prefix"),
				// 			DisplayName: to.Ptr("Prefix"),
				// 		},
				// 	},
				// 	"suffix": {
				// 		Type: to.Ptr(armpolicy.ParameterTypeString),
				// 		Metadata: &armpolicy.ParameterDefinitionsValueMetadata{
				// 			Description: to.Ptr("Resource name suffix"),
				// 			DisplayName: to.Ptr("Suffix"),
				// 		},
				// 	},
				// },
				PolicyRule: map[string]interface{}{
					"if": map[string]interface{}{
						"not": map[string]interface{}{
							"field": "name",
							"like":  "a*b",
							//"like":  "[concat(parameters('prefix'), '*', parameters('suffix'))]",
						},
					},
					"then": map[string]interface{}{
						"effect": "deny",
					},
				},
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

func deployPolicy() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armpolicy.NewAssignmentsClient(subscriptionId, cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	scope := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionId, resourceGroupName)
	log.Printf("scope is %s", scope)
	
	policyDefinitionId := fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Authorization/policyDefinitions/ResourceNaming", subscriptionId)
	log.Printf("policyDefinitionId is %s", policyDefinitionId)

	_, err = client.Create(ctx,
		scope,
		"ResourceName",
		armpolicy.Assignment{
			Properties: &armpolicy.AssignmentProperties{
				Description: to.Ptr("Enforce resource naming conventions"),
				DisplayName: to.Ptr("Enforce Resource Names"),
				Scope: &scope,
				Metadata: map[string]interface{}{
					"assignedBy": "John Doe",
				},
				NonComplianceMessages: []*armpolicy.NonComplianceMessage{
					{
						Message: to.Ptr("A resource name was non-complaint.  It must be in the format 'a*b'."),
					}},
				PolicyDefinitionID: to.Ptr(policyDefinitionId),
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
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
		subscriptionId: subscriptionId,
		location: location,
		resourceGroupName: resourceGroupName,
		deploymentName: "modmdeploy",
		template: template,
		params: parameters,
	}
	if got, err := Create(deployment); err != nil {
		t.Errorf("TestCreateSuccess() returned with an error %s", err)
	} else {
		gotdata, _ := json.Marshal(got)
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
		subscriptionId: subscriptionId,
		location: location,
		resourceGroupName: resourceGroupName,
		deploymentName: "modmdeploy",
		template: template,
		params: parameters,
	}
	if got := DryRun(deployment); got == nil {
		t.Errorf("TestDryRunPolicySuccess() return json with a lenth of 0")
	} else {
		gotdata, _ := json.Marshal(got)
		jsonData := string(gotdata)
		log.Printf("TestDryRunPolicySuccess result - %s", jsonData)
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
		subscriptionId: subscriptionId,
		location: location,
		resourceGroupName: resourceGroupName,
		deploymentName: "modmdeploy",
		template: template,
		params: parameters,
	}
	if got := DryRun(deployment); got == nil {
		t.Errorf("TestDryRunPolicyFailure() return json with a lenth of 0")
	} else {
		gotdata, _ := json.Marshal(got)
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
		subscriptionId: subscriptionId,
		location: location,
		resourceGroupName: resourceGroupName,
		deploymentName: "modmdeploy",
		template: template,
		params: parameters,
	}
	if got := DryRun(deployment); got == nil {
		t.Errorf("DryRun() return json with a lenth of 0")
	} else {
		gotdata, _ := json.Marshal(got)
		jsonData := string(gotdata)
		log.Printf("TestDryRunQuotsFailure result - %s", jsonData)
	}
}