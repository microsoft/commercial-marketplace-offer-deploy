package dispatch

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_DispatchInvokedOperation_NilOperationId(t *testing.T) {
	operation := DispatchInvokedOperation{
		DeploymentId: 1,
		Name:         "test",
		Parameters:   nil,
		Retries:      0,
	}

	assert.Equal(t, uuid.Nil, operation.OperationId)
}
