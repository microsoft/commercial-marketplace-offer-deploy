package messaging

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type ServiceBusPublisher func(message DeploymentMessage) error

func (f ServiceBusPublisher) Publish(message DeploymentMessage) error {
	return f(message)
}

// func NewServiceBusPublisher(connectionString string) (ServiceBusPublisher, error) {
// 	client, err := azservicebus.NewClientFromConnectionString(connectionString)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return func(message DeploymentMessage) error {
// 		sender, err := client.NewSender(azservicebus.TopicName(message.Header.Topic))
// 		if err != nil {
// 			return err
// 		}

// 		return sender.SendMessage(context.Background(), azservicebus.Message{
// 			Body: message.Body,
// 		})
// 	}, nil
// }
