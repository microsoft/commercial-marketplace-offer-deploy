package events

import (
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// subscription model for MODM webhook events
type EventHookMessage struct {
	// the ID of the message
	Id uuid.UUID `json:"id,omitempty"`

	// the ID of the hook
	HookId uuid.UUID `json:"hookId,omitempty"`

	// the type of the event, .e.g. "dryRunCompleted"
	Type string `json:"type,omitempty"`

	// the status of the event, e.g. "success"
	Status string `json:"status,omitempty"`

	// subject is in format like /deployments/{deploymentId}/stages/{stageId}/operations/{operationName}
	// /deployments/{deploymentId}/operations/{operationName}
	Subject string `json:"subject,omitempty"`
	Data    any    `json:"body,omitempty"`
}

// Dry run data
type DryRunData struct {
	Status         string                 `json:"status,omitempty"`
	AdditionalInfo []DryRunAdditionalInfo `json:"additionalInfo,omitempty"`
}

// Dry run message that's part of the dry run data, containing details of the specific dry run results
type DryRunAdditionalInfo struct {
	Info interface{} `json:"info,omitempty"`
	Type string      `json:"type,omitempty"`
}

// all other deployment events

type DeploymentEventData struct {
	DeploymentId int     `json:"deploymentId,omitempty"`
	StageId      *string `json:"stageId,omitempty"`
	OperationId  *string `json:"operationId,omitempty"`
	Message      string  `json:"message,omitempty"`
}

type RetryDeploymentEventData struct {
	DeploymentEventData
	Attempts int `json:"attempt,omitempty"`
}

func (m *EventHookMessage) DeploymentId() (uint, error) {
	if data, ok := m.Data.(DeploymentEventData); ok {
		return uint(data.DeploymentId), nil
	}

	if m.Subject != "" && strings.HasPrefix(m.Subject, "/deployments/") {
		values := strings.Split(strings.TrimPrefix(m.Subject, "/"), "/")

		deploymentId, err := strconv.Atoi(values[1])
		if err != nil {
			return 0, err
		}
		return uint(deploymentId), nil
	}
	return 0, errors.New("unable to get deployment id using data or the subject")
}

func (m *EventHookMessage) SetSubject(deploymentId uint, stageId *uuid.UUID) {
	m.Subject = "/deployments/" + strconv.Itoa(int(deploymentId))
	if stageId != nil {
		m.Subject += "/stages/" + stageId.String()
	}
}
