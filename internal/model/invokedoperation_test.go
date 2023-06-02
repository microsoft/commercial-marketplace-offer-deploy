package model

import (
	"testing"

	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/stretchr/testify/assert"
)

//region running

func Test_InvokedOperation_Running(t *testing.T) {
	invokedOperation := InvokedOperation{
		Retries:  3,
		Attempts: 0,
	}

	invokedOperation.Running()
	assert.Equal(t, sdk.StatusRunning.String(), invokedOperation.Status)
}

func Test_InvokedOperation_MarkRunning_Doesnt_Change_If_Running(t *testing.T) {
	invokedOperation := InvokedOperation{
		Retries:  3,
		Attempts: 0,
		Status:   sdk.StatusScheduled.String(),
	}
	// calling running twice proves that the number of execution results remain the same
	invokedOperation.Running()
	invokedOperation.Running()

	assert.Equal(t, sdk.StatusRunning.String(), invokedOperation.Status)
	assert.Len(t, invokedOperation.Results, 1)
}

func Test_InvokedOperation_MarkRunning_Increments_Attempts(t *testing.T) {
	op := InvokedOperation{
		Retries:  3,
		Attempts: 0,
		Status:   sdk.StatusScheduled.String(),
	}

	op.Running()
	assert.Equal(t, uint(1), op.Attempts)
	assert.Len(t, op.Results, 1)
}

//endregion running

func Test_InvokedOperation_Complete_Prevents_Status_Change(t *testing.T) {
	op := InvokedOperation{
		Retries:  3,
		Attempts: 0,
		Status:   sdk.StatusScheduled.String(),
	}

	op.Running()
	op.Complete()

	// regardless of the status the operation was in, complete is permanent and will prevent status changes
	op.Failed()
	op.Schedule()

	assert.Equal(t, uint(1), op.Attempts)
	assert.Len(t, op.Results, 1)
	assert.Equal(t, true, op.IsCompleted())

	// will be the status that it was in right before complete was called
	assert.Equal(t, sdk.StatusRunning.String(), op.Status)
}

func Test_InvokedOperation_Complete_Prevents_Retrying(t *testing.T) {
	op := InvokedOperation{
		Retries:  3,
		Attempts: 0,
		Status:   sdk.StatusScheduled.String(),
	}

	op.Running()
	op.Success()
	op.Complete()

	// regardless the operation cannot be scheduled if already completed
	err := op.Schedule()

	assert.Error(t, err)
	assert.Equal(t, sdk.StatusSuccess.String(), op.Status)
	assert.Len(t, op.Results, 1)
}

func Test_InvokedOperation_Schedule_Runnning_Multiple_Times_Increments_Attemps(t *testing.T) {
	op := InvokedOperation{
		Retries:  3,
		Attempts: 0,
		Status:   sdk.StatusScheduled.String(),
	}

	op.Running()
	assert.Len(t, op.Results, 1)

	op.Failed()
	op.Schedule()
	op.Running()
	assert.Len(t, op.Results, 2)

	op.Complete()
	assert.Len(t, op.Results, 2)
}

func Test_InvokedOperation_Will_Not_Exceed_Retries(t *testing.T) {
	op := InvokedOperation{
		Retries:  3,
		Attempts: 0,
		Status:   sdk.StatusScheduled.String(),
	}

	for i := 0; i < 5; i++ {
		op.Running()
		op.Failed()
		op.Schedule()
	}
	assert.Len(t, op.Results, int(op.Retries))
}
