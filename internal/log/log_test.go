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

type StacktraceHook struct {
	innerHook logrus.Hook
}

func (h *StacktraceHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *StacktraceHook) Fire(e *logrus.Entry) error {
	if v, found := e.Data[logrus.ErrorKey]; found {
		if err, iserr := v.(error); iserr {
			type stackTracer interface {
				StackTrace() errors.StackTrace
			}
			if st, isst := err.(stackTracer); isst {
				stack := fmt.Sprintf("%+v", st.StackTrace())
				e.Data["stacktrace"] = stack
			}
		}
	}
	h.innerHook.Fire(e)
	return nil
}

func TestLogHook(t *testing.T) {

	env := viper.New()
	env.AddConfigPath("./testdata")
	env.SetConfigName(".env")
	env.SetConfigType("env")

	err := env.ReadInConfig()
	assert.NoError(t, err)

	loggingConfig := &LoggingConfig{
		InstrumentationKey: env.GetString("LOGGING_APP_INSIGHTS_KEY"),
		DefaultLogLevel:    "Info",
	}

	insightsConfig := InsightsConfig{
		Role:               "MODM",
		Version:            "1.0",
		InstrumentationKey: loggingConfig.InstrumentationKey,
	}

	stacktraceHook := &StacktraceHook{
		innerHook: &LogrusHook{
			Client: createTelemetryClient(insightsConfig),
		},
	}

	var output bytes.Buffer
	logrus.SetOutput(&output)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{DisableQuote: true})
	logrus.AddHook(stacktraceHook)

	logrus.WithError(errors.New("Foo")).Error("Wrong")

	outputString := output.String()
	fmt.Println(outputString)

	select {}
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
