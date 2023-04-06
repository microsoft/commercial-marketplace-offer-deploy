package messaging

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/stretchr/testify/assert"
)

func TestGetHandleMethod(t *testing.T) {
	handler := &testHandler{
		t:          t,
		expectedId: "testMethod",
	}
	method := getHandleMethod(handler)

	assert.Equal(t, "func", method.Kind().String())

	// should return with no error
	result := method.Call([]reflect.Value{
		reflect.ValueOf(&testMessage{Id: "testMethod"}),
		reflect.ValueOf(MessageHandlerContext{ReceivedMessage: nil}),
	})[0]

	assert.Nil(t, result.Interface())
}

func TestValidationtHandler(t *testing.T) {

	// pointer to struct should be OK
	err := validateHandler(&testHandler{t: t})
	assert.NoError(t, err)

	// non pointer to struct should be OK
	err = validateHandler(testHandler{t: t})
	assert.Error(t, err)

	//bad
	err = validateHandler(testBadHandler{t: t})
	assert.Error(t, err)
}

func TestShouldExecuteHandler(t *testing.T) {
	handler := &testHandler{
		t:          t,
		expectedId: "e2eTestId",
	}
	serviceBusHandler, err := NewServiceMessageHandler(handler)
	assert.NoError(t, err)

	bytes, _ := json.Marshal(&testMessage{Id: "e2eTestId"})
	message := &azservicebus.ReceivedMessage{
		Body: bytes,
	}

	err = serviceBusHandler.Handle(context.TODO(), message)
	assert.NoError(t, err)
}

//region setup

type testHandler struct {
	t          *testing.T
	expectedId string
}

type testMessage struct {
	Id    string
	Error error
}

func (h *testHandler) Handle(message *testMessage, context MessageHandlerContext) error {
	t := h.t
	assert.Equal(t, h.expectedId, message.Id)
	return message.Error
}

// inccorect handler
type testBadHandler struct {
	t *testing.T
}

func (h *testBadHandler) WrongName(message *testMessage, context MessageHandlerContext) error {
	return nil
}

//endregion setup
