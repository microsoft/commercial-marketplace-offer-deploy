package events

import (
	"strconv"

	"github.com/google/uuid"
)

// Defines an event that occurs in MODM
type EventType string

// the list of available / known event types
const (
	DeploymentDryRunCompletedEventType EventType = "DryRunCompleted"
	DeploymentCreatedEventType         EventType = "Created"
	DeploymentPendingEventType         EventType = "Pending"
	DeploymentStartingEventType        EventType = "Starting"
	DeploymentStartedEventType         EventType = "Started"
	DeploymentCompletedEventType       EventType = "Completed"
	DeploymentErrorEventType           EventType = "Error"
	DeploymentRetryEventType           EventType = "Retry"
)

// Gets the list of events
func GetEventTypes() []string {
	return []string{
		DeploymentDryRunCompletedEventType.String(),
		DeploymentCreatedEventType.String(),
		DeploymentStartedEventType.String(),
		DeploymentCompletedEventType.String(),
		DeploymentErrorEventType.String(),
		DeploymentRetryEventType.String(),
	}
}

func (o EventType) String() string {
	stringValue := string(o)
	return stringValue
}

// subscription model for MODM webhook events
type EventHookMessage struct {
	Id     uuid.UUID `json:"id,omitempty"`
	HookId uuid.UUID `json:"hookId,omitempty"`

	// subject is in format like /deployments/{deploymentId}/stages/{stageId}
	Subject   string `json:"subject,omitempty"`
	EventType string `json:"eventType,omitempty"`
	Body      any    `json:"body,omitempty"`
}

// Dry run
type DryRunCompletedBody struct {
	Messages []EventHookDryRunMessage `json:"messages,omitempty"`
}

type EventHookDryRunMessage struct {
	Type    string `json:"type,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

// all other deployment events

type EventHookDeploymentMessageBody struct {
	DeploymentId int    `json:"deploymentId,omitempty"`
	ResourceId   string `json:"resourceId,omitempty"`
	Status       string `json:"status,omitempty"`
	Message      string `json:"message,omitempty"`
}

func (m *EventHookMessage) SetSubject(deploymentId int, stageId *uuid.UUID, resourceName *string) {
	m.Subject = "/deployments/" + strconv.Itoa(deploymentId)
	if stageId != nil {
		m.Subject += "/stages/" + stageId.String()
	}
	if resourceName != nil {
		m.Subject += "/resources/" + *resourceName
	}
}
