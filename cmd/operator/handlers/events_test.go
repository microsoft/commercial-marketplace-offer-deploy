package handlers

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type testEventsMessageHandler struct {
	db               *gorm.DB
	deployment       *data.Deployment
	invokedOperation *data.InvokedOperation
}

func newTestEventsMessageHandler() *testEventsMessageHandler {
	db := data.NewDatabase(&data.DatabaseOptions{
		UseInMemory: true,
	}).Instance()

	deployment := &data.Deployment{}
	db.Save(deployment)

	invokedOperation := &data.InvokedOperation{
		DeploymentId: deployment.ID,
	}
	db.Save(invokedOperation)

	return &testEventsMessageHandler{
		db:               db,
		deployment:       deployment,
		invokedOperation: invokedOperation,
	}
}

func TestEventsHandlerUpdate(t *testing.T) {
	test := newTestEventsMessageHandler()

	originalMessage := &events.EventHookMessage{
		Id:      uuid.MustParse("22ed40f3-196a-43f8-9b5d-5459ce02ee45"),
		HookId:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		Type:    "deploymentCompleted",
		Status:  "failed",
		Subject: "/deployments/" + strconv.Itoa(int(test.deployment.ID)),
		Data: events.DeploymentEventData{
			Attempts:      1,
			DeploymentId:  1,
			OperationId:   test.invokedOperation.ID,
			CorrelationId: nil,
			StageId:       nil,
			Message:       "",
		},
	}
	bytes, _ := json.Marshal(originalMessage)
	message := &events.EventHookMessage{}
	json.Unmarshal(bytes, message)

	handler := eventsMessageHandler{db: test.db}

	result, err := handler.update(message)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.EqualValues(t, originalMessage.Data.(events.DeploymentEventData).OperationId, result.ID)
}
