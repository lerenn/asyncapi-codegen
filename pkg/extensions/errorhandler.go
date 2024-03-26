package extensions

import (
	"context"
)

// ErrorHandler is the signature of the function that needs to be implemented to use
// errorhandler functionality
type ErrorHandler func(ctx context.Context, topic string, msg *AcknowledgeableBrokerMessage, err error)

// DefaultErrorHandler returns the default error handler, which is a Noop errorhandler.
func DefaultErrorHandler() ErrorHandler {
	return func(ctx context.Context, topic string, msg *AcknowledgeableBrokerMessage, err error) {}
}
