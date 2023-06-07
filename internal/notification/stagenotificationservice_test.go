package notification

import (
	"math/rand"
	"testing"
	"time"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type serviceTestSuite struct {
	suite.Suite

	db             *gorm.DB
	handlerFactory StageNotificationHandlerFactoryFunc
	pump           NotificationPump[model.StageNotification]

	notification *model.StageNotification
}

// entry point for running the test suite
func TestStageNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (suite *serviceTestSuite) SetupSuite() {
	suite.db = data.NewDatabase(&data.DatabaseOptions{UseInMemory: true}).Instance()
	suite.handlerFactory = suite.fakeHandlerFactory()
	suite.pump = suite.fakePump()
}

func (suite *serviceTestSuite) SetupTest() {
	suite.notification = &model.StageNotification{
		Model: gorm.Model{
			ID: uint(rand.Int()),
		},
		OperationId:       [16]byte{},
		CorrelationId:     [16]byte{},
		ResourceGroupName: "",
		Entries:           []model.StageNotificationEntry{},
		Done:              false,
	}
	suite.T().Logf("Test setup with ID: %d", suite.notification.ID)
}

//tests

func (suite *serviceTestSuite) Test_Start_Stop_Is_Wired_Up() {
	service := NewStageNotificationService(suite.pump, suite.handlerFactory)
	service.Start()
	service.Stop()
}

// handler factory method on suite
func (suite *serviceTestSuite) fakeHandlerFactory() StageNotificationHandlerFactoryFunc {
	return func() (NotificationHandler[model.StageNotification], error) {
		return &fakeHandler{t: suite.T()}, nil
	}
}

func (suite *serviceTestSuite) fakePump() NotificationPump[model.StageNotification] {
	return &fakePump{suite: suite}
}

//region fakes

type fakePump struct {
	suite    *serviceTestSuite
	receiver NotificationPumpReceiveFunc[model.StageNotification]
}

func (p *fakePump) Start() {
	p.suite.T().Log("fake pump started")

	time.Sleep(5 * time.Second)
	p.receiver(p.suite.notification)
}

func (p *fakePump) Stop() {
	p.suite.T().Log("fake pump stopped")
}

func (p *fakePump) SetReceiver(receiver NotificationPumpReceiveFunc[model.StageNotification]) {
	p.suite.T().Log("fake pump receiver set")
	p.receiver = receiver
}

type fakeHandler struct {
	t *testing.T
}

func (h *fakeHandler) Handle(context *NotificationHandlerContext[model.StageNotification]) {
	h.t.Logf("fake handler called with ID: %d", context.Notification.ID)
	context.Done(NotificationHandlerResult[model.StageNotification]{
		Notification: context.Notification,
	})
}

//endregion fakes
