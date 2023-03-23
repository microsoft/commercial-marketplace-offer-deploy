package sdk

import (
	"context"
	"encoding/json"
	"log"
	"path/filepath"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/generated"
)

var subscriptionId string
var resourceGroupName string
var location string
var endpoint string = "http://localhost:8080"

func setupForTestDryRun() {
	subscriptionId = "31e9f9a0-9fd2-4294-a0a3-0101246d9700"
	resourceGroupName = "aMODMTestb"
	location = "eastus"

	if exists, _ := utils.DoesResourceGroupExist(subscriptionId, resourceGroupName, location); exists {
		return
	}

	utils.SetupResourceGroup(subscriptionId, resourceGroupName, location)
	utils.DeployPolicyDefinition(subscriptionId)
	utils.DeployPolicy(subscriptionId, resourceGroupName)
}

func TestDryRun(t *testing.T) {
	setupForTestDryRun()

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("Authentication failure: %+v", err)
	}

	client, err := NewClient(endpoint, cred, nil)

	if err != nil {
		log.Print("Client construction failed.")
	}

	testNamePolicyFailure(client)
	testQuotaViolation(client)
}

func testQuotaViolation(client *Client) {
	ctx := context.TODO()

	deployment := createDeployment(ctx, client, quotaViolationPath)
	result, err := client.DryRunDeployment(ctx, *deployment.ID, getParameters(quotaViolationPath))

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Validation Results:\n %s" + *prettify(result.Results))
}

func testNamePolicyFailure(client *Client) {
	ctx := context.TODO()

	deployment := createDeployment(ctx, client, namePolicyViolationPath)
	result, err := client.DryRunDeployment(ctx, *deployment.ID, getParameters(namePolicyViolationPath))

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Validation Results:\n %s" + *prettify(result.Results))
}

func prettify(obj any) *string {
	bytes, _ := json.MarshalIndent(obj, "", "  ")
	result := string(bytes)
	return &result
}

// create the deployment with values
func createDeployment(ctx context.Context, client *Client, templatePath string) *generated.Deployment {
	name := "DryRunDeploymentTest"
	template := getTemplate(templatePath)
	deployment, err := client.CreateDeployment(ctx, generated.CreateDeployment{
		Name:           &name,
		SubscriptionID: &subscriptionId,
		ResourceGroup:  &resourceGroupName,
		Location:       &location,
		Template:       template,
	})

	if err != nil {
		log.Fatal("Failed to create deployment.")
	}

	return deployment
}

const namePolicyViolationPath string = "./test/data/namepolicy/failure/"
const quotaViolationPath string = "./test/data/quota/failure/"

func getParameters(path string) map[string]interface{} {
	paramsPath := filepath.Join(path, "parameters.json")
	parameters, err := utils.ReadJson(paramsPath)
	if err != nil {
		log.Printf("TestDryRun() could not read parameters")
	}
	return parameters
}

func getTemplate(path string) map[string]interface{} {
	fullPath := filepath.Join(path, "mainTemplate.json")
	template, err := utils.ReadJson(fullPath)
	if err != nil {
		log.Printf("couldn't read template")
	}
	return template
}
