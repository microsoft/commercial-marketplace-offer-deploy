package scenarios

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

const RequiredEnvVars = "TEST_AZURE_SUBSCRIPTION_ID,TEST_AZURE_LOCATION"
const DefaultTestVariablesKey = "default"

type AzureTestVariables struct {
	SubscriptionId    string
	ResourceGroupName string
	Location          string
}

// AzureTestSuite is the base test suite for all Azure tests.
type AzureTestSuite struct {
	suite.Suite
	Variables map[string]AzureTestVariables
}

func (suite *AzureTestSuite) AddVariables(testName string, vars AzureTestVariables) {
	if suite.Variables == nil {
		suite.Variables = make(map[string]AzureTestVariables)
	}
	suite.Variables[testName] = vars
}

func (suite *AzureTestSuite) GetVariables(testName string) AzureTestVariables {
	if vars, ok := suite.Variables[testName]; ok {
		return vars
	}
	return suite.Variables[DefaultTestVariablesKey]
}

func (suite *AzureTestSuite) SetupSuite() {

	log.SetLevel(log.DebugLevel)

	suite.RequireEnvVars()

	defaultVariables := AzureTestVariables{
		SubscriptionId:    os.Getenv("TEST_AZURE_SUBSCRIPTION_ID"),
		ResourceGroupName: "modmtest-scenario-" + suite.RandomString(5),
		Location:          os.Getenv("TEST_AZURE_LOCATION"),
	}

	// override resource group name if set
	if os.Getenv("TEST_AZURE_RESOURCE_GROUP") != "" {
		defaultVariables.ResourceGroupName = os.Getenv("TEST_AZURE_RESOURCE_GROUP")
	}

	suite.AddVariables(DefaultTestVariablesKey, defaultVariables)
	suite.T().Logf("SetupSuite: %+v", suite.Variables)
}

func (suite *AzureTestSuite) TearDownSuite() {
	for _, vars := range suite.Variables {
		suite.DeleteResourceGroup(vars)
	}
}

// RequireEnvVars checks that the required environment variables are set as part of the setup.
func (suite *AzureTestSuite) RequireEnvVars() {
	varkeys := strings.Split(RequiredEnvVars, ",")
	for _, varkey := range varkeys {
		suite.Require().NotEmpty(os.Getenv(varkey), "Missing environment variable: "+varkey)
	}
}

func (suite *AzureTestSuite) RandomString(length uint) string {
	base := uuid.New()
	baseString := strings.ReplaceAll(base.String(), "-", "")
	return baseString[0:length]
}

func (suite *AzureTestSuite) CreateOrUpdateResourceGroup(variables AzureTestVariables) {
	suite.T().Logf("Creating Resource Group [%s]", variables.ResourceGroupName)

	client := suite.newResourceGroupClient(variables.SubscriptionId)
	_, err := client.CreateOrUpdate(
		context.Background(),
		variables.ResourceGroupName,
		armresources.ResourceGroup{
			Location: to.Ptr(strings.ToLower(variables.Location)),
		},
		nil,
	)
	suite.Require().NoError(err)
}

func (suite *AzureTestSuite) DeleteResourceGroup(variables AzureTestVariables) {
	suite.T().Logf("Deleting Resource Group [%s]", variables.ResourceGroupName)

	client := suite.newResourceGroupClient(variables.SubscriptionId)

	if res, err := client.CheckExistence(context.Background(), variables.ResourceGroupName, nil); err == nil && res.Success {
		_, err := client.BeginDelete(
			context.Background(),
			variables.ResourceGroupName,
			nil,
		)
		suite.Require().NoError(err)
	}
}

func (suite *AzureTestSuite) ReadJsonFile(dirPath string, fileName string) map[string]any {
	fullPath := filepath.Join(dirPath, fileName)
	contents, err := utils.ReadJson(fullPath)

	if err != nil {
		suite.T().Log(err)
		return make(map[string]any)
	}

	return contents
}

func (suite *AzureTestSuite) ToJson(i any) string {
	bytes, err := json.Marshal(i)
	suite.Require().NoError(err)
	return string(bytes)
}

func (suite *AzureTestSuite) newResourceGroupClient(subscriptionId string) *armresources.ResourceGroupsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	suite.Require().NoError(err)

	client, err := armresources.NewResourceGroupsClient(subscriptionId, cred, nil)
	suite.Require().NoError(err)

	return client
}
