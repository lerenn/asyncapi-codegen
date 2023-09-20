package extensions

import (
	"context"
)

// BrokerMessage is a wrapper that will contain all information regarding a message.
type BrokerMessage struct {
	Headers map[string][]byte
	Payload []byte
}

// IsUninitialized check if the BrokerMessage is at zero value, i.e. the
// uninitialized structure. It can be used to check that a channel is closed.
func (bm BrokerMessage) IsUninitialized() bool {
	return bm.Headers == nil && bm.Payload == nil
}

// BrokerController represents the functions that should be implemented to connect
// the broker to the application or the user.
type BrokerController interface {
	// Publish a message to the broker
	Publish(ctx context.Context, channel string, mw BrokerMessage) error

	// Subscribe to messages from the broker
	Subscribe(ctx context.Context, channel string) (messages chan BrokerMessage, cancel chan any, err error)
}
