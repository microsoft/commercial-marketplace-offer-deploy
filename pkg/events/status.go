package events

// Defines an event that occurs in MODM
type Status string

// the list of available / known event types
const (
	StatusCreated   Status = "created"
	StatusDeleted   Status = "deleted"
	StatusAccepted  Status = "accepted"
	StatusScheduled Status = "scheduled"
	StatusRunning   Status = "running"
	StatusSuccess   Status = "success"
	StatusFailed    Status = "failed"
	StatusError     Status = "error"
)

// Gets the list of events
func GetStatuses() []string {
	return []string{
		StatusCreated.String(),
		StatusAccepted.String(),
		StatusScheduled.String(),
		StatusRunning.String(),
		StatusSuccess.String(),
		StatusError.String(),
	}
}

func (o Status) String() string {
	stringValue := string(o)
	return stringValue
}
