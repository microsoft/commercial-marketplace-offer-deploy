package events

import (
	"strconv"

	"github.com/google/uuid"
)

// Defines an event that occurs in MODM
type EventType string

// the list of available / known event types
const (
	EventTypeCreated   EventType = "Created"
	EventTypeAccepted  EventType = "Accepted"
	EventTypeScheduled EventType = "Scheduled"
	EventTypeRunning   EventType = "Running"
	EventTypeSuccess   EventType = "Success"
	EventTypeFailed    EventType = "Failed"
	EventTypeError     EventType = "Error"
	RetryingEventType  EventType = "Retrying"
)

// Gets the list of events
func GetEventTypes() []string {
	return []string{
		EventTypeCreated.String(),
		EventTypeAccepted.String(),
		EventTypeScheduled.String(),
		EventTypeRunning.String(),
		EventTypeSuccess.String(),
		EventTypeError.String(),
		RetryingEventType.String(),
	}
}

func (o EventType) String() string {
	stringValue := string(o)
	return stringValue
}

// subscription model for MODM webhook events
type EventHookMessage struct {
	// the ID of the message
	Id uuid.UUID `json:"id,omitempty"`

	// the ID of the hook
	HookId uuid.UUID `json:"hookId,omitempty"`

	// subject is in format like /deployments/{deploymentId}/stages/{stageId}/operations/{operationName}
	// /deployments/{deploymentId}/operations/{operationName}
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

type EventHookDeploymentEventMessageBody struct {
	DeploymentId int    `json:"deploymentId,omitempty"`
	Status       string `json:"status,omitempty"`
	Message      string `json:"message,omitempty"`
}

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
