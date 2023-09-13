package extensions

import (
	"context"
)

// BrokerMessage is a wrapper that will contain all information regarding a message.
type BrokerMessage struct {
	Headers map[string][]byte
	Payload []byte
}

// BrokerController represents the functions that should be implemented to connect
// the broker to the application or the user.
type BrokerController interface {
	// Publish a message to the broker
	Publish(ctx context.Context, channel string, mw BrokerMessage) error

	// Subscribe to messages from the broker
	Subscribe(ctx context.Context, channel string) (msgs chan BrokerMessage, stop chan interface{}, err error)
}
