package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	info  = "info"
	debug = "debug"
)

// ConfigureLogLevel checks if an existing log level environment variable is set,
// otherwise it sets the log level to info.
func ConfigureLogLevel(enableDebugLogging bool) error {
	logLevelStr, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		logLevelStr = info
	}
	if enableDebugLogging {
		logLevelStr = debug
	}
	logLevel, err := logrus.ParseLevel(logLevelStr)
	if err != nil {
		return err
	}
	logrus.SetLevel(logLevel)
	return nil
}
