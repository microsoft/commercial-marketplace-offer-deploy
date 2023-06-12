package mapper

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
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
	setSubject(message, invokedOperation)
	return message
}

func setSubject(m *sdk.EventHookMessage, o *model.InvokedOperation) {
	deploymentId := uint(o.DeploymentId)

	if strings.Contains(strings.ToLower(o.Name), "stage") {
		value, ok := o.Parameters[string(model.ParameterKeyStageId)]
		if ok {
			stageId := uuid.MustParse(value.(string))
			m.SetSubject(deploymentId, to.Ptr(stageId))
		}
	} else {
		m.SetSubject(deploymentId, nil)
	}
}

func getEventType(o *model.InvokedOperation) string {
	noun := ""

	if o.Name == sdk.OperationDeploy.String() || o.Name == sdk.OperationRetry.String() {
		noun = "deployment"
	} else if o.Name == sdk.OperationDryRun.String() {
		noun = "dryRun"
	} else if o.Name == sdk.OperationRetryStage.String() || o.Name == sdk.OperationDeployStage.String() {
		noun = "stage"
	}

	verb := ""

	if o.IsRetry() {
		verb = "Retry"
	}

	if o.IsScheduled() {
		verb += "Scheduled"
	} else if o.IsRunning() {
		verb += "Started"
	} else if o.IsCompleted() {
		verb += "Completed"
	} else {
		verb += "Completed"
	}

	return noun + verb
}

func getEventData(invokedOperation *model.InvokedOperation) any {

	switch invokedOperation.Name {
	case sdk.OperationDeploy.String():
		return getDeploymentData(invokedOperation)
	case sdk.OperationDeployStage.String():
		return getDeployStageData(invokedOperation)
	case sdk.OperationDryRun.String():
		return getDryRunData(invokedOperation)
	case sdk.OperationRetry.String():
		return getRetryData(invokedOperation)
	case sdk.OperationRetryStage.String():
		return getRetryStageData(invokedOperation)
	}
	return nil
}

func getRetryData(invokedOperation *model.InvokedOperation) any {
	data := &sdk.DeploymentEventData{
		EventData: getBaseEventData(invokedOperation),
		Message:   fmt.Sprintf("Retry deployment %s", invokedOperation.Status),
	}
	return data
}

func getDeployStageData(invokedOperation *model.InvokedOperation) any {
	data := &sdk.StageEventData{
		EventData: getBaseEventData(invokedOperation),
	}

	parameter, ok := invokedOperation.ParameterValue(model.ParameterKeyStageId)
	if !ok {
		log.Warnf("StageId parameter not found in invoked operation %s", invokedOperation.Name)
		return data
	}

	stageId, err := uuid.Parse(parameter.(string))
	if err != nil {
		log.Warnf("StageId parameter is not a valid UUID in invoked operation %s", invokedOperation.Name)
		return data
	}

	data.StageId = to.Ptr(stageId)

	latestResult := invokedOperation.LatestResult()
	if latestResult != nil {
		if latestResult.Error != "" {
			data.Message = fmt.Sprintf("Error: %s", latestResult.Error)
		}
	}

	if invokedOperation.ParentID != nil {
		data.ParentOperationId = invokedOperation.ParentID
	}

	return data
}

func getRetryStageData(invokedOperation *model.InvokedOperation) any {
	stageId := uuid.MustParse(invokedOperation.Parameters["stageId"].(string))

	data := &sdk.StageEventData{
		EventData: getBaseEventData(invokedOperation),
		StageId:   to.Ptr(stageId),
	}

	latestResult := invokedOperation.LatestResult()
	if latestResult != nil {
		if latestResult.Error != "" {
			data.Message = fmt.Sprintf("Error: %s", latestResult.Error)
		}
	}

	if invokedOperation.ParentID != nil {
		data.ParentOperationId = invokedOperation.ParentID
	}

	return data
}

func getDryRunData(invokedOperation *model.InvokedOperation) any {
	resultStatus := invokedOperation.Status
	latestResult := invokedOperation.LatestResult()

	data := &sdk.DryRunEventData{
		EventData: getBaseEventData(invokedOperation),
		Status:    resultStatus,
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

	return data
}

func getDeploymentData(invokedOperation *model.InvokedOperation) any {
	data := &sdk.DeploymentEventData{
		EventData: getBaseEventData(invokedOperation),
	}

	if invokedOperation.IsRetry() {
		data.Message = fmt.Sprintf("Operation failed, scheduling a retry. Attempt %d", invokedOperation.Attempts)
	} else if invokedOperation.IsRunning() {
		data.Message = fmt.Sprintf("%s started", invokedOperation.Name)
	}

	if len(invokedOperation.LatestResult().Error) > 0 {
		data.Message = fmt.Sprintf("%s. Error: %s", data.Message, invokedOperation.LatestResult().Error)
	}

	return data
}

func getBaseEventData(invokedOperation *model.InvokedOperation) sdk.EventData {
	base := &sdk.EventData{
		DeploymentId: int(invokedOperation.DeploymentId),
		OperationId:  invokedOperation.ID,
		Attempts:     int(invokedOperation.Attempts),
	}

	setTimestamps(base, invokedOperation)

	return *base
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
