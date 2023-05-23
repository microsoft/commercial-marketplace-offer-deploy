package test_scenarios_dryrun

import (
	"context"
	"testing"

	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/microsoft/commercial-marketplace-offer-deploy/test/scenarios"
	"github.com/stretchr/testify/suite"
)

const nameConflictTestName1 = "TestNameConflictTestSuite/Test_Should_Succeed_In_Same_RG"
const nameConflictTestName2 = "TestNameConflictTestSuite/Test_Should_Fail_In_Different_Resource_Group"

type NameConflictTestSuite struct {
	DryRunTestSuite
}

//region setup

func TestNameConflictTestSuite(t *testing.T) {
	suite.Run(t, new(NameConflictTestSuite))
}

func (suite *NameConflictTestSuite) SetupSuite() {
	suite.DryRunTestSuite.SetupSuite()

	suite.TestDataDirPath = "./testdata/nameconflict"

	defaultVars := suite.GetVariables(scenarios.DefaultTestVariablesKey)
	suite.AddVariables(nameConflictTestName2, scenarios.AzureTestVariables{
		SubscriptionId:    defaultVars.SubscriptionId,
		Location:          defaultVars.Location,
		ResourceGroupName: "modmtest-scenario-" + suite.RandomString(10),
	})
}

func (suite *NameConflictTestSuite) TearDownSuite() {
	suite.DryRunTestSuite.TearDownSuite()
}

func (suite *NameConflictTestSuite) SetupTest() {
	suite.DryRunTestSuite.SetupTest()
	suite.deployTemplate()
}

func (suite *NameConflictTestSuite) TearDownTest() {
	suite.DeleteResourceGroup(suite.GetVariables(nameConflictTestName2))
}

func (suite *NameConflictTestSuite) deployTemplate() {
	suite.Deployment.Params["parameters"].(map[string]any)["name"].(map[string]any)["value"] = "modmteststor0" + suite.RandomString(5)

	_, err := deployment.Create(context.Background(), suite.Deployment)
	suite.Require().NoError(err)
}

//endregion setup

//region tests

func (suite *NameConflictTestSuite) Test_Should_Succeed_In_Same_RG() {
	ctx := context.TODO()

	azureDeployment := suite.Deployment
	result, err := deployment.DryRun(ctx, &azureDeployment)
	suite.Assert().NoError(err)

	suite.T().Logf("result: %+v", suite.ToJson(result))

	suite.Assert().Equal(sdk.StatusSuccess.String(), result.Status)
}

func (suite *NameConflictTestSuite) Test_Should_Fail_In_Different_Resource_Group() {
	ctx := context.TODO()

	result, err := deployment.DryRun(ctx, &suite.Deployment)
	suite.Assert().NoError(err)

	suite.T().Logf("result: %+v", suite.ToJson(result))

	suite.Assert().Equal(sdk.StatusFailed.String(), result.Status)
}

//endregion tests
