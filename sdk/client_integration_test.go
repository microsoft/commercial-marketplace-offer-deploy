package sdk

import (
	"context"
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
	ctx := context.TODO()

	if err != nil {
		log.Print("Client construction failed.")
	}

	deployment := createDeployment(ctx, client)

	parameters := getParameters()
	deploymentId := *deployment.ID

	log.Printf("Deployment id: %d", deploymentId)
	client.DryRunDeployment(ctx, deploymentId, parameters)
}

// create the deployment with values
func createDeployment(ctx context.Context, client *Client) *generated.Deployment {
	name := "DryRunDeploymentTest"
	template := getTemplate()
	deployment, err := client.CreateDeployment(ctx, generated.CreateDeployment{
		Name:           &name,
		SubscriptionID: &subscriptionId,
		ResourceGroup:  &resourceGroupName,
		Location:       &location,
		Template:       template,
	})

	log.Printf("Deployment: %d", *deployment.ID)

	if err != nil {
		log.Fatal("Failed to create deployment.")
	}

	return deployment
}

var testDataPath string = "./test/data/namepolicy/failure/"

func getParameters() map[string]interface{} {
	paramsPath := filepath.Join(testDataPath, "parameters.json")
	parameters, err := utils.ReadJson(paramsPath)
	if err != nil {
		log.Printf("TestDryRun() could not read parameters")
	}
	return parameters
}

func getTemplate() map[string]interface{} {
	path := filepath.Join(testDataPath, "mainTemplate.json")
	template, err := utils.ReadJson(path)
	if err != nil {
		log.Printf("couldn't read template")
	}
	return template
}
