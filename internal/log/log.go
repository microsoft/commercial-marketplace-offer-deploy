package log

import (
	"encoding/json"
	"fmt"
	"os"

	//	log "github.com/sirupsen/logrus"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
	"github.com/sirupsen/logrus"
	
)

type LoggingConfig struct {
	InstrumentationKey string
	DefaultLogLevel    string
}

type appInsightsOptions struct {
	InstrumentationKey string

	Role    string
	Version string
}

// todo, AppConfig drives there
func ConfigureLogging(config *LoggingConfig) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetReportCaller(true)
	//logrus.SetFormatter(&logrus.JSONFormatter{})

	if len(config.InstrumentationKey) == 0 {
		insightsConfig := appInsightsOptions{
			Role:               "MODM",
			Version:            "1.0",
			InstrumentationKey: config.InstrumentationKey,
		}
		hook := &LogrusHook{
			Client: createTelemetryClient(insightsConfig),
		}
		logrus.AddHook(hook)
	}
}

type LogMessage struct {
	Message string
	Level   logrus.Level
}

type LogPublisher interface {
	// publishes a message to all web hook subscriptions
	Publish(message *LogMessage) error
	Log(message string)
	LogWarning(message string)
	LogInfo(message string)
}

// PanicLevel Level = iota
// 	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
// 	// logging level is set to Panic.
// 	FatalLevel
// 	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
// 	// Commonly used for hooks to send errors to an error tracking service.
// 	ErrorLevel
// 	// WarnLevel level. Non-critical entries that deserve eyes.
// 	WarnLevel
// 	// InfoLevel level. General operational entries about what's going on inside the
// 	// application.
// 	InfoLevel
// 	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
// 	DebugLevel
// 	// TraceLevel level. Designates finer-grained informational events than the Debug.
// 	TraceLevel

func (p *appInsightsOptions) Log(message string) {
	p.Publish(&LogMessage{
		Message: message,
		Level:   logrus.InfoLevel,
	})
}

func (p *appInsightsOptions) LogError(message string) {
	p.Publish(&LogMessage{
		Message: message,
		Level:   logrus.ErrorLevel,
	})
}

func (p *appInsightsOptions) LogInfo(message string) {
	p.Publish(&LogMessage{
		Message: message,
		Level:   logrus.InfoLevel,
	})
}

func (p *appInsightsOptions) LogWarning(message string) {
	p.Publish(&LogMessage{
		Message: message,
		Level:   logrus.WarnLevel,
	})
}

func (p *appInsightsOptions) Publish(message *LogMessage) error {
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

func createTelemetryClient(options appInsightsOptions) appinsights.TelemetryClient {
	client := appinsights.NewTelemetryClient(options.InstrumentationKey)

	if len(options.Role) > 0 {
		client.Context().Tags.Cloud().SetRole(options.Role)
	}

	if len(options.Version) > 0 {
		client.Context().Tags.Application().SetVer(options.Version)
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
