package log

import (
	"encoding/json"
	"fmt"
	"os"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
	"github.com/sirupsen/logrus"
)

const (
	LogPath = "./"
	LogName     = "modmlog"
	LogFileName = LogName + ".txt"
)

type LoggingConfig struct {
	FilePath		   string
	InstrumentationKey string
	DefaultLogLevel    string
}

type InsightsConfig struct {
	InstrumentationKey string

	Role    string
	Version string
}


func ConfigureLogging(config *LoggingConfig) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetReportCaller(true)

	if config == nil {
		fmt.Println("No logging configuration provided")
		return
	}

	if len(config.FilePath) > 0 {
		stacktraceHook := &StacktraceHook{
			innerHook: &FileHook{
				fileName: config.FilePath,
			},
		}
		logrus.AddHook(stacktraceHook)
	}
}

type LogMessage struct {
	Message string
	Level   logrus.Level
}

type LogrusHook struct {
	Client appinsights.TelemetryClient
}

func (hook *LogrusHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *LogrusHook) Fire(entry *logrus.Entry) error {
	if _, ok := entry.Data["message"]; !ok {
		entry.Data["message"] = entry.Message
	}

	level := convertSeverityLevel(entry.Level)
	telemetry := appinsights.NewTraceTelemetry(entry.Message, level)

	for key, value := range entry.Data {
		value = formatData(value)
		telemetry.Properties[key] = fmt.Sprintf("%v", value)
	}

	hook.Client.Track(telemetry)
	return nil
}

func convertSeverityLevel(level logrus.Level) contracts.SeverityLevel {
	switch level {
	case logrus.PanicLevel:
		return contracts.Critical
	case logrus.FatalLevel:
		return contracts.Critical
	case logrus.ErrorLevel:
		return contracts.Error
	case logrus.WarnLevel:
		return contracts.Warning
	case logrus.InfoLevel:
		return contracts.Information
	case logrus.DebugLevel, logrus.TraceLevel:
		return contracts.Verbose
	default:
		return contracts.Information
	}
}

func formatData(value interface{}) (formatted interface{}) {
	switch value := value.(type) {
	case json.Marshaler:
		return value
	case error:
		return value.Error()
	case fmt.Stringer:
		return value.String()
	default:
		return value
	}
}
