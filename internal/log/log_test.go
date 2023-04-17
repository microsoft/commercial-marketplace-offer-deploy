package log

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLogger(t *testing.T) {
	myLogger := NewLoggerPublisher()

	//
	myMsg := LogMessage{
		Message: "317 test",
		Level:   logrus.WarnLevel,
	}
	//
	myLogger.Publish(&myMsg)
}

// logging level: Error, info
// service name: operator, apiserver
