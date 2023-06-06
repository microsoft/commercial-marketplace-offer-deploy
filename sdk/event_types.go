package sdk

// Defines an event that occurs in MODM
type EventTypeName string

// the list of available / known event types
const (
	EventTypeDeploymentCreated   EventTypeName = "deploymentCreated"
	EventTypeDeploymentUpdated   EventTypeName = "deploymentUpdated"
	EventTypeDeploymentDeleted   EventTypeName = "deploymentDeleted"
	EventTypeDeploymentScheduled EventTypeName = "deploymentScheduled"
	EventTypeDeploymentStarted   EventTypeName = "deploymentStarted"
	EventTypeDeploymentCompleted EventTypeName = "deploymentCompleted"
	EventTypeStageCompleted      EventTypeName = "stageCompleted"
	EventTypeStageRetried        EventTypeName = "stageRetried"
	EventTypeStageStarted        EventTypeName = "stageStarted"
	EventTypeStageScheduled      EventTypeName = "stageScheduled"
	EventTypeDeploymentRetried   EventTypeName = "deploymentRetried"
	EventTypeDryRunCompleted     EventTypeName = "dryRunCompleted"

	EventTypeDeploymentOperationReceived EventTypeName = "deploymentOperationReceived"
	EventTypeDeploymentEventReceived     EventTypeName = "deploymentEventReceived"
)

func (e EventTypeName) String() string {
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
		EventTypeDeploymentOperationReceived.String(),
		EventTypeDeploymentEventReceived.String(),
	}
}
