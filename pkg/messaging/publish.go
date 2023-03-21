package messaging

type Publisher interface {
	Publish(message DeploymentMessage) error
}

type DeploymentMessage struct {
	Header DeploymentMessageHeader  `json:"header"`
	Body interface{} `json:"body"`
}

type DeploymentMessageHeader struct {
	Topic string `json:"topic"`
}

