package messaging

import (
	"errors"
	"fmt"
)


type Publisher interface {
	Publish(message DeploymentMessage) error
}

type DeploymentMessage struct {
	Header DeploymentMessageHeader  `json:"header"`
	Body any `json:"body"`
}

type DeploymentMessageHeader struct {
	Topic string `json:"topic"`
}

type PublisherConfig struct {
	Type string
	TypeConfig any
}

func CreatePublisher(config PublisherConfig) (Publisher, error) {
	switch config.Type {
		case "servicebus": {
			serviceBusConfig := config.TypeConfig.(ServiceBusConfig)
			return NewServiceBusPublisher(serviceBusConfig.namespace, serviceBusConfig.queueName)
		}
		default: {
			errorMsg := fmt.Sprintf("There was no publisher for the provided type %s", config.Type)
			return nil, errors.New(errorMsg)
		}
	}
}

