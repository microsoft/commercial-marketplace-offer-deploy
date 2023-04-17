package log

import "testing"

func TestLogger(t *testing.T) {
	myLogger := NewLoggerPublisher()

	//
	myMsg := LogMessage{
		JSONPayload: "hello world, it works!!!!",
	}
	//
	myLogger.Publish(&myMsg, "high")
}

// logging level: Error, info
// service name: operator, apiserver
