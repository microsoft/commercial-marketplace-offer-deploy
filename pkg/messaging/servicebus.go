package messaging

import (
	"context"
	"encoding/json"
	//"os"

	//"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type ServiceBusConfig struct {
	Namespace string
	QueueName string
}

type ServiceBusPublisher func(message DeploymentMessage) error

func (f ServiceBusPublisher) Publish(message DeploymentMessage) error {
	return f(message)
}

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
