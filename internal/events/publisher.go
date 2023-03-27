package events

import (
	"context"
	"sync"
)

type Publisher interface {
	Publish(eventType EventType, message *EventSubscriptionMessage) error
}

type webHookPublisher struct {
	subscriptionsProvider SubscriptionsProvider
	sender                MessageSender
}

func NewWebHookPublisher(sender MessageSender, subscriptionsProvider SubscriptionsProvider) Publisher {
	publisher := &webHookPublisher{sender: sender, subscriptionsProvider: subscriptionsProvider}

	return publisher
}

func (p *webHookPublisher) Publish(eventType EventType, message *EventSubscriptionMessage) error {
	subscriptions, err := p.subscriptionsProvider.GetSubscriptions(eventType)

	if err != nil {
		return err
	}

	subscriptionsCount := len(subscriptions)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(subscriptionsCount)

	var ctx context.Context = context.Background()

	for i := 0; i < subscriptionsCount; i++ {
		go func(i int) {
			defer waitGroup.Done()
			subscription := subscriptions[i]
			p.sender.Send(ctx, &subscription)
		}(i)
	}
	waitGroup.Wait()

	return nil
}
