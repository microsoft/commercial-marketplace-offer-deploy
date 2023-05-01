package messaging

// this represents the message (command) that is sent to the operator to be executed
type ExecuteInvokedOperation struct {
	OperationId string `json:"operationId"`
}
