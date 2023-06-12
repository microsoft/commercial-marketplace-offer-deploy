package operation

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/microsoft/commercial-marketplace-offer-deploy/test/fake"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type executeTestFakes struct {
	db         *gorm.DB
	notifyFake *fake.FakeNotifyFunc // underlying fake for notify
	sender     *fake.FakeMessageSender
	log        *log.Entry
	manager    *OperationManager
	repository Repository
	runFunc    OperationFunc
}

type operationExecuteTest struct {
	suite.Suite
	actionReturnsError bool
	fakes              executeTestFakes
	operation          *Operation
}

func TestExecutor(t *testing.T) {
	suite.Run(t, new(operationExecuteTest))
}

func (suite *operationExecuteTest) SetupTest() {
	suite.newFakes()

	suite.fakes.db.Save(&model.Deployment{})

	//create a new operation with all the underlying fakes to control it
	operation, err := suite.fakes.repository.New(sdk.OperationDeploy, func(i *model.InvokedOperation) error {
		i.Retries = 3
		i.Status = sdk.StatusNone.String()
		i.DeploymentId = 1
		return nil
	})
	suite.Require().NoError(err)

	suite.operation = operation
}

func (suite *operationExecuteTest) Test_Dispatch_Count_Matches_Attempts() {
	operation := suite.operation

	suite.Assert().Equal(uint(3), operation.Retries)
	suite.Assert().Equal(uint(0), operation.Attempts)

	suite.actionReturnsError = true
	operation.Task(NewOperationTask(OperationTaskOptions{
		Run: suite.fakes.runFunc,
	}))

	for !operation.AttemptsExceeded() {
		operation.Execute()
	}
	suite.Assert().Equal(operation.Attempts, uint(suite.fakes.sender.Count()+1))
}

func (suite *operationExecuteTest) Test_Notifications_For_Retries() {
	operation := suite.operation

	suite.Assert().Equal(uint(3), operation.Retries)
	suite.Assert().Equal(uint(0), operation.Attempts)

	suite.actionReturnsError = true
	operation.Task(NewOperationTask(OperationTaskOptions{
		Run: suite.fakes.runFunc,
	}))

	for !operation.AttemptsExceeded() {
		operation.Execute()
	}
	suite.Assert().Equal(8, suite.fakes.notifyFake.Count())

	for _, message := range suite.fakes.notifyFake.Messages() {
		bytes, err := json.MarshalIndent(message, "", "  ")
		suite.Require().NoError(err)

		suite.T().Logf("%+v", string(bytes))
	}
}

//region fakes

func (suite *operationExecuteTest) newFakes() {
	notifyFake := fake.NewFakeNotifyFunc(suite.T())

	db := fake.NewFakeDatabase(suite.T()).Instance()

	sender := fake.NewFakeMessageSender(suite.T())

	manager, err := NewManager(db, sender, notifyFake.Notify)
	suite.Require().NoError(err)

	repository, err := NewRepository(manager, nil)
	suite.Require().NoError(err)

	log.SetLevel(log.TraceLevel)

	suite.fakes = executeTestFakes{
		db:         fake.NewFakeDatabase(suite.T()).Instance(),
		notifyFake: notifyFake,
		sender:     sender,
		manager:    manager,
		repository: repository,
		runFunc: func(context ExecutionContext) error {
			suite.T().Log("operationAction called")
			if suite.actionReturnsError {
				return &RetriableError{Err: errors.New("fake error"), RetryAfter: 0}
			}
			return nil
		},
		log: log.WithFields(log.Fields{
			"test": "executorTest",
		}),
	}
}

//endregion fakes
