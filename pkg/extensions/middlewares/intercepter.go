package middlewares

import (
	"context"

	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions"
)

// Intercepter is a middleware that intercepts messages in reception and in publication.
func Intercepter(ch chan extensions.BrokerMessage) extensions.Middleware {
	return func(_ context.Context, msg *extensions.BrokerMessage, _ extensions.NextMiddleware) error {
		// Send the message to the channel
		ch <- *msg

		// Do not interrupt the operations
		return nil
	}
}
