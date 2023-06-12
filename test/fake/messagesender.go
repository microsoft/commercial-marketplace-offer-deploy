package fake

import (
	"context"
	"testing"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
)

type FakeMessageSender struct {
	t         *testing.T
	messages  []any
	queueName string
	sendCalls uint
}

// all messages passed to notify
func (fake *FakeMessageSender) Messages() []any {
	return fake.messages
}

// message count
func (fake *FakeMessageSender) Count() uint {
	return fake.sendCalls
}

func (fake *FakeMessageSender) Last() any {
	return fake.messages[len(fake.messages)-1]
}

func (fake *FakeMessageSender) Send(ctx context.Context, queueName string, messages ...any) ([]messaging.SendMessageResult, error) {
	fake.sendCalls++
	fake.t.Logf("Send call [%d] on fake message sender", fake.sendCalls)

	fake.queueName = queueName
	fake.messages = append(fake.messages, messages...)

	results := make([]messaging.SendMessageResult, len(messages))

	return results, nil
}

func (s *FakeMessageSender) EnsureTopology(ctx context.Context, queueName string) error {
	s.t.Log("EnsureTopology called")
	return nil
}

func NewFakeMessageSender(t *testing.T) *FakeMessageSender {
	return &FakeMessageSender{
		t:        t,
		messages: []any{},
	}
}
