package notification

import (
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
}

//tests

func (suite *serviceTestSuite) TestStart() {
	service := NewStageNotificationService(suite.pump, suite.handlerFactory)
	service.Start()
}

// handler factory method on suite
func (suite *serviceTestSuite) fakeHandlerFactory() StageNotificationHandlerFactoryFunc {
	return func() (NotificationHandler[model.StageNotification], error) {
		return &fakeHandler{t: suite.T()}, nil
	}
}

func (suite *serviceTestSuite) fakePump() NotificationPump[model.StageNotification] {
	return &fakePump{t: suite.T()}
}

//region fakes

type fakePump struct {
	t        *testing.T
	receiver NotificationPumpReceiveFunc[model.StageNotification]
}

func (p *fakePump) Start() {
	p.t.Log("fake pump started")

	time.Sleep(5 * time.Second)
	p.receiver(&model.StageNotification{})
}

func (p *fakePump) Stop() {
	p.t.Log("fake pump stopped")
}

func (p *fakePump) SetReceiver(receiver NotificationPumpReceiveFunc[model.StageNotification]) {
	p.t.Log("fake pump receiver set")
	p.receiver = receiver
}

type fakeHandler struct {
	t *testing.T
}

func (h *fakeHandler) Handle(context *NotificationHandlerContext[model.StageNotification]) {
	h.t.Log("fake handler called")
	context.Done(NotificationHandlerResult[model.StageNotification]{
		Notification: context.Notification,
	})
}

//endregion fakes
