package deployment

import (
	"context"
	"encoding/json"
	"io/ioutil"

	// "encoding/json"
	// "io/ioutil"
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
	deployPolicy()

	exitVal := m.Run()
	log.Println("Cleaning up resources after the tests here")
	
	// err = DeleteResourceGroup(subscriptionId, resourceGroupName)
	// if err != nil {
	// 	log.Println("Error deleting resource group ")
	// }

	os.Exit(exitVal)
}

func setupResourceGroup() {
	// err := DeleteResourceGroup(subscriptionId, resourceGroupName)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	resp, err := CreateResourceGroup(subscriptionId, resourceGroupName, location)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("%s was created", *resp.Name)
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

	scope := "subscriptions/" + subscriptionId + "/resourceGroups/" + resourceGroupName
	log.Printf("scope is %s", scope)
	policyDefinitionId := "/providers/Microsoft.Authorization/policyDefinitions/037eea7a-bd0a-46c5-9a66-03aea78705d3"
	log.Printf("policyDefinitionId is %s", policyDefinitionId)

	_, err = client.Create(ctx,
		scope,
		"CogSvcs",
		armpolicy.Assignment{
			Properties: &armpolicy.AssignmentProperties{
				Description: to.Ptr("Forces cog service deployments to be locked down"),
				DisplayName: to.Ptr("Protect Cog Svcs"),
				Scope: &scope,
				Metadata: map[string]interface{}{
					"assignedBy": "John Doe",
				},
				NonComplianceMessages: []*armpolicy.NonComplianceMessage{
					{
						Message: to.Ptr("Cognitive Services must be locked down."),
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

func TestDryRun(t *testing.T) {
	log.Printf("Inside TestDryRun")

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