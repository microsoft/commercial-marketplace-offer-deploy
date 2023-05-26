package operation

import (
	"fmt"

	"github.com/microsoft/commercial-marketplace-offer-deploy/sdk"
)

// performs mapping of the invoked operation to the correct event hook message

func mapToMessage(invokedOperation *InvokedOperation) (*sdk.EventHookMessage, error) {
	message := &sdk.EventHookMessage{
		Status: invokedOperation.Status,
	}
	message.SetSubject(uint(invokedOperation.DeploymentId), nil)

	data, err := getEventMessageData(invokedOperation)
	if err != nil {
		return nil, err
	}

	message.Data = data

	return message, nil

}

func getEventMessageData(invokedOperation *InvokedOperation) (any, error) {
	var data any

	if invokedOperation.Name == "deploy" {
		deploymentData := &sdk.DeploymentEventData{
			DeploymentId: int(invokedOperation.DeploymentId),
			OperationId:  invokedOperation.ID,
		}

		if invokedOperation.IsRetry() {
			deploymentData.Message = fmt.Sprintf("Retrying %s. Attempt %d of %d", invokedOperation.Name, invokedOperation.Attempts, invokedOperation.Retries)
		}

		data = deploymentData
	}

	return data, nil
}
