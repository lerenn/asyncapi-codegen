package extensions

import "context"

// Prefix is the prefix used for all context keys in order to avoid collision
// with other keys that can be present in context.
const Prefix = "asyncapi-"

// ContextKey is the type of the keys used in the context.
type ContextKey string

const (
	// ContextKeyIsProvider is the name of the provider this data is coming from.
	// When coming from a generated user, it is `asyncapi`.
	ContextKeyIsProvider ContextKey = Prefix + "provider"
	// ContextKeyIsChannel is the name of the channel this data is coming from.
	ContextKeyIsChannel ContextKey = Prefix + "channel"
	// ContextKeyIsMessageDirection is the direction this data is coming from.
	// It can be either "publication" or "reception".
	ContextKeyIsMessageDirection ContextKey = Prefix + "operation"
	// ContextKeyIsBrokerMessage is the message that has been sent or received from/to the broker.
	ContextKeyIsBrokerMessage ContextKey = Prefix + "broker-message"
	// ContextKeyIsMessage is the message that has been sent or received.
	ContextKeyIsMessage ContextKey = Prefix + "message"
	// ContextKeyIsCorrelationID is the correlation ID of the message.
	ContextKeyIsCorrelationID ContextKey = Prefix + "correlationID"
)

// String returns the string representation of the key.
func (k ContextKey) String() string {
	return string(k)
}

// IfContextSetWith executes the function if the key is set in the context as a string.
func IfContextSetWith[T any](ctx context.Context, key ContextKey, fn func(value T)) {
	// Get value
	value := ctx.Value(key)
	if value == nil {
		return
	}

	// Get value as type T
	if tValue, ok := value.(T); ok {
		fn(tValue)
	}
}

// IfContextValueEquals executes the function if the key is set in the context
// as a given type and the value is equal to the expected value.
func IfContextValueEquals[T comparable](ctx context.Context, key ContextKey, expected T, fn func()) {
	IfContextSetWith(ctx, key, func(value T) {
		if value == expected {
			fn()
		}
	})
}
