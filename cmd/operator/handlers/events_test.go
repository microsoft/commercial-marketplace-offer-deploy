package handlers

import (
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type testEventsMessageHandler struct {
	db               *gorm.DB
	deployment       *model.Deployment
	invokedOperation *model.InvokedOperation
	message          *sdk.EventHookMessage
}

func (h *testEventsMessageHandler) maxOutRetries() {
	h.invokedOperation.Retries = 0
	h.db.Save(h.invokedOperation)
}

func (h *testEventsMessageHandler) makeRetriable() {
	h.invokedOperation.Retries = 10
	h.invokedOperation.Attempts = 0
	h.db.Save(h.invokedOperation)
}

func newTestEventsMessageHandler() *testEventsMessageHandler {
	db := data.NewDatabase(&data.DatabaseOptions{
		UseInMemory: true,
	}).Instance()

	deployment := &model.Deployment{}
	db.Save(deployment)

	invokedOperation := &model.InvokedOperation{
		DeploymentId: deployment.ID,
	}
	db.Save(invokedOperation)

	message := &sdk.EventHookMessage{
		Id:      uuid.MustParse("22ed40f3-196a-43f8-9b5d-5459ce02ee45"),
		HookId:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		Type:    "deploymentCompleted",
		Status:  "failed",
		Subject: "/deployments/" + strconv.Itoa(int(deployment.ID)),
		Data: sdk.DeploymentEventData{
			Attempts:      1,
			DeploymentId:  1,
			OperationId:   invokedOperation.ID,
			CorrelationId: nil,
			StageId:       nil,
			Message:       "",
		},
	}

	return &testEventsMessageHandler{
		db:               db,
		deployment:       deployment,
		invokedOperation: invokedOperation,
		message:          message,
	}
}

func TestEventsHandlerUpdate(t *testing.T) {
	test := newTestEventsMessageHandler()

	handler := eventsMessageHandler{db: test.db}

	result, err := handler.update(test.message)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.EqualValues(t, test.message.Data.(sdk.DeploymentEventData).OperationId, result.ID)
}

func TestEventsHandlerShouldNotRetryIfRetriesExceeded(t *testing.T) {
	test := newTestEventsMessageHandler()
	test.maxOutRetries()

	handler := eventsMessageHandler{db: test.db}
	test.maxOutRetries()

	result, err := handler.update(test.message)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.EqualValues(t, test.message.Status, result.Status)
}

func TestEventsHandlerMessageStatusScheduledIfRetriable(t *testing.T) {
	test := newTestEventsMessageHandler()
	test.maxOutRetries()

	handler := eventsMessageHandler{db: test.db}

	test.makeRetriable()
	result, err := handler.update(test.message)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.EqualValues(t, "scheduled", result.Status)
}
