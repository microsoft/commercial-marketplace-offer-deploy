package events

// Defines an event that occurs in MODM
type EventType string

// the list of available / known event types
const (
	EventTypeDeploymentCreated   EventType = "deploymentCreated"
	EventTypeDeploymentUpdated   EventType = "deploymentUpdated"
	EventTypeDeploymentDeleted   EventType = "deploymentDeleted"
	EventTypeDeploymentStarted   EventType = "deploymentStarted"
	EventTypeDeploymentCompleted EventType = "deploymentCompleted"
	EventTypeStageCompleted      EventType = "stageCompleted"

	EventTypeDeploymentRetried EventType = "deploymentRetried"
	EventTypeDryRunCompleted   EventType = "dryRunCompleted"
	EventTypeDryRunRetrying    EventType = "dryRunRetrying"
	EventTypeDryRunRetried     EventType = "dryRunRetried"

	EventTypeDeploymentOperationReceived EventType = "deploymentOperationReceived"
	EventTypeDeploymentEventReceived     EventType = "deploymentEventReceived"
)

func (e EventType) String() string {
	stringValue := string(e)
	return stringValue
}

func GetEventTypes() []string {
	return []string{
		EventTypeDeploymentCreated.String(),
		EventTypeDeploymentUpdated.String(),
		EventTypeDeploymentDeleted.String(),
		EventTypeDeploymentStarted.String(),
		EventTypeDeploymentCompleted.String(),
		EventTypeDeploymentRetried.String(),
		EventTypeDryRunCompleted.String(),
		EventTypeDryRunRetrying.String(),
		EventTypeDryRunRetried.String(),
		EventTypeDeploymentOperationReceived.String(),
		EventTypeDeploymentEventReceived.String(),
	}
}
