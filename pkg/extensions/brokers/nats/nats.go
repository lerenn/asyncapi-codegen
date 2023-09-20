package nats

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers"
	"github.com/nats-io/nats.go"
)

// Controller is the Controller implementation for asyncapi-codegen.
type Controller struct {
	connection *nats.Conn
	logger     extensions.Logger
	queueGroup string
}

// ControllerOption is a function that can be used to configure a NATS controller
// Examples: WithQueueGroup(), WithLogger().
type ControllerOption func(controller *Controller)

// NewController creates a new NATS controller.
func NewController(url string, options ...ControllerOption) *Controller {
	// Connect to NATS
	nc, err := nats.Connect(url)
	if err != nil {
		panic(err)
	}

	// Creates default controller
	controller := &Controller{
		connection: nc,
		queueGroup: brokers.DefaultQueueGroupID,
		logger:     extensions.DummyLogger{},
	}

	// Execute options
	for _, option := range options {
		option(controller)
	}

	return controller
}

// WithQueueGroup set a custom queue group for channel subscription.
func WithQueueGroup(name string) ControllerOption {
	return func(controller *Controller) {
		controller.queueGroup = name
	}
}

// WithLogger set a custom logger that will log operations on broker controller.
func WithLogger(logger extensions.Logger) ControllerOption {
	return func(controller *Controller) {
		controller.logger = logger
	}
}

// Publish a message to the broker.
func (c *Controller) Publish(_ context.Context, channel string, bm extensions.BrokerMessage) error {
	msg := nats.NewMsg(channel)

	// Set message headers and content
	for k, v := range bm.Headers {
		msg.Header.Set(k, string(v))
	}
	msg.Data = bm.Payload

	// Publish message
	if err := c.connection.PublishMsg(msg); err != nil {
		return err
	}

	// Flush the queue
	return c.connection.Flush()
}

// Subscribe to messages from the broker.
func (c *Controller) Subscribe(ctx context.Context, channel string) (
	messages chan extensions.BrokerMessage,
	cancel chan any,
	err error,
) {
	// Initialize channels
	messages = make(chan extensions.BrokerMessage, brokers.BrokerMessagesQueueSize)
	cancel = make(chan any, 1)

	// Subscribe on subject
	sub, err := c.connection.QueueSubscribe(channel, c.queueGroup, messagesHandler(messages))
	if err != nil {
		return nil, nil, err
	}

	go func() {
		// Wait for cancel request
		<-cancel

		// Drain the NATS subscription
		if err := sub.Drain(); err != nil {
			c.logger.Error(ctx, err.Error())
		}

		// Close messages in order to avoid new messages
		close(messages)

		// Close cancel to let listeners know that the cancellation is complete
		close(cancel)
	}()

	return messages, cancel, nil
}

func messagesHandler(messages chan extensions.BrokerMessage) nats.MsgHandler {
	return func(msg *nats.Msg) {
		// Get headers
		headers := make(map[string][]byte, len(msg.Header))
		for k, v := range msg.Header {
			if len(v) > 0 {
				headers[k] = []byte(v[0])
			}
		}

		// Create and transmit message to user
		messages <- extensions.BrokerMessage{
			Headers: headers,
			Payload: msg.Data,
		}
	}
}
