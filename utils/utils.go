package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Contains returns true if the given string is contained in the given slice of strings.
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// configureLogLevel If an existing log level environment variable is present, re-use that to configure logrus.
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
