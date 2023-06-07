package notification

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/hook"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/deployment"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/microsoft/commercial-marketplace-offer-deploy/test/azuresuite"
	"github.com/stretchr/testify/suite"
)

type handlerTestSuite struct {
	azuresuite.AzureTestSuite
	notifyFunc        hook.NotifyFunc
	deploymentsClient *armresources.DeploymentsClient
	correlationId     uuid.UUID
	settingsName      string
	notification      *model.StageNotification
}

func TestStageNotificationHandler(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}

func (suite *handlerTestSuite) SetupSuite() {
	suite.AzureTestSuite.SetupSuite()

	settings, key := suite.NewSettings()
	suite.settingsName = key

	suite.notifyFunc = func(ctx context.Context, message *sdk.EventHookMessage) (uuid.UUID, error) {
		suite.T().Logf("NotifyFunc called with message: %v", message)
		return uuid.New(), nil
	}

	suite.setDeploymentsClient(settings)
	suite.CreateOrUpdateResourceGroup(settings)

	// now create a deployment
	suite.deployTemplate(settings)

	// this id comes from testdata/ value in the template.json file
	stageId := uuid.MustParse("31e9f9a0-9fd2-4294-a0a3-0101246d9700")

	suite.notification = &model.StageNotification{
		ResourceGroupName: settings.ResourceGroupName,
		CorrelationId:     suite.correlationId,
		Entries: []model.StageNotificationEntry{
			{
				StageId: stageId,
				Message: sdk.EventHookMessage{
					Type:   string(sdk.EventTypeStageStarted),
					Status: string(sdk.StatusRunning),
					Data: sdk.DeploymentEventData{
						DeploymentId: 1,
						StageId:      &stageId,
					},
				},
			},
		},
	}
}

//region tests

func (suite *handlerTestSuite) Test_StageNotificationHandler_getAzureDeployments() {
	settings, ok := suite.SettingsByName(suite.settingsName)
	suite.Require().True(ok)

	handler := &stageNotificationHandler{
		notify:            suite.notifyFunc,
		deploymentsClient: suite.deploymentsClient,
	}

	result, err := handler.getAzureDeploymentResources(context.Background(), &model.StageNotification{
		ResourceGroupName: settings.ResourceGroupName,
		CorrelationId:     suite.correlationId,
	})
	suite.Assert().NoError(err)
	suite.Assert().Len(result, 2)
}

func (suite *handlerTestSuite) Test_StageNotificationHandler_Handle() {
	handler := &stageNotificationHandler{
		notify:            suite.notifyFunc,
		deploymentsClient: suite.deploymentsClient,
	}

	context := NewNotificationHandlerContext[model.StageNotification](context.Background(), suite.notification)

	channel := context.Channel()

	go handler.Handle(context)
	<-channel

	suite.Assert().True(context.Notification.IsDone)
}

//endregion tests

func (suite *handlerTestSuite) setDeploymentsClient(settings azuresuite.AzureTestSettings) {
	credential := suite.GetCredential()
	deploymentsClient, err := armresources.NewDeploymentsClient(settings.SubscriptionId, credential, nil)
	suite.Require().NoError(err)

	suite.deploymentsClient = deploymentsClient
}

// suite method that creates an azure deployment
func (suite *handlerTestSuite) deployTemplate(settings azuresuite.AzureTestSettings) {
	suite.T().Logf("Deploying template to: \nResource Group: %s", settings.ResourceGroupName)

	testdir := "./testdata/testdeployment"

	azureDeployment := deployment.AzureDeployment{
		SubscriptionId:    settings.SubscriptionId,
		ResourceGroupName: settings.ResourceGroupName,
		Location:          settings.Location,
		DeploymentName:    "test-deploy-" + suite.RandomString(5),
		Template:          suite.ReadJsonFile(testdir, "template.json"),
		Params:            suite.ReadJsonFile(testdir, "parameters.json"),
	}

	deployer, err := deployment.NewDeployer(deployment.DeploymentTypeARM, settings.SubscriptionId)
	suite.Require().NoError(err)

	ctx := context.Background()
	begin, err := deployer.Begin(ctx, azureDeployment)
	suite.Require().NoError(err)
	suite.Require().NotNil(begin.CorrelationID)

	id := *begin.CorrelationID
	suite.T().Logf("Beginning template deployment to setup test  [%v]", id)
	suite.T().Log("This will take a minute. Make sure the timeout of the test is long enough to wait for the deployment to complete.")

	suite.correlationId = uuid.MustParse(id)

	_, err = deployer.Wait(ctx, &begin.ResumeToken)
	suite.Require().NoError(err)
}
