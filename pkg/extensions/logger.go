package extensions

import "context"

// LogInfo is a key-value pair that will be added to the log.
type LogInfo struct {
	Key   string
	Value any
}

// Logger is the interface that must be implemented by a logger.
type Logger interface {
	// Info logs information based on a message and key-value elements
	Info(ctx context.Context, msg string, info ...LogInfo)

	// Warning logs information based on a message and key-value elements
	// This levels indicates a non-expected state but that does not prevent the
	// application to work properly
	Warning(ctx context.Context, msg string, info ...LogInfo)

	// Error logs error based on a message and key-value elements
	Error(ctx context.Context, msg string, info ...LogInfo)
}

// DummyLogger is a logger that does not log anything.
type DummyLogger struct {
}

// Info logs information based on a message and key-value elements.
func (dl DummyLogger) Info(_ context.Context, _ string, _ ...LogInfo) {}

// Warning logs information based on a message and key-value elements.
func (dl DummyLogger) Warning(_ context.Context, _ string, _ ...LogInfo) {}

// Error logs error based on a message and key-value elements.
func (dl DummyLogger) Error(_ context.Context, _ string, _ ...LogInfo) {}
