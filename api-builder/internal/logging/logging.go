// Package logging provides logging configurations for the application.
package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func Init(debug bool) *logrus.Logger {
	Logger = logrus.New()

	Logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	Logger.SetOutput(os.Stdout)

	level := logrus.InfoLevel
	if debug {
		level = logrus.DebugLevel
	}
	Logger.SetLevel(level)

	if os.Getenv("LOG_FORMAT") == "json" {
		Logger.SetFormatter(&logrus.JSONFormatter{})
	}

	return Logger
}
