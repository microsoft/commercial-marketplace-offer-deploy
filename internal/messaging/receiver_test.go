package messaging

import (
	"context"
	"testing"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/stretchr/testify/assert"
)

func TestReceiveError(t *testing.T) {

	namespace := "bobjacgps"
	queueName := "testqueue"
	ctx := context.TODO()

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		assert.NoError(t, err)
	}
	errorMessages := []string{}

	client, err := azservicebus.NewClient(namespace, cred, nil)
	if err != nil {
		errorMessages = append(errorMessages, err.Error())
		assert.NoError(t, err)
	}

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	if err != nil {
		errorMessages = append(errorMessages, err.Error())
		assert.NoError(t, err)
	}

	var messages []*azservicebus.ReceivedMessage = []*azservicebus.ReceivedMessage{}
	messages, err = receiver.ReceiveMessages(ctx, 1, nil)
	assert.True(t, len(messages) == 0)
}