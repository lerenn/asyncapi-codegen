package middlewares

import (
	"context"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

// Recovery is a middleware that recovers from panic in middlewares coming after
// it and user code from subscription.
func Recovery(logger extensions.Logger) extensions.Middleware {
	return func(ctx context.Context, msg *extensions.BrokerMessage, next extensions.NextMiddleware) error {
		// Recover in case of panic
		defer func() {
			if r := recover(); r != nil {
				logger.Error(ctx, fmt.Sprintf("Recovered from panic: %v", r))
			}
		}()

		// Call next middleware
		return next()
	}
}
