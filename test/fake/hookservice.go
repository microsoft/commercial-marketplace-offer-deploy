package fake

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

type FakeHookService struct {
	t        *testing.T
	messages []sdk.EventHookMessage
}

func NewFakeHookService(t *testing.T) *FakeHookService {
	return &FakeHookService{
		t:        t,
		messages: []sdk.EventHookMessage{},
	}
}

func (q *FakeHookService) Messages() []sdk.EventHookMessage {
	return q.messages
}

func (q *FakeHookService) Notify(ctx context.Context, message *sdk.EventHookMessage) (uuid.UUID, error) {
	message.Id = uuid.New()
	q.t.Logf("Fake HookService called with message: %v", message)

	bytes, _ := json.Marshal(message)
	unmarshaled := sdk.EventHookMessage{}
	json.Unmarshal(bytes, &unmarshaled)

	q.messages = append(q.messages, unmarshaled)
	return message.Id, nil
}
