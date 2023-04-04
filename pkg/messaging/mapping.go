package messaging

import (
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

// maps to a function that takes a message and returns a mapped type and an error
func MapTo[T any](message *azservicebus.ReceivedMessage) (*T, error) {
	mapped := new(T)
	err := json.Unmarshal(message.Body, &mapped)
	return mapped, err
}
