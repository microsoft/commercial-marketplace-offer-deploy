package sdk

import (
	"context"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/generated"
)

var subscriptionId string
var resourceGroupName string
var location string
var endpoint string = "http://localhost:8080"

func SetupDryTest() {
	subscriptionId = "31e9f9a0-9fd2-4294-a0a3-0101246d9700"
	resourceGroupName = "aMODMTestb"
	location = "eastus"

	if exists, _ := utils.DoesResourceGroupExist(subscriptionId, resourceGroupName, location); exists {
		return
	}

	utils.SetupResourceGroup(subscriptionId, resourceGroupName, location)
	utils.DeployPolicyDefinition(subscriptionId)
	utils.DeployPolicy(subscriptionId, resourceGroupName)

	//exitVal := m.Run()
	//log.Println("Cleaning up resources after the tests here")

	//os.Exit(exitVal)
}

func TestDryRun(t *testing.T) {
	//err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Println("Cannot load environment variables from .env")
	// }
	SetupDryTest()

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
	client.DryRunDeployment(ctx, *deployment.ID, parameters)
}

// create the deployment with values
func createDeployment(ctx context.Context, client *Client) *generated.Deployment {
	name := "Test"
	template := getTemplate()

	log.Printf("creating deployment with template %v", template)

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

func getParameters() map[string]interface{} {
	paramsPath := "./test/data/namepolicy/success/parameters.json"
	parameters, err := utils.ReadJson(paramsPath)
	if err != nil {
		log.Printf("TestDryRun() could not read parameters")
	}
	return parameters
}

func getTemplate() map[string]interface{} {
	path := "../test/deployment/resourcename/failure/mainTemplate.json"
	template, err := utils.ReadJson(path)
	if err != nil {
		log.Printf("couldn't read template")
	}
	return template
}
