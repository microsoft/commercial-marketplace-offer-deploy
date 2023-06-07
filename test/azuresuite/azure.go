package azuresuite

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
	"github.com/stretchr/testify/suite"
)

const RequiredEnvVars = "TEST_AZURE_SUBSCRIPTION_ID,TEST_AZURE_LOCATION"

type AzureTestSettings struct {
	SubscriptionId    string
	ResourceGroupName string
	Location          string
}

// AzureTestSuite is the base test suite for all Azure tests.
type AzureTestSuite struct {
	suite.Suite

	// store settings by test name
	TestSettings map[string]AzureTestSettings
}

func (suite *AzureTestSuite) AddSettings(settings AzureTestSettings) string {
	testName := suite.T().Name()

	suite.T().Logf("Adding settings for %s", testName)

	if suite.TestSettings == nil {
		suite.TestSettings = make(map[string]AzureTestSettings)
	}
	suite.TestSettings[testName] = settings

	return testName
}

// creates new azure settings based on environment variables
func (suite *AzureTestSuite) NewSettings() (AzureTestSettings, string) {
	settings := AzureTestSettings{
		SubscriptionId:    os.Getenv("TEST_AZURE_SUBSCRIPTION_ID"),
		ResourceGroupName: "modmtest-scenario-" + suite.RandomString(5),
		Location:          os.Getenv("TEST_AZURE_LOCATION"),
	}

	// override resource group name if set
	if os.Getenv("TEST_AZURE_RESOURCE_GROUP") != "" {
		settings.ResourceGroupName = os.Getenv("TEST_AZURE_RESOURCE_GROUP")
	}

	key := suite.AddSettings(settings)

	return settings, key
}

// get settings using the built in active test name
func (suite *AzureTestSuite) Settings() AzureTestSettings {
	if s, ok := suite.TestSettings[suite.T().Name()]; ok {
		return s
	}
	return AzureTestSettings{}
}

func (suite *AzureTestSuite) SettingsByName(testName string) (AzureTestSettings, bool) {
	if s, ok := suite.TestSettings[testName]; ok {
		return s, true
	}
	return AzureTestSettings{}, false
}

func (suite *AzureTestSuite) SetupSuite() {
	suite.RequireEnvVars()
}

func (suite *AzureTestSuite) TearDownSuite() {
	for _, s := range suite.TestSettings {
		suite.DeleteResourceGroup(s)
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

func (suite *AzureTestSuite) CreateOrUpdateResourceGroup(settings AzureTestSettings) {
	suite.T().Logf("Creating Resource Group [%s]", settings.ResourceGroupName)

	client := suite.newResourceGroupClient(settings.SubscriptionId)
	_, err := client.CreateOrUpdate(
		context.Background(),
		settings.ResourceGroupName,
		armresources.ResourceGroup{
			Location: to.Ptr(strings.ToLower(settings.Location)),
		},
		nil,
	)
	suite.Require().NoError(err)
}

func (suite *AzureTestSuite) DeleteResourceGroup(settings AzureTestSettings) {
	suite.T().Logf("Deleting Resource Group [%s]", settings.ResourceGroupName)

	client := suite.newResourceGroupClient(settings.SubscriptionId)

	if res, err := client.CheckExistence(context.Background(), settings.ResourceGroupName, nil); err == nil && res.Success {
		_, err := client.BeginDelete(
			context.Background(),
			settings.ResourceGroupName,
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

func (suite *AzureTestSuite) GetCredential() *azidentity.DefaultAzureCredential {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	suite.Require().NoError(err)

	return cred
}
