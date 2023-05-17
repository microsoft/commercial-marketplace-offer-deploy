package fakes

import (
	"context"
	"testing"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
)

type FakeMessageSender struct {
	t *testing.T
}

func (s *FakeMessageSender) Send(ctx context.Context, queueName string, messages ...any) ([]messaging.SendMessageResult, error) {
	s.t.Log("Send called")
	return nil, nil
}

func (s *FakeMessageSender) EnsureTopology(ctx context.Context, queueName string) error {
	s.t.Log("EnsureTopology called")
	return nil
}

func NewFakeMessageSender(t *testing.T) *FakeMessageSender {
	return &FakeMessageSender{}
}
