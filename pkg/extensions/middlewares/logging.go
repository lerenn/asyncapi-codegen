package middlewares

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

// Logging is a middleware that logs messages in reception and in publication.
func Logging(logger extensions.Logger) extensions.Middleware {
	return func(ctx context.Context, msg *extensions.BrokerMessage, _ extensions.NextMiddleware) error {
		// Log if this is a received message
		extensions.IfContextValueEquals(ctx, extensions.ContextKeyIsDirection, "reception", func() {
			logger.Info(ctx, "Received a message")
		})

		// Log if this is a published message
		extensions.IfContextValueEquals(ctx, extensions.ContextKeyIsDirection, "publication", func() {
			logger.Info(ctx, "Published a message")
		})

		// Do not interrupt the operations
		return nil
	}
}
