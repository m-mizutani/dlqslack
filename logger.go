package dlqslack

import "github.com/sirupsen/logrus"

// Logger is configurable from external.
var Logger = logrus.New()

// SetLogLevel should be used to change log level by environment variable
func SetLogLevel(level string) {
	switch level {
	case "TRACE":
		Logger.SetLevel(logrus.TraceLevel)
	case "DEBUG":
		Logger.SetLevel(logrus.DebugLevel)
	case "INFO":
		Logger.SetLevel(logrus.InfoLevel)
	case "WARN":
		Logger.SetLevel(logrus.WarnLevel)
	case "ERROR":
		Logger.SetLevel(logrus.ErrorLevel)
	}
}
