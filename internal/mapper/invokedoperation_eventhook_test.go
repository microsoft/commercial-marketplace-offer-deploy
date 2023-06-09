package mapper

import (
	"testing"

	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
	"github.com/stretchr/testify/assert"
)

func TestGetDryRunData_should_return_status_of_invokedoperation(t *testing.T) {
	// Arrange
	invokedOperation := &model.InvokedOperation{
		Status: "Started",
	}

	// Act
	result := getDryRunData(invokedOperation)

	// Assert
	assert.EqualValues(t, "Started", result.(*sdk.DryRunEventData).Status)
}

func TestGetDryRunData_should_return_status_of_invokedoperation_Results(t *testing.T) {
	// Arrange
	invokedOperation := &model.InvokedOperation{
		Status: "Started",
		Attempts: uint(1),
	}

	sdkResult := &sdk.DryRunResult{
		Status: "Success",
	}

	resultsMap := make(map[uint]*model.InvokedOperationResult)
	resultsMap[uint(1)] = &model.InvokedOperationResult{
		Attempt: uint(1),
		Status: "Completed",
		Value: sdkResult,
	}
	invokedOperation.Results = resultsMap

	// Act
	result := getDryRunData(invokedOperation)

	// Assert
	assert.EqualValues(t, "Success", result.(*sdk.DryRunEventData).Status)
}
