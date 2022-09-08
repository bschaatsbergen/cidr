package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// ConfigureLogLevel If an existing log level environment variable is present, use that to configure logrus.
func ConfigureLogLevel(debugLogsEnabled bool) {
	logLevelStr, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		logLevelStr = "info"
	}
	if debugLogsEnabled {
		logLevelStr = "debug"
	}
	logLevel, err := logrus.ParseLevel(logLevelStr)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)
}
