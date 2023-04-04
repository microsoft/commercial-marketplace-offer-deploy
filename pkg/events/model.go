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

// subscription model for MODM webhook events
type WebHookEventMessage struct {
	Id             uuid.UUID `json:"id,omitempty"`
	SubscriptionId uuid.UUID `json:"subscriptionId,omitempty"`
	DeploymentId   int32     `json:"deploymentId,omitempty"`
	StageId        uuid.UUID `json:"stageId,omitempty"`
	EventType      string    `json:"eventType,omitempty"`
	Body           any       `json:"body,omitempty"`
}

// Dry run
type WebHookDryRunCompletedBody struct {
	Messages []WebHookDryRunMessage `json:"messages,omitempty"`
}

type WebHookDryRunMessage struct {
	Type    string `json:"type,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

// all other deployment events

type WebHookDeploymentEventMessageBody struct {
	ResourceId string `json:"ResourceId,omitempty"`
	Status     int32  `json:"status,omitempty"`
	Message    string `json:"message,omitempty"`
}
