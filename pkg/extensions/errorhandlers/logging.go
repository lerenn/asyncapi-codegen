package errorhandlers

import (
	"context"
	"fmt"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

// Logging is an errorhandler that logs errors from processing subscription function
func Logging(logger extensions.Logger) extensions.ErrorHandler {
	return func(ctx context.Context, topic string, msg *extensions.AcknowledgeableBrokerMessage, err error) {
		logger.Error(ctx, fmt.Sprintf("error on processing subscription function for topic: %s - error: %s",
			topic, err.Error()))
	}
}
