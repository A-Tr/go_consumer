package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

// Logger Config includes necessary data to
// create a logger
type LoggerConfig struct {
	Component string
}

// We call this at the start of our program
// to set the logger globally
func Init(o io.Writer, level logrus.Level) {
	logrus.SetOutput(o)
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		DisableColors: true,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "@timestamp",
		},
	})

	// Initialize hooks
	stackHook := NewStackHook()
	logrus.AddHook(stackHook)

}

// New logger returns takes a config a
// returns a logger ready with extra fields
func NewLogger(config LoggerConfig, traceId string) *logrus.Entry {
	fields := logrus.Fields{
		"component":  config.Component,
		"@timestamp": time.Now().UTC().Format(time.RFC3339),
		"traceid":    traceId,
	}

	logger := logrus.WithFields(fields)

	return logger
}
