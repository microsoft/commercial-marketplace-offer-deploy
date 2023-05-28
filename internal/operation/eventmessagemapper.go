package operation

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

// performs mapping of the invoked operation to the correct event hook message

func mapToMessage(invokedOperation *Operation) *sdk.EventHookMessage {
	message := &sdk.EventHookMessage{
		Status: invokedOperation.Status,
		Type:   getEventType(invokedOperation),
		Error:  invokedOperation.LatestResult().Error,
		Data:   getEventData(invokedOperation),
	}
	message.SetSubject(uint(invokedOperation.DeploymentId), nil)

	return message
}

func getEventType(invokedOperation *Operation) string {
	eventType := ""

	if invokedOperation.Name == sdk.OperationDeploy.String() {
		eventType = string(sdk.EventTypeDeploymentCompleted)
		if invokedOperation.Attempts > 1 {
			eventType = string(sdk.EventTypeDeploymentRetried)
		}
	} else if invokedOperation.Name == sdk.OperationRetry.String() {
		eventType = string(sdk.EventTypeDeploymentRetried)
	} else if invokedOperation.Name == sdk.OperationDryRun.String() {
		eventType = string(sdk.EventTypeDryRunCompleted)
	} else if invokedOperation.Name == sdk.OperationRetryStage.String() {
		return string(sdk.EventTypeStageRetried)
	}

	return eventType
}

func getEventData(invokedOperation *Operation) any {
	if invokedOperation.Name == sdk.OperationDeploy.String() {
		return getDeploymentData(invokedOperation)
	}
	if invokedOperation.Name == sdk.OperationDryRun.String() {
		return getDryRunData(invokedOperation)
	}
	if invokedOperation.Name == sdk.OperationRetry.String() {
		return getRetryData(invokedOperation)
	}
	if invokedOperation.Name == sdk.OperationRetryStage.String() {
		return getRetryStageData(invokedOperation)
	}
	return nil
}

func getRetryData(operation *Operation) any {
	return &sdk.DeploymentEventData{
		DeploymentId: int(operation.DeploymentId),
		OperationId:  operation.ID,
		Message:      fmt.Sprintf("Retry deployment %s", operation.Status),
	}
}

func getRetryStageData(operation *Operation) any {
	return &sdk.DeploymentEventData{
		DeploymentId: int(operation.DeploymentId),
		StageId:      to.Ptr(uuid.MustParse(operation.Parameters["stageId"].(string))),
		OperationId:  operation.ID,
		Message:      fmt.Sprintf("Retry deployment %s", operation.Status),
	}
}

func getDryRunData(invokedOperation *Operation) any {
	resultStatus := sdk.StatusError.String()
	result := invokedOperation.LatestResult().Value

	data := &sdk.DryRunEventData{
		DeploymentId: int(invokedOperation.DeploymentId),
		OperationId:  invokedOperation.ID,
		Status:       resultStatus,
		Attempts:     invokedOperation.Attempts,
		StartedAt:    invokedOperation.CreatedAt.UTC(),
		CompletedAt:  invokedOperation.UpdatedAt.UTC(),
	}

	if result != nil {
		if dryRunResult, ok := result.(*sdk.DryRunResult); ok {
			data.Status = dryRunResult.Status
			data.Errors = dryRunResult.Errors
		}
	}
	return data
}

func getDeploymentData(invokedOperation *Operation) any {
	var data any

	if invokedOperation.Name == "deploy" {
		deploymentData := &sdk.DeploymentEventData{
			DeploymentId: int(invokedOperation.DeploymentId),
			OperationId:  invokedOperation.ID,
		}

		if invokedOperation.IsRetry() {
			deploymentData.Message = fmt.Sprintf("%s is being retried. Attempt %d of %d", invokedOperation.Name, invokedOperation.Attempts, invokedOperation.Retries)
		} else if invokedOperation.IsRunning() {
			deploymentData.Message = fmt.Sprintf("%s started", invokedOperation.Name)
		}

		if len(invokedOperation.LatestResult().Error) > 0 {
			deploymentData.Message = fmt.Sprintf("%s. Error: %s", deploymentData.Message, invokedOperation.LatestResult().Error)
		}

		data = deploymentData
	}

	data = nil
	return data
}
