package events

import (
	"context"
	"log"
	"sync"
)

type WebHookPublisher interface {
	// publishes a message to all web hook subscriptions
	Publish(message *EventSubscriptionMessage) error
}

type webHookPublisher struct {
	subscriptionsProvider SubscriptionsProvider
	sender                MessageSender
}

func NewWebHookPublisher(sender MessageSender, subscriptionsProvider SubscriptionsProvider) WebHookPublisher {
	publisher := &webHookPublisher{sender: sender, subscriptionsProvider: subscriptionsProvider}

	return publisher
}

func (p *webHookPublisher) Publish(message *EventSubscriptionMessage) error {
	subscriptions, err := p.subscriptionsProvider.GetSubscriptions(message.EventType)

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
			err := p.sender.Send(ctx, &message)

			if err != nil {
				log.Printf("error sending message to subscription [%s]", subscription.Name)
			}
		}(i)
	}
	waitGroup.Wait()

	return nil
}
