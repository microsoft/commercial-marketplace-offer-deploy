package log

import "testing"

func TestLogger(t *testing.T) {
	myLogger := NewLogPublisher("", "")
	myMsg := LogMessage{
		JSONPayload: "hello world!!!",
	}
	myLogger.Publish(&myMsg)
}
