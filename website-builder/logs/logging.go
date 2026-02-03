// Package logs
package logs

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	Logging *logrus.Logger
}

func Init(debug bool) *Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)

	level := logrus.InfoLevel
	if debug {
		level = logrus.DebugLevel
	}
	logger.SetLevel(level)

	if os.Getenv("LOG_FORMAT") == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	return &Logger{
		Logging: logger,
	}
}

func (l *Logger) ErrorLogger(err error, message string) error {
	l.Logging.WithError(err).Error(message)
	return fmt.Errorf("%s", message)
}

func (l *Logger) ErrorMessage(message string) error {
	l.Logging.Println(message)
	return fmt.Errorf("%s", message)
}
