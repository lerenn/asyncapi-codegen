package extensions

import (
	"context"
)

// BrokerMessage is a wrapper that will contain all information regarding a message
type BrokerMessage struct {
	Headers map[string][]byte
	Payload []byte
}

// BrokerController represents the functions that should be implemented to connect
// the broker to the application or the client
type BrokerController interface {
	// SetLogger set a logger that will log operations on broker controller
	SetLogger(logger Logger)

	// Publish a message to the broker
	Publish(ctx context.Context, channel string, mw BrokerMessage) error

	// Subscribe to messages from the broker
	Subscribe(ctx context.Context, channel string) (msgs chan BrokerMessage, stop chan interface{}, err error)

	// SetQueueName sets the name of the queue that will be used by the broker
	SetQueueName(name string)
}
