package messaging

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type MessageHandler interface {
	Handle(ctx context.Context, message *azservicebus.ReceivedMessage) error
}

type serviceBusMessageHandler struct {
}

// Handle implements MessageHandler
func (*serviceBusMessageHandler) Handle(ctx context.Context, message *azservicebus.ReceivedMessage) error {
	//msg, err := MapTo[DeploymentMessage](message)
	panic("Not implemented")
}

func NewServiceBusHandler() MessageHandler {
	return &serviceBusMessageHandler{}
}
