// Package logger for logging

package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

var log *Logger

func init() {
	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.JSONFormatter{})
	logrusLogger.SetOutput(os.Stdout)
	logrusLogger.SetLevel(logrus.InfoLevel)
	log = &Logger{logrusLogger}
}

// NewLogger returns a new logger instance
func NewLogger() *Logger {
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
