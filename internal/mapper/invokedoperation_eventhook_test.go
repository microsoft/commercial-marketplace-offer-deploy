package mapper

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
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
		Status:   "Started",
		Attempts: uint(1),
	}

	sdkResult := &sdk.DryRunResult{
		Status: "Success",
	}

	resultsMap := make(map[uint]*model.InvokedOperationResult)
	resultsMap[uint(1)] = &model.InvokedOperationResult{
		Attempt: uint(1),
		Status:  "Completed",
		Value:   sdkResult,
	}
	invokedOperation.Results = resultsMap

	// Act
	result := getDryRunData(invokedOperation)

	// Assert
	assert.EqualValues(t, "Success", result.(*sdk.DryRunEventData).Status)
}

func Test_getBaseEventData_sets_ids_and_attempts(t *testing.T) {
	// Arrange
	invokedOperation := &model.InvokedOperation{
		BaseWithGuidPrimaryKey: model.BaseWithGuidPrimaryKey{
			ID: uuid.New(),
		},
		Name:         sdk.OperationDeploy.String(),
		DeploymentId: uint(1),
		Status:       "Started",
		Attempts:     uint(1),
	}

	result := getBaseEventData(invokedOperation)

	assert.Equal(t, invokedOperation.ID, result.OperationId)
	assert.Equal(t, invokedOperation.DeploymentId, uint(result.DeploymentId))
	assert.Equal(t, invokedOperation.Attempts, uint(result.Attempts))
}

func Test_getBaseEventData_Subject_has_stageId_when_stage_operation_and_stageId_param_exists(t *testing.T) {
	stageId := uuid.New().String()

	operation := model.InvokedOperation{
		BaseWithGuidPrimaryKey: model.BaseWithGuidPrimaryKey{
			ID: uuid.New(),
		},
		Name:         sdk.OperationDeploy.String(),
		DeploymentId: uint(1),
		Attempts:     uint(1),
		Parameters: map[string]interface{}{
			string(model.ParameterKeyStageId): stageId,
		},
	}

	result := MapInvokedOperation(&operation)

	assert.Contains(t, result.Subject, "deployments/1")
	assert.NotContains(t, result.Subject, fmt.Sprintf("stages/%s", stageId))

	operation2 := operation
	operation2.Name = sdk.OperationRetryStage.String()
	result = MapInvokedOperation(&operation2)

	assert.Contains(t, result.Subject, fmt.Sprintf("stages/%s", stageId))
}
