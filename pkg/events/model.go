package events

import (
	"github.com/google/uuid"
)

// Defines an event that occurs in MODM
type EventType string

// the list of available / known event types
const (
	DeploymentDryRunCompletedEventType EventType = "deployment.dryruncompleted"
	DeploymentCreatedEventType         EventType = "deployment.created"
	DeploymentStartedEventType         EventType = "deployment.started"
	DeploymentCompletedEventType       EventType = "deployment.completed"
	DeploymentErrorEventType           EventType = "deployment.error"
	DeploymentRetryEventType           EventType = "deployment.retry"
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

// subscription model

type EventSubscriptionMessage struct {
	Id        uuid.UUID `json:"id,omitempty"`
	EventType `json:"eventType,omitempty"`
	Payload   map[string]any `json:"payload,omitempty"`
}
