package middleware

import (
	"context"

	apiContext "github.com/lerenn/asyncapi-codegen/pkg/context"
	"github.com/lerenn/asyncapi-codegen/pkg/log"
)

// Logging is a middleware that logs messages in reception and in publication
func Logging(logger log.Interface) Middleware {
	return func(ctx context.Context, next Next) context.Context {
		apiContext.IfEquals(ctx, apiContext.KeyIsMessageDirection, "reception", func() {
			logger.Info(ctx, "Received a message")
		})

		apiContext.IfEquals(ctx, apiContext.KeyIsMessageDirection, "publication", func() {
			logger.Info(ctx, "Published a message")
		})

		return ctx
	}
}
