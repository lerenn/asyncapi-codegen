package errorhandlers

import (
	"context"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

// Noop is an errorhandler as default implementation that does nothing
func Noop() extensions.ErrorHandler {
	return func(ctx context.Context, topic string, msg *extensions.AcknowledgeableBrokerMessage, err error) {}
}
