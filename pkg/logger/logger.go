// Package logger Everything logger related
package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

// DefaultLogger receives the default log of the GIN framework
func (l *Logger) DefaultLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()

		cost := time.Since(start)
		fields := []zapcore.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("domain-errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		}
		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				l.Desugar().Error(e, fields...)
			}
		} else {
			l.Desugar().Info(path, fields...)
		}
	}
}