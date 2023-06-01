package operation

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

// performs mapping of the invoked operation to the correct event hook message

func mapToMessage(invokedOperation *model.InvokedOperation) *sdk.EventHookMessage {
	message := &sdk.EventHookMessage{
		Status: invokedOperation.Status,
		Type:   getEventType(invokedOperation),
		Error:  invokedOperation.LatestResult().Error,
		Data:   getEventData(invokedOperation),
	}
	message.SetSubject(uint(invokedOperation.DeploymentId), nil)

	return message
}

func getEventType(o *model.InvokedOperation) string {
	noun := ""

	if o.Name == sdk.OperationDeploy.String() || o.Name == sdk.OperationRetry.String() {
		noun = "deployment"
	} else if o.Name == sdk.OperationDryRun.String() {
		noun = "dryRun"
	} else if o.Name == sdk.OperationRetryStage.String() {
		noun = "stage"
	}

	verb := ""

	if o.IsScheduled() {
		verb = "Scheduled"
	} else if o.IsRunning() {
		verb = "Started"
	} else if o.IsCompleted() {
		verb = "Completed"
	} else if o.IsRetriable() {
		verb = "Retried"
	} else {
		verb = "Completed"
	}

	return noun + verb
}

func getEventData(invokedOperation *model.InvokedOperation) any {
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

func getRetryData(operation *model.InvokedOperation) any {
	return &sdk.DeploymentEventData{
		DeploymentId: int(operation.DeploymentId),
		OperationId:  operation.ID,
		Message:      fmt.Sprintf("Retry deployment %s", operation.Status),
	}
}

func getRetryStageData(operation *model.InvokedOperation) any {
	return &sdk.DeploymentEventData{
		DeploymentId: int(operation.DeploymentId),
		StageId:      to.Ptr(uuid.MustParse(operation.Parameters["stageId"].(string))),
		OperationId:  operation.ID,
		Message:      fmt.Sprintf("Retry deployment %s", operation.Status),
	}
}

func getDryRunData(invokedOperation *model.InvokedOperation) any {
	resultStatus := sdk.StatusError.String()
	result := invokedOperation.LatestResult().Value

	data := &sdk.DryRunEventData{
		DeploymentId: int(invokedOperation.DeploymentId),
		OperationId:  invokedOperation.ID,
		Status:       resultStatus,
		Attempts:     int(invokedOperation.Attempts),
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

func getDeploymentData(invokedOperation *model.InvokedOperation) any {
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
