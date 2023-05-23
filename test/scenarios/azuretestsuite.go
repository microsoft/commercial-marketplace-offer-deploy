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
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const RequiredEnvVars = "TEST_AZURE_SUBSCRIPTION_ID,TEST_AZURE_LOCATION"

type AzureTestVars struct {
	SubscriptionId    string
	ResourceGroupName string
	Location          string
}

// AzureTestSuite is the base test suite for all Azure tests.
type AzureTestSuite struct {
	suite.Suite
	AzureVars AzureTestVars
}

func (suite *AzureTestSuite) SetupSuite() {
	suite.RequireEnvVars()

	suite.AzureVars = AzureTestVars{
		SubscriptionId:    os.Getenv("TEST_AZURE_SUBSCRIPTION_ID"),
		ResourceGroupName: "modm-test-scenario-" + suite.RandomString(5),
		Location:          os.Getenv("TEST_AZURE_LOCATION"),
	}

	suite.T().Logf("SetupSuite: %+v", suite.AzureVars)

	suite.CreateOrUpdateResourceGroup(suite.AzureVars.ResourceGroupName)
}

func (suite *AzureTestSuite) TearDownSuite() {
	suite.DeleteResourceGroup(suite.AzureVars.ResourceGroupName)
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

func (suite *AzureTestSuite) CreateOrUpdateResourceGroup(name string) {
	suite.T().Logf("Creating Resource Group [%s]", name)

	client := suite.newResourceGroupClient()
	_, err := client.CreateOrUpdate(
		context.Background(),
		name,
		armresources.ResourceGroup{
			Location: to.Ptr(suite.AzureVars.Location),
		},
		nil,
	)
	suite.Require().NoError(err)
}

func (suite *AzureTestSuite) DeleteResourceGroup(name string) {
	suite.T().Logf("Deleting Resource Group [%s]", name)

	client := suite.newResourceGroupClient()
	_, err := client.BeginDelete(
		context.Background(),
		name,
		nil,
	)
	suite.Require().NoError(err)
}

func (suite *AzureTestSuite) ReadJsonFile(dirPath string, fileName string) map[string]any {
	fullPath := filepath.Join(dirPath, fileName)
	template, err := utils.ReadJson(fullPath)

	require.NoError(suite.T(), err)

	return template
}

func (suite *AzureTestSuite) ToJson(i any) string {
	bytes, err := json.Marshal(i)
	suite.Require().NoError(err)
	return string(bytes)
}

func (suite *AzureTestSuite) newResourceGroupClient() *armresources.ResourceGroupsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	suite.Require().NoError(err)

	client, err := armresources.NewResourceGroupsClient(suite.AzureVars.SubscriptionId, cred, nil)
	suite.Require().NoError(err)

	return client
}
