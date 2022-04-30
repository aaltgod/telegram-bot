package logger

import "github.com/sirupsen/logrus"

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		DisableColors:   false,
		DisableSorting:  true,
		FullTimestamp:   true,
		TimestampFormat: "2006/01/02 15:04:05",
	}
	return logger
}
