package extensions

import (
	"context"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/errorhandlers"
)

// ErrorHandler is the signature of the function that needs to be implemented to use
// errorhandler functionality
type ErrorHandler func(ctx context.Context, topic string, msg *AcknowledgeableBrokerMessage, err error)

// DefaultErrorHandler returns the default error handler, which is Noop one.
func DefaultErrorHandler() ErrorHandler {
	return errorhandlers.Noop()
}
