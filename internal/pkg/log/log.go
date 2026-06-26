package log

import (
	"github.com/sirupsen/logrus"
)

var logger = logrus.StandardLogger()

func Error(args ...any) {
	logger.Error(args...)
}

func Errorf(format string, args ...any) {
	logger.Errorf(format, args...)
}

func Info(args ...any) {
	logger.Info(args...)
}

func Infof(format string, args ...any) {
	logger.Infof(format, args...)
}
