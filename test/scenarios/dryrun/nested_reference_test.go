package test_scenarios_dryrun

import (
	"context"
	"testing"

	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/microsoft/commercial-marketplace-offer-deploy/test/scenarios"
	"github.com/stretchr/testify/suite"
)

type NestedReferenceTestSuite struct {
	DryRunTestSuite
}

//region setup

func TestNestedReferenceTestSuite(t *testing.T) {
	suite.Run(t, new(NestedReferenceTestSuite))
}

func (suite *NestedReferenceTestSuite) SetupSuite() {
	suite.DryRunTestSuite.SetupSuite()
	suite.TestDataDirPath = "./testdata/nestedreference"

	suite.AddVariables("resourceGroupWithExistingStorageAccount", scenarios.AzureTestVariables{
		SubscriptionId:    suite.DefaultVariables().SubscriptionId,
		Location:          suite.DefaultVariables().Location,
		ResourceGroupName: "modmtest-scenario-" + suite.RandomString(10),
	})
}

func (suite *NestedReferenceTestSuite) TearDownSuite() {
	suite.DryRunTestSuite.TearDownSuite()
}

func (suite *NestedReferenceTestSuite) SetupTest() {
	suite.DryRunTestSuite.SetupTest()

	existingStorageAccountName := suite.setupExistingStorageAccount()

	// storage account 2 uses the output variable value from storage account 1, which is coming from name2 parameter
	suite.Deployment.Params = map[string]interface{}{
		"name1": "teststore0" + suite.RandomString(5),
		"name2": existingStorageAccountName,
	}
}

func (suite *NestedReferenceTestSuite) setupExistingStorageAccount() string {
	storageAccountName := "existmodm0" + suite.RandomString(10)
	vars := suite.GetVariables("resourceGroupWithExistingStorageAccount")
	suite.CreateOrUpdateResourceGroup(vars)

	azureDeployment := suite.NewDeployment(variablesKey3)
	azureDeployment.Template = suite.readJsonFile("template.existing")
	azureDeployment.Params = map[string]interface{}{
		"name": storageAccountName,
	}

	_, err := deployment.Create(context.Background(), azureDeployment)
	suite.Require().NoError(err)

	suite.T().Logf("created storage account [%s] in [%s]", storageAccountName, vars.ResourceGroupName)

	return storageAccountName
}

//endregion setup

//region tests

func (suite *NestedReferenceTestSuite) Test_DryRun_With_Nested_Reference_To_Name_Conflict() {
	ctx := context.TODO()

	azureDeployment := suite.Deployment
	result, err := deployment.DryRun(ctx, &azureDeployment)
	suite.Assert().NoError(err)

	suite.T().Logf("result: %+v", suite.ToJson(result))
	suite.Assert().Equal(sdk.StatusFailed.String(), result.Status)
}

//endregion tests
