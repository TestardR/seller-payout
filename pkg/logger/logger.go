package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -package=mock -source=logger.go -destination=$MOCK_FOLDER/logger.go Logger

// Logger is the Log interface.
type Logger interface {
	Info(...interface{})
	Error(...interface{})
	Fatal(...interface{})
}

// New returns a new client instance.
func New(appName string) Logger {
	var log = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.JSONFormatter{},
		Level:     logrus.InfoLevel,
	}

	entry := log.WithFields(logrus.Fields{
		"appname": appName,
	})

	return entry
}
