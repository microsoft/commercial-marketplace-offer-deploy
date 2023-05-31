package dispatch

import "github.com/google/uuid"

// We want to take a request to invoke and operation and dispatch.
// This is the command to do that.
type DispatchInvokedOperation struct {
	OperationId  uuid.UUID
	DeploymentId uint
	Name         string
	Parameters   map[string]any
	Attributes   map[string]any
	Retries      int
}
