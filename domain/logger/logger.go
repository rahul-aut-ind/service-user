// Package logger Everything logger related
package logger

import (
	"go.uber.org/zap"
)

type ILogger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

type Logger struct {
	*zap.SugaredLogger
}

// New Creates new instance of Logger
func New() *Logger {
	var logger *zap.Logger

	logger, _ = zap.NewDevelopment()
	sugarLogger := logger.Sugar()

	return &Logger{sugarLogger}
}
