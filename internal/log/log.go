package log

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
)

type LogMessage struct {
	Message string
	Level   logrus.Level
}

type LogPublisher interface {
	// publishes a message to all web hook subscriptions
	Publish(message *LogMessage) error
}

type InsightsConfig struct {
	InstrumentationKey string

	Role    string
	Version string
}

func NewLoggerPublisher() LogPublisher {

	insightsConfig := &InsightsConfig{
		Role:    "NAMEOFYOURAPP",
		Version: "1.0",

		InstrumentationKey: "e2af1774-2ab3-4eca-aa0b-7c75e6e6b8c5",
		// TODO: Move to ENV file
	}

	hook := &LogrusHook{
		Client: createTelemetryClient(insightsConfig),
	}

	logrus.AddHook(hook)

	return insightsConfig
}

func (p *InsightsConfig) Publish(message *LogMessage) error {
	switch message.Level {
	case logrus.PanicLevel:
		logrus.Error(message.Message)
	case logrus.FatalLevel:
		logrus.Error(message.Message)
	case logrus.ErrorLevel:
		logrus.Error(message.Message)
	case logrus.WarnLevel:
		logrus.Warn(message.Message)
	case logrus.InfoLevel:
		logrus.Info(message.Message)
	case logrus.DebugLevel, logrus.TraceLevel:
		logrus.Warn(message.Message)
	default:
		logrus.Info(message.Message)
	}

	return nil
}

func createTelemetryClient(config *InsightsConfig) appinsights.TelemetryClient {
	client := appinsights.NewTelemetryClient(config.InstrumentationKey)

	if len(config.Role) > 0 {
		client.Context().Tags.Cloud().SetRole(config.Role)
	}

	if len(config.Version) > 0 {
		client.Context().Tags.Application().SetVer(config.Version)
	}

	return client
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
