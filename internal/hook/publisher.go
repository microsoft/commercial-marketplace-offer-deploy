package hook

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/data"
	model "github.com/microsoft/commercial-marketplace-offer-deploy/pkg/events"
)

// Publishes event hook messages to all web hooks registered in the system.
type Publisher interface {
	// publishes a message to all web hook subscriptions
	Publish(message *model.EventHookMessage) error
}

type publisher struct {
	provider EventHooksProvider
	senders  map[uuid.UUID]hookSender
}

func NewEventHookPublisher(subscriptionsProvider EventHooksProvider) Publisher {
	publisher := &publisher{senders: map[uuid.UUID]hookSender{}, provider: subscriptionsProvider}

	return publisher
}

func (p *publisher) Publish(message *model.EventHookMessage) error {
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
			hook := hooks[i]
			message.Id = uuid.New().String()
			message.HookId = hook.ID.String()

			sender := p.getSender(*hook)

			log.Debugf("sending message [%s] to %s - '%s'", message.Id, message.HookId, hook.Callback)
			err := sender.Send(ctx, &message)

			if err != nil {
				log.Error("error posting to callback '%s' [%s]", hook.Callback, hook.Name)
			}
		}(i)
	}
	waitGroup.Wait()

	return nil
}

func (p *publisher) getSender(subscription data.EventHook) hookSender {
	if _, ok := p.senders[subscription.ID]; !ok {
		p.senders[subscription.ID] = newHookSender(subscription.Callback, subscription.ApiKey)
	}
	return p.senders[subscription.ID]
}
