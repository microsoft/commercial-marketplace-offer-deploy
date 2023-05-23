package test_scenarios_dryrun

import (
	"context"
	"testing"

	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/stretchr/testify/suite"
)

type DirectTemplateParamsTestSuite struct {
	DryRunTestSuite
}

//region setup

func TestDirectTemplateParamsTestSuite(t *testing.T) {
	suite.Run(t, new(DirectTemplateParamsTestSuite))
}

func (suite *DirectTemplateParamsTestSuite) SetupSuite() {
	suite.DryRunTestSuite.SetupSuite()
	suite.TestDataDirPath = "./testdata/directtemplateparams"
}

func (suite *DirectTemplateParamsTestSuite) TearDownSuite() {
	suite.DryRunTestSuite.TearDownSuite()
}

func (suite *DirectTemplateParamsTestSuite) SetupTest() {
	suite.DryRunTestSuite.SetupTest()
	suite.Deployment.Params = map[string]interface{}{
		"name": "teststore0" + suite.RandomString(5),
	}
}

//endregion setup

//region tests

func (suite *DirectTemplateParamsTestSuite) Test_DryRun_Supports_Direct_Params() {
	ctx := context.TODO()

	azureDeployment := suite.Deployment
	result, err := deployment.DryRun(ctx, &azureDeployment)
	suite.Assert().NoError(err)

	suite.T().Logf("result: %+v", suite.ToJson(result))
	suite.Assert().Equal(sdk.StatusSuccess.String(), result.Status)
}

//endregion tests
