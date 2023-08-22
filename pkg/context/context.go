package context

import "context"

const Prefix = "aapi-"

// Key is the type of the keys used in the context
type Key string

const (
	// Module is the name of the module this data is coming from.
	// When coming from a generated client, it is `asyncapi`
	KeyIsModule Key = Prefix + "module"
	// Provider is the name of the provider this data is coming from.
	// When coming from generated code, it is `app`, `client` or `broker`
	KeyIsProvider Key = Prefix + "provider"
	// Action is the name of the action this data is coming from.
	// When coming from generated code, it is the name of the channel
	KeyIsAction Key = Prefix + "action"
	// Operation is the name of the operation this data is coming from.
	// When coming from generated code, it is `subscribe`, `publish`, `wait-for`, etc
	KeyIsOperation Key = Prefix + "operation"
	// Message is the message that has been sent or received
	KeyIsMessage Key = Prefix + "message"
	// CorrelationID is the correlation ID of the message
	KeyIsCorrelationID Key = Prefix + "correlationID"
)

// IfSet executes the function if the key is set in the context
func IfSet(ctx context.Context, key Key, f func(value any)) {
	value := ctx.Value(key)
	if value != nil {
		f(value)
	}
}
