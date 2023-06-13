package sdk

// Defines an event that occurs in MODM
type EventTypeName string

// deployment event types
const (
	EventTypeDeploymentCreated        EventTypeName = "deploymentCreated"
	EventTypeDeploymentUpdated        EventTypeName = "deploymentUpdated"
	EventTypeDeploymentDeleted        EventTypeName = "deploymentDeleted"
	EventTypeDeploymentScheduled      EventTypeName = "deploymentScheduled"
	EventTypeDeploymentRetryScheduled EventTypeName = "deploymentRetryScheduled"
	EventTypeDeploymentPending        EventTypeName = "deploymentPending"
	EventTypeDeploymentStarted        EventTypeName = "deploymentStarted"
	EventTypeDeploymentRetryStarted   EventTypeName = "deploymentRetryScheduled"
	EventTypeDeploymentCompleted      EventTypeName = "deploymentCompleted"
)

// stage event types
const (
	EventTypeStageScheduled      EventTypeName = "stageScheduled"
	EventTypeStageRetryScheduled EventTypeName = "stageRetryScheduled"
	EventTypeStagePending        EventTypeName = "stagePending"
	EventTypeStageStarted        EventTypeName = "stageStarted"
	EventTypeStageRetryStarted   EventTypeName = "stageRetryStarted"
	EventTypeStageCompleted      EventTypeName = "stageCompleted"
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
		EventTypeDeploymentRetryScheduled.String(),
		EventTypeDeploymentPending.String(),
		EventTypeDeploymentStarted.String(),
		EventTypeDeploymentRetryStarted.String(),
		EventTypeDeploymentCompleted.String(),
		EventTypeStageScheduled.String(),
		EventTypeStageRetryScheduled.String(),
		EventTypeStageStarted.String(),
		EventTypeStageRetryStarted.String(),
		EventTypeStageCompleted.String(),
		EventTypeDryRunSchedule.String(),
		EventTypeDryRunStarted.String(),
		EventTypeDryRunCompleted.String(),
		EventTypeDeploymentEventReceived.String(),
	}
}
