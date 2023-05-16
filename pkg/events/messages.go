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
	Data    any    `json:"data,omitempty"`
}

// all other deployment events

type DeploymentEventData struct {
	DeploymentId int        `json:"deploymentId,omitempty" mapstructure:"deploymentId"`
	StageId      *uuid.UUID `json:"stageId,omitempty" mapstructure:"stageId"`
	OperationId  uuid.UUID  `json:"operationId,omitempty" mapstructure:"operationId"`

	// the correlation ID used to track azure deployments. This may be nil if the particular event data is about something that
	// happened inside MODM and not azure.
	CorrelationId *uuid.UUID `json:"correlationId,omitempty" mapstructure:"correlationId"`
	Attempts      int        `json:"attempts,omitempty" mapstructure:"attempts"`
	Message       string     `json:"message,omitempty" mapstructure:"message"`
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
