package handlers

import (
	"testing"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model/operation"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/microsoft/commercial-marketplace-offer-deploy/test/fake"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type operationsTestSuite struct {
	suite.Suite
	db *gorm.DB

	hookService  *fake.FakeHookService
	sender       *fake.FakeMessageSender
	service      *operation.OperationService
	funcProvider operation.OperationFuncProvider

	invokedOperation model.InvokedOperation
}

func TestOperationMessageHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(operationsTestSuite))
}

func (suite *operationsTestSuite) SetupSuite() {
	suite.funcProvider = suite.newFakeOperationFuncProvider()
	suite.db = data.NewDatabase(&data.DatabaseOptions{UseInMemory: true}).Instance()
	suite.hookService = fake.NewFakeHookService(suite.T())
	suite.sender = fake.NewFakeMessageSender(suite.T())

	service, err := operation.NewService(suite.db, suite.sender, suite.hookService.Notify)
	suite.Assert().NoError(err)

	suite.service = service

	suite.invokedOperation = model.InvokedOperation{
		BaseWithGuidPrimaryKey: model.BaseWithGuidPrimaryKey{
			ID: uuid.MustParse("53a65e21-6a7a-4994-8ff9-f51c610b6067"),
		},
		Name:         "deploy",
		Retries:      3,
		Attempts:     0,
		Parameters:   map[string]interface{}{},
		DeploymentId: 1,
	}
}

func (suite *operationsTestSuite) SetupTest() {
	suite.db.Save(&suite.invokedOperation)
}

func (suite *operationsTestSuite) Test_Handle_Completes_Operation_When_No_Error_In_Operation() {
	operationFactory, err := operation.NewRepository(suite.service, suite.funcProvider)
	suite.Assert().NoError(err)

	handler := &operationMessageHandler{
		operationFactory: operationFactory,
	}

	handler.Handle(&messaging.ExecuteInvokedOperation{
		OperationId: suite.invokedOperation.ID,
	}, messaging.MessageHandlerContext{
		ReceivedMessage: nil,
	})

	result := &model.InvokedOperation{}
	suite.db.First(result, suite.invokedOperation.ID)

	suite.Assert().Equal(uint(1), result.Attempts)
	suite.Assert().Equal(sdk.StatusSuccess.String(), result.Status)
	suite.Assert().True(result.IsCompleted())
}

func (suite *operationsTestSuite) newFakeOperationFuncProvider() operation.OperationFuncProvider {
	return &fakeOperationFuncProvider{
		t: suite.T(),
	}
}

type fakeOperationFuncProvider struct {
	t *testing.T
}

func (f *fakeOperationFuncProvider) Get(operationType sdk.OperationType) (operation.OperationFunc, error) {
	return func(context *operation.ExecutionContext) error {
		f.t.Logf("Executing operation %s", operationType)
		return nil
	}, nil
}
