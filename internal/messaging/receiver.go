package messaging

import (
	"context"
	"fmt"
	"log"

	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
)

const banner = `
🄶🄾🄻🄳🄴🄽 🅁🄴🄲🄴🄸🅅🄴🅁`

type MessageReceiver interface {
	Start()
	Stop()
	GetName() string
}

type MessageReceiverOptions struct {
	QueueName string
}

type ServiceBusMessageReceiverOptions struct {
	MessageReceiverOptions
	FullyQualifiedNamespace string
}

//region servicebus receiver

type serviceBusReceiver struct {
	stopped    bool
	stop       chan bool
	ctx        context.Context
	queueName  string
	namespace  string
	handler    ServiceBusMessageHandler
	credential azcore.TokenCredential
}

func (r *serviceBusReceiver) Start() {
	fmt.Println(banner)
	log.Printf("starting message receiver: %s", r.queueName)

	r.stopped = false
	receiver, err := r.getQueueReceiver()
	if receiver != nil {
		defer receiver.Close(r.ctx)
	}

	// if the receiver was created, then don't bother continuing
	if err != nil {
		r.Stop()
		return
	}

	for {
		select {
		case <-r.stop:
			return
		default:
			for {
				if r.stopped {
					break
				}

				var messages []*azservicebus.ReceivedMessage = []*azservicebus.ReceivedMessage{}
				messages, err = receiver.ReceiveMessages(r.ctx, 1, nil)
				if err != nil {
					log.Printf("%s - error receiving: %s\n", r.queueName, err)
				}

				for _, message := range messages {
					log.Printf("%s - Received message: %s\n", r.queueName, message.MessageID)

					err := r.handler.Handle(r.ctx, message)

					if err != nil {
						log.Println(err)
					}
					err = receiver.CompleteMessage(context.TODO(), message, nil)
					if err != nil {
						var sbErr *azservicebus.Error
						if errors.As(err, &sbErr) && sbErr.Code == azservicebus.CodeLockLost {
							// The message lock has expired. This isn't fatal for the client, but it does mean
							// that this message can be received by another Receiver (or potentially this one!).
							log.Printf("Message lock expired\n")
							// You can extend the message lock by calling receiver.RenewMessageLock(msg) before the
							// message lock has expired.
							continue
						}
					}
					log.Printf("Completed message: %s\n", message.MessageID)
				}
			}
		}
	}
}

func (r *serviceBusReceiver) Stop() {
	log.Println("Stopping message receiver")

	r.ctx.Done()
	r.stop <- true
	<-r.stop
	close(r.stop)
}

func (r *serviceBusReceiver) getQueueReceiver() (*azservicebus.Receiver, error) {
	errorMessages := []string{}

	client, err := azservicebus.NewClient(r.namespace, r.credential, nil)
	if err != nil {
		errorMessages = append(errorMessages, err.Error())
		log.Println("failure creating client")
	}

	receiver, err := client.NewReceiverForQueue(r.queueName, nil)
	if err != nil {
		errorMessages = append(errorMessages, err.Error())
		log.Println("failure creating receiver")
	}

	if len(errorMessages) > 0 {
		return nil, utils.NewAggregateError(errorMessages)
	}

	return receiver, nil
}

func (r *serviceBusReceiver) GetName() string {
	return r.queueName
}

func NewServiceBusReceiver(handler any, credential azcore.TokenCredential, options ServiceBusMessageReceiverOptions) (MessageReceiver, error) {
	serviceBusMessageHandler, err := NewServiceMessageHandler(handler)
	if err != nil {
		return nil, err
	}

	receiver := serviceBusReceiver{
		stop:       make(chan bool),
		stopped:    true,
		queueName:  options.QueueName,
		namespace:  options.FullyQualifiedNamespace,
		ctx:        context.TODO(),
		credential: credential,
		handler:    serviceBusMessageHandler,
	}
	return &receiver, nil
}

//endregion servicebus receiver
