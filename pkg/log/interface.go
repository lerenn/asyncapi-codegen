package log

import "context"

// AdditionalInfo is a key-value pair that will be added to the log
type AdditionalInfo struct {
	Key   string
	Value interface{}
}

// Logger is the interface that must be implemented by a logger
type Interface interface {
	// Info logs information based on a message and key-value elements
	Info(ctx context.Context, msg string, info ...AdditionalInfo)

	// Error logs error based on a message and key-value elements
	Error(ctx context.Context, msg string, info ...AdditionalInfo)
}
