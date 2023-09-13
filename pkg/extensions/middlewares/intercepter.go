package middlewares

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

// Intercepter is a middleware that intercepts messages in reception and in publication.
func Intercepter(ch chan extensions.BrokerMessage) extensions.Middleware {
	return func(ctx context.Context, next extensions.NextMiddleware) context.Context {
		// If there is a broker message, then send it to the channel
		extensions.IfContextSetWith(ctx, extensions.ContextKeyIsBrokerMessage, func(brokerMessage extensions.BrokerMessage) {
			ch <- brokerMessage
		})

		return ctx
	}
}
