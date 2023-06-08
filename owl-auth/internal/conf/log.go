package conf

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Initialize configures the global logrus logger instance
func Initialize(isProduction bool) {
	var level logrus.Level
	var formatter logrus.Formatter

	if isProduction {
		level = logrus.InfoLevel
		formatter = &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "msg",
			},
			TimestampFormat: "2006-01-02T15:04:05Z",
		}
	} else {
		level = logrus.DebugLevel
		formatter = &logrus.TextFormatter{
			ForceColors:     true,
			TimestampFormat: "2006-01-02 15:04:05.000",
			FullTimestamp:   true,
		}
	}
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(level)
	logrus.SetFormatter(formatter)
	logrus.SetReportCaller(true)
}
