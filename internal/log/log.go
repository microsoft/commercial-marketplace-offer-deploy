package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	LogPath     = "./"
	LogName     = "modmlog"
	LogFileName = LogName + ".txt"
)

type LoggingOptions struct {
	FilePath        string
	DefaultLogLevel string
}

type InsightsConfig struct {
	InstrumentationKey string

	Role    string
	Version string
}

func ConfigureLogging(config *LoggingOptions) {	
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)

	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	logrus.SetFormatter(formatter)

	if config == nil {
		fmt.Println("No logging configuration provided")
		return
	}

	fmt.Println("default log level: ", config.DefaultLogLevel)

	logLevel := logrus.TraceLevel
	if len(config.DefaultLogLevel) > 0 {
		logLevel = getLogorusLevel(config.DefaultLogLevel)
	}
	logrus.SetLevel(logLevel)

	fmt.Println("log level set to: ", logLevel)

	if len(config.FilePath) > 0 {
		stacktraceHook := &StacktraceHook{
			innerHook: &FileHook{
				fileName: config.FilePath,
			},
		}
		logrus.AddHook(stacktraceHook)
	}
}

func getLogorusLevel(level string) logrus.Level {
    switch strings.ToLower(level) {
    case "debug":
        return logrus.DebugLevel
    case "info":
        return logrus.InfoLevel
    case "warn":
        return logrus.WarnLevel
    case "error":
        return logrus.ErrorLevel
	case "trace":
		return logrus.TraceLevel
    default:
        return logrus.InfoLevel
    }
}
