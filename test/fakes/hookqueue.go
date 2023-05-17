package fakes

import (
	"context"
	"testing"

	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
)

type FakeHookQueue struct {
	t        *testing.T
	messages []events.EventHookMessage
}

func NewFakeHookQueue(t *testing.T) *FakeHookQueue {
	return &FakeHookQueue{
		t:        t,
		messages: []events.EventHookMessage{},
	}
}

func (q *FakeHookQueue) Messages() []events.EventHookMessage {
	return q.messages
}

func (q *FakeHookQueue) Add(ctx context.Context, message *events.EventHookMessage) error {
	q.t.Logf("fakeHookQueue.Add called with message: %v", message)
	q.messages = append(q.messages, *message)
	return nil
}
