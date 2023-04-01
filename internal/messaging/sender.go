package messaging

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type MessageSender interface {
	Send(ctx context.Context, message ...any) error
}

// TODO: need to get the message bus from configuration

func NewMessageSender(credential azcore.TokenCredential) MessageSender {
	return nil
}
