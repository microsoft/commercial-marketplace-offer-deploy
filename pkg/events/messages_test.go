package events

import (
	"testing"

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
