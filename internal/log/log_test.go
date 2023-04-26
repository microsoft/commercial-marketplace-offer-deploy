package log

import (
	"bytes"
	"fmt"

	//"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func TestLogHook(t *testing.T) {
	stacktraceHook := &StacktraceHook{
		innerHook: &FileHook{
			fileName: "./testdata/test.log",
		},
	}

	var output bytes.Buffer
	logrus.SetOutput(&output)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.AddHook(stacktraceHook)

	logrus.WithError(errors.New("Foo")).Error("Wrong")

	outputString := output.String()
	fmt.Println(outputString)

	// select {}
}

func TestLogFormatter(t *testing.T) {
	var output bytes.Buffer
	logrus.SetOutput(&output)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.Error(errors.New("This was a great be failure"))

	outputString := output.String()
	fmt.Println(outputString)
}
