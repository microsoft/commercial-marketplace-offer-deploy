package deployment

import (
	"context"
	"encoding/json"
	"io/ioutil"
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
	resp, err := CreateResourceGroup(subscriptionId, resourceGroupName, location)
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

	scope := "/subscriptions/" + subscriptionId + "/resourceGroups/" + resourceGroupName
	log.Printf("scope is %s", scope)
	policyDefinitionId := "/subscriptions/" + subscriptionId + "/providers/Microsoft.Authorization/policyDefinitions/ResourceNaming"
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

func readJson(path string) (map[string]interface{}, error) {
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

func TestDryRunSuccess(t *testing.T) {
	log.Printf("Inside TestDryRunSuccess")

	templatePath := "../../test/deployment/cogsvcspub/success/mainTemplate.json" 
	template, err := readJson(templatePath)
	if err != nil {
		t.Errorf("DryRun() could not read templateFile")
	} 

	parametersPath :=  "../../test/deployment/cogsvcspub/success/parameters.json"
	parameters, err := readJson(parametersPath)
	if err != nil {
		t.Errorf("DryRun() could not read parameters")
	} 

	deployment := &AzureDeployment{
		subscriptionId: subscriptionId,
		location: location,
		resourceGroupName: resourceGroupName,
		deploymentName: "modmdeploy",
		template: template,
		params: parameters,
	}
	if got := DryRun(deployment); len(got) == 0{
		t.Errorf("DryRun() return json with a lenth of 0")
	} else {
		log.Printf("TestDryRunResult - %s", got)
	}
}

func TestDryRunFailure(t *testing.T) {
	log.Printf("Inside TestDryRunFailure")

	templatePath := "../../test/deployment/cogsvcspub/failure/mainTemplate.json" 
	template, err := readJson(templatePath)
	if err != nil {
		t.Errorf("DryRun() could not read templateFile")
	} 

	parametersPath :=  "../../test/deployment/cogsvcspub/failure/parameters.json"
	parameters, err := readJson(parametersPath)
	if err != nil {
		t.Errorf("DryRun() could not read parameters")
	} 

	deployment := &AzureDeployment{
		subscriptionId: subscriptionId,
		location: location,
		resourceGroupName: resourceGroupName,
		deploymentName: "modmdeploy",
		template: template,
		params: parameters,
	}
	if got := DryRun(deployment); len(got) == 0{
		t.Errorf("DryRun() return json with a lenth of 0")
	} else {
		log.Printf("TestDryRunResult - %s", got)
	}
}