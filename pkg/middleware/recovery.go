package middleware

import (
	"context"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/log"
)

// Recovery is a middleware that recovers from panic in middlewares coming after
// it and user code from subscription.
func Recovery(logger log.Interface) Middleware {
	return func(ctx context.Context, next Next) context.Context {
		// Recover in case of panic
		defer func() {
			if r := recover(); r != nil {
				logger.Error(ctx, fmt.Sprintf("Recovered from panic: %v", r))
			}
		}()

		// Call next middleware
		next(ctx)

		return ctx
	}
}
