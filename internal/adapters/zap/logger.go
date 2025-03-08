package zap

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"log"
)

const (
	FailedLogMsg = `{"msg" : "failed to log event. logger is invalid or un-initialized"}`
)

// Logger provides an implementation for sigma.Logger interface.
// It uses zap.SugaredLogger as the underlying logging implementation.
type Logger struct {
	sugar *zap.SugaredLogger
}

// NewLogger creates a new application logger.
// It chooses the logger type based on the application environment.
// It returns an error if the logger cannot be initialized.
func NewLogger(env string) (*Logger, error) {
	// Choose logger type based on application environment
	conf := zap.NewProductionConfig()
	if env == "development" {
		conf = zap.NewDevelopmentConfig()
	}

	// Override any custom config
	l, err := conf.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	logger := l.Sugar()
	return &Logger{
		sugar: logger,
	}, nil
}

// Close cleans up the logger resources.
// For zap logger it flushes the logger buffer, if any.
// It should generally be called before main exits.
// It returns any error encountered while flushing the buffer.
func (l *Logger) Close() error {
	// Check if logger was initialized
	if l == nil || l.sugar == nil {
		return errors.New("cannot properly close logger. logger is invalid or un-initialized")
	}
	err := l.sugar.Sync()
	return err
}

// Debug implements sigma.Logger.Debug interface by invoking the underlying logging
// implementation for debug level logging
func (l *Logger) Debug(ctx context.Context, msg string, args ...interface{}) {
	if l == nil || l.sugar == nil {
		log.Println(FailedLogMsg)
		return
	}
	l.sugar.Debugw(msg, args...)
}

// Info implements sigma.Logger.Info interface by invoking the underlying logging
// implementation for info level logging
func (l *Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l == nil || l.sugar == nil {
		log.Println(FailedLogMsg)
		return
	}
	l.sugar.Infow(msg, args...)
}

// Warn implements sigma.Logger.Warn interface by invoking the underlying logging
// implementation for warn level logging
func (l *Logger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l == nil || l.sugar == nil {
		log.Println(FailedLogMsg)
		return
	}
	l.sugar.Warnw(msg, args...)
}

// Error implements sigma.Logger.Error interface by invoking the underlying logging
// implementation for error level logging
func (l *Logger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l == nil || l.sugar == nil {
		log.Println(FailedLogMsg)
		return
	}
	l.sugar.Errorw(msg, args...)
}

// Fatal implements sigma.Logger.Fatal interface by invoking the underlying logging
// implementation for fatal level logging
func (l *Logger) Fatal(ctx context.Context, msg string, args ...interface{}) {
	if l == nil || l.sugar == nil {
		log.Println(FailedLogMsg)
		return
	}
	l.sugar.Fatalw(msg, args...)
}
