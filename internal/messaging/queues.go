package messaging

type QueueName string

const (
	QueueNameEvents     QueueName = "events"
	QueueNameOperations QueueName = "operations"
)
