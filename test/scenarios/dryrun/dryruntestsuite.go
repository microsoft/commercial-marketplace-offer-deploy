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

func (suite *DryRunTestSuite) newAzureDeployment() deployment.AzureDeployment {
	suite.T().Log(" - Creating new AzureDeployment")

	d := deployment.AzureDeployment{
		SubscriptionId:    suite.AzureVars.SubscriptionId,
		Location:          suite.AzureVars.Location,
		ResourceGroupName: suite.AzureVars.ResourceGroupName,
		DeploymentName:    "DryRunTest",
		DeploymentType:    deployment.AzureResourceManager,
		Template:          suite.readJsonFile("template"),
		Params:            suite.readJsonFile("parameters"),
	}
	suite.T().Logf(" - %+v", d)

	return d
}

func (suite *DryRunTestSuite) readJsonFile(name string) map[string]any {
	return suite.AzureTestSuite.ReadJsonFile(suite.TestDataDirPath, fmt.Sprintf("%s.json", name))
}

//endregion helpers

func (suite *DryRunTestSuite) SetupSuite() {
	suite.AzureTestSuite.SetupSuite()

	t := suite.T()

	t.Log("SetupSuite")
	suite.TestDataDirPath = "./testdata/unavailableresource"

}

func (suite *DryRunTestSuite) TearDownSuite() {
	suite.T().Log("TearDown")
	suite.AzureTestSuite.TearDownSuite()
}

func (suite *DryRunTestSuite) SetupTest() {
	suite.T().Logf("Setup Test - [%s]", suite.T().Name())
	suite.Deployment = suite.newAzureDeployment()
}

func (suite *DryRunTestSuite) TearDownTest() {
	suite.T().Log("TearDown Test")
}
