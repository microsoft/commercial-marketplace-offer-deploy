package fake

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

type FakeNotifyFunc struct {
	t        *testing.T
	messages []*sdk.EventHookMessage
	count    int
}

func (fake *FakeNotifyFunc) Notify(ctx context.Context, message *sdk.EventHookMessage) (uuid.UUID, error) {
	fake.count++

	fake.t.Logf("Notify called [%d]", fake.count)

	id := uuid.New()
	message.Id = id
	fake.messages = append(fake.messages, message)

	return id, nil
}

// all messages passed to notify
func (fake *FakeNotifyFunc) Messages() []*sdk.EventHookMessage {
	return fake.messages
}

// message count
func (fake *FakeNotifyFunc) Count() int {
	return fake.count
}

func NewFakeNotifyFunc(t *testing.T) *FakeNotifyFunc {
	return &FakeNotifyFunc{t: t}
}
