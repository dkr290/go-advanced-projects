package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}

// NewLogger returns a new logger instance
func NewLogger() *logrus.Logger {
	return log
}

// SetLogLevel sets the log level
func SetLogLevel(level string) {
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
}

// SetTextFormatter sets the logger to use text format instead of JSON
func SetTextFormatter() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
