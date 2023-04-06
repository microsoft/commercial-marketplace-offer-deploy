package events

import (
	"context"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	model "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
)

type WebHookPublisher interface {
	// publishes a message to all web hook subscriptions
	Publish(message *model.WebHookEventMessage) error
}

type webHookPublisher struct {
	subscriptionsProvider SubscriptionsProvider
	senders               map[uuid.UUID]WebHookSender
}

func NewWebHookPublisher(subscriptionsProvider SubscriptionsProvider) WebHookPublisher {
	publisher := &webHookPublisher{senders: map[uuid.UUID]WebHookSender{}, subscriptionsProvider: subscriptionsProvider}

	return publisher
}

func (p *webHookPublisher) Publish(message *model.WebHookEventMessage) error {
	subscriptions, err := p.subscriptionsProvider.GetSubscriptions()

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
			message.SubscriptionId = subscription.ID
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

func (p *webHookPublisher) getSender(subscription data.EventSubscription) WebHookSender {
	if _, ok := p.senders[subscription.ID]; !ok {
		p.senders[subscription.ID] = NewMessageSender(subscription.Callback, subscription.ApiKey)
	}
	return p.senders[subscription.ID]
}
