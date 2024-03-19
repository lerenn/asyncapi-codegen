package extensions

import "context"

// ErrorHandler is the signature of the function that needs to be implemented to use
// errorhandler functionality
type ErrorHandler func(ctx context.Context, topic string, msg *AcknowledgeableBrokerMessage, err error)
