package log

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLogger(t *testing.T) {
	myLogger := NewLoggerPublisher()

	//
	myMsg := LogMessage{
		Message: "Testing from TestLogger()",
		Level:   logrus.WarnLevel,
	}

	myLogger.Publish(&myMsg)
}
