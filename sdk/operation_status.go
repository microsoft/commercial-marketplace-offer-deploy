package sdk

// Gets the list of events
func GetStatuses() []string {
	return []string{
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
