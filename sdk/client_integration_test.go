package sdk

import (
	"context"
	"log"
	"os"
	"testing"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/stretchr/testify/require"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	"github.com/joho/godotenv"
)

var subscriptionId string
var resourceGroupName string
var location string
var endpoint string

func SetupDryTest() {
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

	client, err := NewClient("http://localhost:8080", cred, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	paramsPath := "./test/data/namepolicy/success/parameters.json"
	parameters, err := utils.ReadJson(paramsPath)
	if err != nil {
		t.Errorf("TestDryRun() could not read parameters")
	} 

	client.DryRunDeployment(context.TODO(), 1, parameters)
}
