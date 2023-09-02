package broker

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/log"
)

// Controller represents the functions that should be implemented to connect
// the broker to the application or the client
type Controller interface {
	// SetLogger set a logger that will log operations on broker controller
	SetLogger(logger log.Interface)

	// Publish a message to the broker
	Publish(ctx context.Context, channel string, mw Message) error

	// Subscribe to messages from the broker
	Subscribe(ctx context.Context, channel string) (msgs chan Message, stop chan interface{}, err error)

	// SetQueueName sets the name of the queue that will be used by the broker
	SetQueueName(name string)
}
