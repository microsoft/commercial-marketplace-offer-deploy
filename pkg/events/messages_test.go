package events

import (
	"encoding/json"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetDeploymentIdUsingSubjectOnMessage(t *testing.T) {
	message := EventHookMessage{
		Subject: "/deployments/1234/stages/5678/operations/91011",
	}

	deploymentId, err := message.DeploymentId()
	assert.NoError(t, err)
	assert.Equal(t, uint(1234), deploymentId)
}

func TestGetDeploymentIdUsingDataOnMessage(t *testing.T) {
	message := EventHookMessage{
		Subject: "",
		Data:    DeploymentEventData{DeploymentId: 1},
	}

	deploymentId, err := message.DeploymentId()
	assert.NoError(t, err)
	assert.Equal(t, uint(1), deploymentId)
}

func Test_EventHookMessage_DryRunEventData_Marshaling(t *testing.T) {
	original := EventHookMessage{
		Subject: "",
		Type:    EventTypeDryRunCompleted.String(),
		Data: DryRunEventData{
			DeploymentId: 1,
			OperationId:  uuid.New(),
			Attempts:     1,
			Result: DryRunEventDataResult{
				Status: to.Ptr("failed"),
				Error: &DryRunEventDataError{
					Code:           to.Ptr("code"),
					AdditionalInfo: []*DryRunEventDataErrorAdditionalInfo{},
				},
			},
		},
	}

	bytes, _ := json.Marshal(original)
	jsonString := string(bytes)

	unmarshaled := &EventHookMessage{}
	_ = json.Unmarshal([]byte(jsonString), unmarshaled)
	t.Logf("marshaled: %+v", unmarshaled)

	data, err := unmarshaled.DryRunEventData()
	assert.NoError(t, err)

	assert.Equal(t, 1, data.DeploymentId)
	assert.Equal(t, "code", *data.Result.Error.Code)
}

func Test_EventHookMessage_DryRunEventData_Fails_With_WrongType(t *testing.T) {
	original := EventHookMessage{
		Subject: "",
		Type:    "anything but dryRunCompleted",
		Data: DryRunEventData{
			DeploymentId: 1,
			OperationId:  uuid.New(),
			Attempts:     1,
			Result: DryRunEventDataResult{
				Status: to.Ptr("failed"),
				Error: &DryRunEventDataError{
					AdditionalInfo: []*DryRunEventDataErrorAdditionalInfo{},
				},
			},
		},
	}

	bytes, _ := json.Marshal(original)
	jsonString := string(bytes)

	unmarshaled := &EventHookMessage{}
	_ = json.Unmarshal([]byte(jsonString), unmarshaled)
	t.Logf("marshaled: %+v", unmarshaled)

	_, err := unmarshaled.DryRunEventData()
	assert.Error(t, err)
}

// deployment event data parsing

func Test_EventHookMessage_DeploymentEventData_Marshaling(t *testing.T) {
	types := []EventType{
		EventTypeDeploymentCompleted,
		EventTypeDeploymentCreated,
		EventTypeDeploymentDeleted,
		EventTypeDeploymentEventReceived,
		EventTypeDeploymentRetried,
		EventTypeDeploymentUpdated,
	}

	for _, eventType := range types {

		original := EventHookMessage{
			Subject: "",
			Type:    eventType.String(),
			Data: DeploymentEventData{
				DeploymentId: 1,
				OperationId:  uuid.New(),
				Attempts:     1,
			},
		}

		bytes, _ := json.Marshal(original)
		jsonString := string(bytes)

		unmarshaled := &EventHookMessage{}
		_ = json.Unmarshal([]byte(jsonString), unmarshaled)
		t.Logf("marshaled: %+v", unmarshaled)

		data, err := unmarshaled.DeploymentEventData()
		assert.NoError(t, err)

		assert.Equal(t, 1, data.DeploymentId)
		assert.Equal(t, original.Data.(DeploymentEventData).OperationId, data.OperationId)
	}
}
