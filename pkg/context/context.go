package context

import "context"

const Prefix = "aapi-"

// Key is the type of the keys used in the context
type Key string

const (
	// KeyIsModule is the name of the module this data is coming from.
	// When coming from a generated client, it is `asyncapi`
	KeyIsModule Key = Prefix + "module"
	// KeyIsProvider is the name of the provider this data is coming from.
	// When coming from generated code, it is `app`, `client` or `broker`
	KeyIsProvider Key = Prefix + "provider"
	// KeyIsChannel is the name of the channel this data is coming from.
	KeyIsChannel Key = Prefix + "channel"
	// KeyIsOperation is the name of the operation this data is coming from.
	// When coming from generated code, it is `subscribe`, `publish`, `wait-for`, etc
	KeyIsOperation Key = Prefix + "operation"
	// KeyIsMessage is the message that has been sent or received
	KeyIsMessage Key = Prefix + "message"
	// KeyIsCorrelationID is the correlation ID of the message
	KeyIsCorrelationID Key = Prefix + "correlationID"
	// KeyIsDirection is the direction of the message
	// It can be either "publication" or "reception"
	KeyIsDirection Key = Prefix + "direction"
)

// String returns the string representation of the key
func (k Key) String() string {
	return string(k)
}

// IfSetAsString executes the function if the key is set in the context as a string
func IfSet[T any](ctx context.Context, key Key, fn func(value T)) {
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

// IfEquals executes the function if the key is set in the context as a string
// and the value is equal to the expected value
func IfEquals[T comparable](ctx context.Context, key Key, expected T, fn func()) {
	IfSet(ctx, key, func(value T) {
		if value == expected {
			fn()
		}
	})
}
