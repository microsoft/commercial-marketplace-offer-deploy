package events

// Defines an event that occurs in MODM
type DeploymentEvent string

const (
	DeploymentDryRunCompletedEvent DeploymentEvent = "DeploymentDryRunCompleted"
	DeploymentCreatedEvent         DeploymentEvent = "DeploymentCreated"
	DeploymentStartedEvent         DeploymentEvent = "DeploymentStarted"
	DeploymentCompleted            DeploymentEvent = "DeploymentCompleted"
)

// Gets the list of events
func GetEvents() []string {
	return []string{
		DeploymentDryRunCompletedEvent.String(),
		DeploymentCreatedEvent.String(),
		DeploymentStartedEvent.String(),
		DeploymentCompleted.String(),
	}
}

func (o DeploymentEvent) String() string {
	stringValue := string(o)
	return stringValue
}
