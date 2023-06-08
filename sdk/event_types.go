package sdk

// Defines an event that occurs in MODM
type EventTypeName string

// deployment event types
const (
	EventTypeDeploymentCreated   EventTypeName = "deploymentCreated"
	EventTypeDeploymentUpdated   EventTypeName = "deploymentUpdated"
	EventTypeDeploymentDeleted   EventTypeName = "deploymentDeleted"
	EventTypeDeploymentScheduled EventTypeName = "deploymentScheduled"
	EventTypeDeploymentStarted   EventTypeName = "deploymentStarted"
	EventTypeDeploymentCompleted EventTypeName = "deploymentCompleted"
	EventTypeDeploymentRetried   EventTypeName = "deploymentRetried"
)

// stage event types
const (
	EventTypeStageScheduled EventTypeName = "stageScheduled"
	EventTypeStageStarted   EventTypeName = "stageStarted"
	EventTypeStageCompleted EventTypeName = "stageCompleted"
	EventTypeStageRetried   EventTypeName = "stageRetried"
)

// dry run event types
const (
	EventTypeDryRunSchedule  EventTypeName = "dryRunScheduled"
	EventTypeDryRunStarted   EventTypeName = "dryRunStarted"
	EventTypeDryRunCompleted EventTypeName = "dryRunCompleted"
)

// generic event type
const (
	EventTypeDeploymentEventReceived EventTypeName = "deploymentEventReceived"
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
		EventTypeDeploymentScheduled.String(),
		EventTypeDeploymentStarted.String(),
		EventTypeDeploymentCompleted.String(),
		EventTypeDeploymentRetried.String(),
		EventTypeStageScheduled.String(),
		EventTypeStageStarted.String(),
		EventTypeStageCompleted.String(),
		EventTypeStageRetried.String(),
		EventTypeDryRunSchedule.String(),
		EventTypeDryRunStarted.String(),
		EventTypeDryRunCompleted.String(),
		EventTypeDeploymentEventReceived.String(),
	}
}
