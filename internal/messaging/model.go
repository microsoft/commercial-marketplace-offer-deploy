package messaging

import "github.com/google/uuid"

// this represents the message (command) that is sent to the operator to be executed
type ExecuteInvokedOperation struct {
	OperationId uuid.UUID `json:"operationId"`
}
