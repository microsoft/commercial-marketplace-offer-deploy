package messaging

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/utils"
)

const banner = `
🄶🄾🄻🄳🄴🄽 🅁🄴🄲🄴🄸🅅🄴🅁`

const DefaultMaxMessages = 10

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
	MaxMessages             int
}

//region servicebus receiver

type serviceBusReceiver struct {
	stopped     bool
	stop        chan bool
	ctx         context.Context
	queueName   string
	namespace   string
	handler     ServiceBusMessageHandler
	credential  azcore.TokenCredential
	logger      *log.Entry
	maxMessages int
}

func (r *serviceBusReceiver) Start() {
	fmt.Println(banner)
	r.logger.Trace("Starting")

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
				r.processMessages(receiver)
			}
		}
	}
}

func (r *serviceBusReceiver) processMessages(receiver *azservicebus.Receiver) {
	var messages []*azservicebus.ReceivedMessage = []*azservicebus.ReceivedMessage{}
	messages, err := receiver.ReceiveMessages(r.ctx, r.maxMessages, nil)
	if err != nil {
		log.Errorf("%s - error receiving: %v", r.queueName, err)
		return
	}

	log.Tracef("%s - received [%d] messages", r.queueName, len(messages))

	// waitGroup := sync.WaitGroup{}
	// waitGroup.Add(len(messages))

	for _, message := range messages {
		go func(message *azservicebus.ReceivedMessage) {
			//defer waitGroup.Done()
			err := r.executedHandler(receiver, message)
			if err != nil {
				log.Errorf("Error handling message: %v", err)
			}
		}(message)
	}
	//waitGroup.Wait()
}

func (r *serviceBusReceiver) executedHandler(receiver *azservicebus.Receiver, message *azservicebus.ReceivedMessage) error {
	if message == nil {
		return errors.New("message is nil")
	}

	log.WithFields(log.Fields{
		"queueName": r.queueName,
		"messageId": message.MessageID,
		"body":      string(message.Body),
	}).Trace("Received message")

	err := r.handler.Handle(r.ctx, message)

	if err != nil {
		log.Error(err)
	}

	err = receiver.CompleteMessage(r.ctx, message, nil)

	if err != nil {
		log.Errorf("Error completing message: %v", err)

		var sbErr *azservicebus.Error
		if errors.As(err, &sbErr) && sbErr.Code == azservicebus.CodeLockLost {
			// The message lock has expired. This isn't fatal for the client, but it does mean
			// that this message can be received by another Receiver (or potentially this one!).
			log.Warn("Message lock expired\n")
			// You can extend the message lock by calling receiver.RenewMessageLock(msg) before the
			// message lock has expired.
		}
		return err
	}
	log.Tracef("Completed message [%s]", message.MessageID)
	return nil
}

func (r *serviceBusReceiver) Stop() {
	log.Debug("Stopping message receiver")

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
		log.Error("failure creating client")
	}

	receiver, err := client.NewReceiverForQueue(r.queueName, nil)
	if err != nil {
		errorMessages = append(errorMessages, err.Error())
		log.Error("failure creating receiver")
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

	maxMessages := DefaultMaxMessages
	if options.MaxMessages > 0 {
		maxMessages = options.MaxMessages
	}

	receiver := serviceBusReceiver{
		stop:        make(chan bool),
		stopped:     true,
		queueName:   options.QueueName,
		namespace:   options.FullyQualifiedNamespace,
		ctx:         context.Background(),
		credential:  credential,
		handler:     serviceBusMessageHandler,
		logger:      log.WithField("queue", options.QueueName),
		maxMessages: maxMessages,
	}
	return &receiver, nil
}

//endregion servicebus receiver
