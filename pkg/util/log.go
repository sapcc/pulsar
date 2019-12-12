package util

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// NewLogger returns a new Logger with log level configured.
func NewLogger() log.Logger {
	logger := log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	logLevel := level.AllowInfo()
	if v, ok := os.LookupEnv("DEBUG"); ok && v != "false" {
		logLevel = level.AllowDebug()
	}
	level.NewFilter(logger, logLevel)

	return logger
}
