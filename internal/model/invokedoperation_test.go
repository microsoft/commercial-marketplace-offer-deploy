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

	changed := invokedOperation.Running()
	assert.Equal(t, true, changed)
	assert.Equal(t, sdk.StatusRunning.String(), invokedOperation.Status)
}

func Test_InvokedOperation_MarkRunning_Doesnt_Change_If_Running(t *testing.T) {
	invokedOperation := InvokedOperation{
		Retries:  3,
		Attempts: 0,
		Status:   sdk.StatusRunning.String(),
	}

	changed := invokedOperation.Running()
	assert.Equal(t, false, changed)
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
