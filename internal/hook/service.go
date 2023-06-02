package hook

import (
	"context"
	"errors"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/config"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/messaging"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	log "github.com/sirupsen/logrus"
)

var (
	instance     Service
	instanceOnce sync.Once
	instanceErr  error
)

// notify is the function signature for the event hook Add
type NotifyFunc func(ctx context.Context, message *sdk.EventHookMessage) error

const eventsQueueName = string(messaging.QueueNameEvents)

// This implementation is to make the semantics clear that this is the lifecycle of a hook message:
// eventHookMessage --> added to queue --> received --> executed handler (events) --> publish the message using Publisher.Publish()
// queue for adding hook messages to be published
type Service interface {
	// adds a message to the hooks queue
	Notify(ctx context.Context, message *sdk.EventHookMessage) error
}

type service struct {
	queueName string
	sender    messaging.MessageSender
}

// notification
func (q *service) Notify(ctx context.Context, message *sdk.EventHookMessage) error {
	if message == nil {
		return errors.New("message is nil")
	}

	results, err := q.sender.Send(ctx, q.queueName, message)
	if err != nil {
		log.Errorf("Error attempting toadd event message to queue [%s]: %v", q.queueName, err)
		return err
	} else {
		log.Tracef("EventHook message sent [%s]", message.Id)
	}
	if len(results) > 0 {
		for _, result := range results {
			if result.Error != nil {
				log.Errorf("Error sending event message: %v", result.Error)
				return result.Error
			}
		}
	}
	return nil
}

// enqueues a message to the event hooks service
func Notify(ctx context.Context, message *sdk.EventHookMessage) error {
	if instance == nil {
		return errors.New("hook queue not configured. call Configure() first")
	}

	if message != nil {
		if message.Id == uuid.Nil {
			message.Id = uuid.New()
		}
	}

	return instance.Notify(ctx, message)
}

func Configure(appConfig *config.AppConfig) error {
	instanceOnce.Do(func() {
		credential, err := getAzureCredential()
		if err != nil {
			log.Errorf("Error creating Azure credential for hook.Queue: %v", err)
		}

		sender, err := messaging.NewServiceBusMessageSender(credential, messaging.MessageSenderOptions{
			SubscriptionId:          appConfig.Azure.SubscriptionId,
			Location:                appConfig.Azure.Location,
			ResourceGroupName:       appConfig.Azure.ResourceGroupName,
			FullyQualifiedNamespace: appConfig.Azure.GetFullQualifiedNamespace(),
		})
		if err != nil {
			log.Errorf("Error creating message sender for hook.Queue: %v", err)
			instanceErr = err
			return
		}

		instance = NewService(sender)
	})
	return instanceErr
}

func SetInstance(i Service) {
	instance = i
}

func NewService(messageSender messaging.MessageSender) Service {
	return &service{
		queueName: eventsQueueName,
		sender:    messageSender,
	}
}

func getAzureCredential() (azcore.TokenCredential, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	return credential, err
}
