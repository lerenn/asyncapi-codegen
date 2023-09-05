package middlewares

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

// Logging is a middleware that logs messages in reception and in publication
func Logging(logger extensions.Logger) extensions.Middleware {
	return func(ctx context.Context, next extensions.NextMiddleware) context.Context {
		extensions.IfContextValueEquals(ctx, extensions.ContextKeyIsMessageDirection, "reception", func() {
			logger.Info(ctx, "Received a message")
		})

		extensions.IfContextValueEquals(ctx, extensions.ContextKeyIsMessageDirection, "publication", func() {
			logger.Info(ctx, "Published a message")
		})

		return ctx
	}
}
