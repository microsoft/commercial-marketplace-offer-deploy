package log

import (
	"bytes"
	"fmt"

	//"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)


func TestLogHook(t *testing.T) {

	env := viper.New()
	env.AddConfigPath("./testdata")
	env.SetConfigName(".env")
	env.SetConfigType("env")

	err := env.ReadInConfig()
	assert.NoError(t, err)

	stacktraceHook := &StacktraceHook{
		innerHook: &FileHook{
			fileName: "./test.log",
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

//	select {}
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
