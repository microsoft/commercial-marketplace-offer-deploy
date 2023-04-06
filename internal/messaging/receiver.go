package messaging

import (
	"context"
	"errors"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type MessageReceiver interface {
	Start()
	Stop()
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
	stopped   bool
	stop      chan bool
	ctx       context.Context
	queueName string
	namespace string
	handler   ServiceBusMessageHandler
}

func (r *serviceBusReceiver) Start() {
	r.stopped = false
	log.Println("starting the receiver")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Println("error getting default credential: ", err)
	}

	client, err := azservicebus.NewClient(r.namespace, cred, nil)

	if err != nil {
		log.Println("failure creating client")
	}

	receiver, err := client.NewReceiverForQueue(r.queueName, nil)
	if err != nil {
		log.Println("failure creating receiver")
	}

	defer receiver.Close(r.ctx)

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
					log.Printf("Error receiving messages: %s\n", err)
				}

				for _, message := range messages {
					log.Printf("Received message: %s\n", message.MessageID)

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
	r.ctx.Done()
	r.stop <- true
	<-r.stop
	close(r.stop)
}

func NewServiceBusReceiver(handler any, options ServiceBusMessageReceiverOptions) (MessageReceiver, error) {
	serviceBusMessageHandler, err := NewServiceMessageHandler(handler)
	if err != nil {
		return nil, err
	}

	receiver := serviceBusReceiver{
		stop:      make(chan bool),
		stopped:   true,
		queueName: options.QueueName,
		namespace: options.FullyQualifiedNamespace,
		ctx:       context.TODO(),
		handler:   serviceBusMessageHandler,
	}
	return &receiver, nil
}

//endregion servicebus receiver
