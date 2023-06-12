package operation

import "github.com/google/uuid"

// this represents the message (command) that is sent to the operator to be executed
//
//	remarks: this will be sent over the messaging infrastructure
type ExecuteOperationCommand struct {
	OperationId uuid.UUID `json:"operationId"`
}
