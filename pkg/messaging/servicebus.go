package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	//"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type ServiceBusConfig struct {
	Namespace string
	QueueName string
}

type ServiceBusReceiver struct {
	stop chan bool
	queueName string 
	namespace string
}

func (r *ServiceBusReceiver) Start() error {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}

	client, err := azservicebus.NewClient(r.namespace, cred, nil)

	if err != nil {
		return err
	}

	go func() {
		receiver, err := client.NewReceiverForQueue(r.queueName, nil)
		if err != nil {
			return
		}
		defer receiver.Close(context.TODO())
		
		ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
		
		for {
			select {
				case <-r.stop: {
					fmt.Println("Stopping receiver")
					log.Printf("Logging the stop of receiver")
					cancel()
					break
				}
				default: {
					messages, err := receiver.ReceiveMessages(ctx, 1, nil)
					log.Printf("%d messages received\n", len(messages))
					if err != nil {
						log.Println("err was true.  Calling cancel()")
					//	cancel()
						log.Println("returning")
						break 
					}
					for _, message := range messages {
						var deploymentMessage DeploymentMessage
						err := json.Unmarshal(message.Body, &deploymentMessage)
						if err != nil {
							cancel()
							return
						}
						log.Printf("Received message: %s\n", deploymentMessage)
						err = receiver.CompleteMessage(context.TODO(), message, nil)
						if err != nil {
							cancel()
							return
						}
						log.Printf("Completed message: %s\n", deploymentMessage)
					}
				}
			}
			log.Println("inside of for loop")
		}	
		//fmt.Println("outside of for loop")
	}()

	return nil
}

func (r *ServiceBusReceiver) Stop() {
	r.stop <- true
}

func NewServiceBusReceiver(config ServiceBusConfig) (*ServiceBusReceiver, error) {
	receiver := ServiceBusReceiver{
		stop: make(chan bool),
		queueName: config.QueueName,
		namespace: config.Namespace,
	}
	return &receiver, nil
}

type ServiceBusPublisher func(message DeploymentMessage) error
//type ServiceBusBackgroundReceiver func(chan bool) error

func (f ServiceBusPublisher) Publish(message DeploymentMessage) error {
	return f(message)
}

// func (f ServiceBusBackgroundReceiver) Start(stop chan bool) error {
// 	return f(stop)
// }

func NewServiceBusPublisher(ns string, queueName string) (ServiceBusPublisher, error) {
	// ns := os.Getenv("SERVICEBUS_ENDPOINT")
	// var credsToAdd []azcore.TokenCredential
	// cliCred, err := azidentity.NewAzureCLICredential(nil)
	// if err != nil {
	// 	return nil, err
	// }
	// envCred, err := azidentity.NewEnvironmentCredential(nil)
	// if err != nil {
	// 	return nil, err
	// }

	// //todo: adjust client credentials in accordance to api
	// credsToAdd = append(credsToAdd, cliCred, envCred)
	// cred, err := azidentity.NewChainedTokenCredential(credsToAdd, nil)
	// if err != nil {
	// 	return nil, err
	// }

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

// func NewServiceBusBackgroundReceiver(ns string, queueName string, stop chan bool) (ServiceBusBackgroundReceiver, error) {
// 	cred, err := azidentity.NewDefaultAzureCredential(nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	client, err := azservicebus.NewClient(ns, cred, nil)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return func(stop chan bool) error {
// 		go func() {
// 			receiver, err := client.NewReceiverForQueue(queueName, nil)
// 			if err != nil {
// 				return
// 			}
// 			defer receiver.Close(context.TODO())
			
// 			ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
			
// 			for {
// 				select {
// 					case <-stop: {
// 						fmt.Println("Stopping receiver")
// 						log.Printf("Logging the stop of receiver")
// 						cancel()
// 						break
// 					}
// 					default: {
// 						messages, err := receiver.ReceiveMessages(ctx, 1, nil)
// 						log.Printf("%d messages received\n", len(messages))
// 						if err != nil {
// 							log.Println("err was true.  Calling cancel()")
// 							cancel()
// 							log.Println("returning")
// 							return 
// 						}
// 						for _, message := range messages {
// 							var deploymentMessage DeploymentMessage
// 							err := json.Unmarshal(message.Body, &deploymentMessage)
// 							if err != nil {
// 								cancel()
// 								return
// 							}
// 							log.Printf("Received message: %s\n", deploymentMessage)
// 							err = receiver.CompleteMessage(context.TODO(), message, nil)
// 							if err != nil {
// 								cancel()
// 								return
// 							}
// 							log.Printf("Completed message: %s\n", deploymentMessage)
// 						}
// 					}
// 				}
// 				log.Println("inside of for loop")
// 			}	
// 			//fmt.Println("outside of for loop")
// 		}()
// 		return nil
// 	}, nil
// }
