package mapper

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

// performs mapping of the invoked operation to the correct event hook message

func MapInvokedOperation(invokedOperation *model.InvokedOperation) *sdk.EventHookMessage {
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
	} else if o.IsRetry() {
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

func getRetryData(invokedOperation *model.InvokedOperation) any {
	data := &sdk.DeploymentEventData{
		EventData: sdk.EventData{
			DeploymentId: int(invokedOperation.DeploymentId),
			OperationId:  invokedOperation.ID,
		},
		Message: fmt.Sprintf("Retry deployment %s", invokedOperation.Status),
	}
	setTimestamps(&data.EventData, invokedOperation)

	return data
}

func getRetryStageData(invokedOperation *model.InvokedOperation) any {
	data := &sdk.DeploymentEventData{
		EventData: sdk.EventData{
			DeploymentId: int(invokedOperation.DeploymentId),
			OperationId:  invokedOperation.ID,
		},
		StageId: to.Ptr(uuid.MustParse(invokedOperation.Parameters["stageId"].(string))),
		Message: fmt.Sprintf("Retry deployment %s", invokedOperation.Status),
	}
	setTimestamps(&data.EventData, invokedOperation)
	return data
}

func getDryRunData(invokedOperation *model.InvokedOperation) any {
	resultStatus := invokedOperation.Status
	latestResult := invokedOperation.LatestResult()

	data := &sdk.DryRunEventData{
		EventData: sdk.EventData{
			DeploymentId: int(invokedOperation.DeploymentId),
			OperationId:  invokedOperation.ID,
			Attempts:     int(invokedOperation.Attempts),
		},
		Status: resultStatus,
	}

	if latestResult != nil {
		result := latestResult.Value
		if result != nil {
			if dryRunResult, ok := result.(*sdk.DryRunResult); ok {
				data.Status = dryRunResult.Status
				data.Errors = dryRunResult.Errors
			}
		}
	}

	setTimestamps(&data.EventData, invokedOperation)
	return data
}

func getDeploymentData(invokedOperation *model.InvokedOperation) any {
	data := &sdk.DeploymentEventData{
		EventData: sdk.EventData{
			DeploymentId: int(invokedOperation.DeploymentId),
			OperationId:  invokedOperation.ID,
		},
	}

	if invokedOperation.IsRetry() {
		data.Message = fmt.Sprintf("%s is being retried. Attempt %d of %d", invokedOperation.Name, invokedOperation.Attempts, invokedOperation.Retries)
	} else if invokedOperation.IsRunning() {
		data.Message = fmt.Sprintf("%s started", invokedOperation.Name)
	}

	if len(invokedOperation.LatestResult().Error) > 0 {
		data.Message = fmt.Sprintf("%s. Error: %s", data.Message, invokedOperation.LatestResult().Error)
	}

	setTimestamps(&data.EventData, invokedOperation)

	return data
}

func setTimestamps(data *sdk.EventData, operation *model.InvokedOperation) {
	if operation.IsScheduled() {
		data.ScheduledAt = to.Ptr(operation.CreatedAt.UTC())
	} else if operation.IsRunning() {
		latest := operation.LatestResult()
		if latest != nil {
			data.StartedAt = to.Ptr(latest.StartedAt.UTC())
		}
	} else if operation.IsCompleted() {
		first := operation.FirstResult()
		if first != nil {
			data.StartedAt = to.Ptr(first.StartedAt.UTC())
		}
		latest := operation.LatestResult()
		if latest != nil {
			data.CompletedAt = to.Ptr(latest.CompletedAt.UTC())
		}
	}
}
