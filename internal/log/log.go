package log

import (
	"log"
)

type LogMessage struct {
	JSONPayload string
}

type LogPublisher interface {
	// publishes a message to all web hook subscriptions
	Publish(message *LogMessage) error
}

type logPublisher struct {
	logMode string
	sender  string
}

func NewLogPublisher(sender string, logMode string) LogPublisher {
	publisher := &logPublisher{sender: sender, logMode: logMode}

	return publisher
}

func (p *logPublisher) Publish(message *LogMessage) error {
	log.Printf("recieved logged mesage: %s", message)

	// TODO: write to open telemetry which will write to app insights

	return nil
}
