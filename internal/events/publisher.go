package events

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	model "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
)

type EventHookPublisher interface {
	// publishes a message to all web hook subscriptions
	Publish(message *model.EventHookMessage) error
}

type eventHookPublisher struct {
	provider EventHooksProvider
	senders  map[uuid.UUID]WebHookSender
}

func NewWebHookPublisher(subscriptionsProvider EventHooksProvider) EventHookPublisher {
	publisher := &eventHookPublisher{senders: map[uuid.UUID]WebHookSender{}, provider: subscriptionsProvider}

	return publisher
}

func (p *eventHookPublisher) Publish(message *model.EventHookMessage) error {
	hooks, err := p.provider.Get()

	if err != nil {
		return err
	}

	hookCount := len(hooks)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(hookCount)

	var ctx context.Context = context.Background()

	for i := 0; i < hookCount; i++ {
		go func(i int) {
			defer waitGroup.Done()
			subscription := hooks[i]
			message.HookId = subscription.ID
			sender := p.getSender(*subscription)
			err := sender.Send(ctx, &message)

			if err != nil {
				log.Printf("error sending message to subscription [%s]", subscription.Name)
			}
		}(i)
	}
	waitGroup.Wait()

	return nil
}

func (p *eventHookPublisher) getSender(subscription data.EventHook) WebHookSender {
	if _, ok := p.senders[subscription.ID]; !ok {
		p.senders[subscription.ID] = NewMessageSender(subscription.Callback, subscription.ApiKey)
	}
	return p.senders[subscription.ID]
}
