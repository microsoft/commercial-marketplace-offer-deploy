package test_scenarios_dryrun

import (
	"fmt"

	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/test/scenarios"
)

type DryRunTestSuite struct {
	scenarios.AzureTestSuite

	Deployment      deployment.AzureDeployment
	TestDataDirPath string
}

//region helpers

func (suite *DryRunTestSuite) NewDeployment(variablesKey string) deployment.AzureDeployment {
	suite.T().Log(" - Creating AzureDeployment")

	variables := suite.GetVariables(variablesKey)

	d := deployment.AzureDeployment{
		SubscriptionId:    variables.SubscriptionId,
		Location:          variables.Location,
		ResourceGroupName: variables.ResourceGroupName,
		DeploymentName:    "DryRunTest-" + suite.RandomString(5),
		DeploymentType:    deployment.AzureResourceManager,
		Template:          suite.readJsonFile("template"),
		Params:            suite.readJsonFile("parameters"),
	}
	return d
}

func (suite *DryRunTestSuite) readJsonFile(name string) map[string]any {
	return suite.AzureTestSuite.ReadJsonFile(suite.TestDataDirPath, fmt.Sprintf("%s.json", name))
}

//endregion helpers

func (suite *DryRunTestSuite) SetupSuite() {
	suite.AzureTestSuite.SetupSuite()
}

func (suite *DryRunTestSuite) TearDownSuite() {
	suite.T().Log("TearDown")
	suite.AzureTestSuite.TearDownSuite()
}

func (suite *DryRunTestSuite) SetupTest() {
	testName := suite.T().Name()

	vars := suite.GetVariables(testName)
	suite.CreateOrUpdateResourceGroup(vars)

	suite.T().Logf("Setup Test - [%s]", testName)
	suite.Deployment = suite.NewDeployment(testName)
}

func (suite *DryRunTestSuite) TearDownTest() {
	testName := suite.T().Name()
	suite.T().Log("TearDown Test")

	vars := suite.GetVariables(testName)
	suite.DeleteResourceGroup(vars)
}
