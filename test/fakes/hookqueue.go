package fakes

import (
	"context"
	"testing"

	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

type FakeHookQueue struct {
	t        *testing.T
	messages []sdk.EventHookMessage
}

func NewFakeHookQueue(t *testing.T) *FakeHookQueue {
	return &FakeHookQueue{
		t:        t,
		messages: []sdk.EventHookMessage{},
	}
}

func (q *FakeHookQueue) Messages() []sdk.EventHookMessage {
	return q.messages
}

func (q *FakeHookQueue) Add(ctx context.Context, message *sdk.EventHookMessage) error {
	q.t.Logf("fakeHookQueue.Add called with message: %v", message)
	q.messages = append(q.messages, *message)
	return nil
}
