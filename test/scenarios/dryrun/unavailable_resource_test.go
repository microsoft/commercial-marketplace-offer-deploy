package test_scenarios_dryrun

import (
	"context"
	"testing"

	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/stretchr/testify/suite"
)

type UnavailableResourceTestSuite struct {
	DryRunTestSuite
}

//region setup

func TestUnavailableResourceTestSuite(t *testing.T) {
	suite.Run(t, new(UnavailableResourceTestSuite))
}

func (suite *UnavailableResourceTestSuite) SetupSuite() {
	suite.DryRunTestSuite.SetupSuite()
	suite.TestDataDirPath = "./testdata/unavailableresource"
}

//endregion setup

//region tests

func (suite *UnavailableResourceTestSuite) Test_Should_Fail_If_Resource_Sku_Not_In_Region() {
	ctx := context.TODO()
	result, err := deployment.DryRun(ctx, &suite.Deployment)
	suite.Assert().NoError(err)

	suite.T().Logf("result: %+v", suite.ToJson(result))

	suite.Assert().Equal(sdk.StatusFailed.String(), result.Status)
}

//endregion tests
