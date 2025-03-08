package ports

import "context"

// Logger is the interface that wraps the basic logging methods.
type Logger interface {
	Close() error
	Debug(ctx context.Context, msg string, args ...interface{})
	Info(ctx context.Context, msg string, args ...interface{})
	Warn(ctx context.Context, msg string, args ...interface{})
	Error(ctx context.Context, msg string, args ...interface{})
	Fatal(ctx context.Context, msg string, args ...interface{})
}
