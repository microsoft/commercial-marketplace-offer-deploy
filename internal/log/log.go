package log

import (
	"fmt"
	"os"

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
	logrus.SetLevel(logrus.DebugLevel)
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

	if len(config.FilePath) > 0 {
		stacktraceHook := &StacktraceHook{
			innerHook: &FileHook{
				fileName: config.FilePath,
			},
		}
		logrus.AddHook(stacktraceHook)
	}
}
