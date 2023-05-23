package test_scenarios_dryrun

import (
	"context"
	"testing"

	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/stretchr/testify/suite"
)

type NameConflictTestSuite struct {
	DryRunTestSuite
	ResourceName string
}

//region setup

func TestNameConflictTestSuite(t *testing.T) {
	suite.Run(t, new(NameConflictTestSuite))
}

func (suite *NameConflictTestSuite) SetupSuite() {
	suite.DryRunTestSuite.SetupSuite()
	suite.TestDataDirPath = "./testdata/nameconflict"
}

func (suite *NameConflictTestSuite) SetupTest() {
	suite.DryRunTestSuite.SetupTest()

	suite.ResourceName = "modmtest0" + suite.RandomString(10)

	suite.Deployment.DeploymentName = "deploy-" + suite.ResourceName
	suite.Deployment.Params["parameters"].(map[string]any)["name"].(map[string]any)["value"] = suite.ResourceName

	suite.T().Logf(" - Storage Account Name: %s", suite.ResourceName)
	// deploy so we can run dry run against the deployed storage account
	_, err := deployment.Create(context.Background(), suite.Deployment)
	suite.Require().NoError(err)
}

//endregion setup

//region tests

func (suite *NameConflictTestSuite) Test_Should_Catch_Service_Name_Conflict() {
	ctx := context.TODO()

	result, err := deployment.DryRun(ctx, &suite.Deployment)
	suite.Assert().NoError(err)

	suite.T().Logf("result: %+v", suite.ToJson(result))

	suite.Assert().Equal(sdk.StatusFailed.String(), result.Status)
}

//endregion tests
