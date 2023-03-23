package sdk

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/joho/godotenv"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
)

var subscriptionId string
var resourceGroupName string
var location string
var endpoint string = "http://localhost:8080"

func SetupDryTest() {
	// log.Println("Test setup beginning")
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Println("Cannot load environment variables from .env")
	// }

	subscriptionId = "31e9f9a0-9fd2-4294-a0a3-0101246d9700"
	resourceGroupName = "MODMTest"
	location = "eastus"

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

	deployment, err := client.CreateDeployment(ctx, internal.CreateDeployment{})

	if err != nil {
		log.Fatal("Failed to create deployment.")
	}

	parameters := getParameters()
	client.DryRunDeployment(ctx, *deployment.ID, parameters)
}

func getParameters() map[string]interface{} {
	paramsPath := "./test/data/namepolicy/success/parameters.json"
	parameters, err := utils.ReadJson(paramsPath)
	if err != nil {
		log.Printf("TestDryRun() could not read parameters")
	}
	return parameters
}
