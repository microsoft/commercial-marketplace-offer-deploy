package events

// Defines an event that occurs in MODM
type Event string

const (
	DeploymentDryRunCompletedEvent Event = "DeploymentDryRunCompleted"
	DeploymentCreatedEvent         Event = "DeploymentCreated"
	DeploymentStartedEvent         Event = "DeploymentStarted"
	DeploymentCompleted            Event = "DeploymentCompleted"
)
