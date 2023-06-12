package operation

import (
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
)

type OperationError struct {
	operation model.InvokedOperation
}

func (e *OperationError) Error() string {
	return ""
}

func NewError(operation *Operation) *OperationError {
	return &OperationError{
		operation: operation.InvokedOperation,
	}
}
