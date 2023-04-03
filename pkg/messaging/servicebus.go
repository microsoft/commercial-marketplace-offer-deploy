package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
//	"github.com/sqs/goreturns/returns"
)

type ServiceBusConfig struct {
	Namespace string
	QueueName string
}

type ServiceBusReceiver struct {
	stopped bool
	stop chan bool
	ctx context.Context
	queueName string 
	namespace string
	handler func(message *azservicebus.ReceivedMessage) (any, error)
}

func MapDeploymentMessage(message *azservicebus.ReceivedMessage) (any, error) {
	var deploymentMessage DeploymentMessage
	err := json.Unmarshal(message.Body, &deploymentMessage)
	if err != nil {
		log.Println("Failure")
	}
	return deploymentMessage, err
}

func (r *ServiceBusReceiver) Start()  {
	r.stopped = false
	log.Println("Starting the receiver")

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Println("Failure")
	}

	client, err := azservicebus.NewClient(r.namespace, cred, nil)

	if err != nil {
		log.Println("Failure")
	}
	log.Println("Created the client")
	log.Println("Starting the receiver loop")
	receiver, err := client.NewReceiverForQueue(r.queueName, nil)
	if err != nil {
		log.Println("Failure")
	}
	log.Println("Created the receiver from client")
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
				log.Println("inside of default")
				var messages []*azservicebus.ReceivedMessage = []*azservicebus.ReceivedMessage{}
				messages, err = receiver.ReceiveMessages(r.ctx, 1, nil)
				if err != nil {
					log.Printf("Error receiving messages: %s\n", err)
				}
				log.Println("received messages completed in anonymous function")
				log.Println("after receive messages")
				log.Printf("%d messages received\n", len(messages))
				for _, message := range messages {
					deploymentMessage, err := r.handler(message)
					//var deploymentMessage DeploymentMessage
					//err := json.Unmarshal(message.Body, &deploymentMessage)
					if err != nil {
						log.Println("Failure")
					}
					log.Printf("Received message: %s\n", deploymentMessage)
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
					log.Printf("Completed message: %s\n", deploymentMessage)
				}
			}
		}
	}	
}

func (r *ServiceBusReceiver) Stop() {
	log.Println("Stopping the receiver")
	//r.stopped = true
	r.ctx.Done()
	log.Println("Done with the context")
	r.stop <- true
	<-r.stop
	close(r.stop)
}


func NewServiceBusReceiver(config ServiceBusConfig) (*ServiceBusReceiver, error) {
	receiver := ServiceBusReceiver{
		stop: make(chan bool),
		stopped: true,
		queueName: config.QueueName,
		namespace: config.Namespace,
		ctx: context.TODO(),
		handler: MapDeploymentMessage,
	}
	return &receiver, nil
}

type ServiceBusPublisher func(message DeploymentMessage) error

func (f ServiceBusPublisher) Publish(message DeploymentMessage) error {
	return f(message)
}

func NewServiceBusPublisher(ns string, queueName string) (ServiceBusPublisher, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	client, err := azservicebus.NewClient(ns, cred, nil)

	if err != nil {
		return nil, err
	}

	return func(message DeploymentMessage) error {
		sender, err := client.NewSender(queueName, nil)
		if err != nil {
			return err
		}
		defer sender.Close(context.TODO())

		jsonContent, err := json.Marshal(message)
		if err != nil {
			return err
		}

		sbMessage := &azservicebus.Message {
			Body: []byte(jsonContent),
		}

		return sender.SendMessage(context.TODO(), sbMessage, nil)
	}, nil
}

