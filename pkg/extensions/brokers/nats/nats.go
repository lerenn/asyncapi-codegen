package nats

import (
	"context"
	"errors"

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
	receivedBrokerMessages chan extensions.BrokerMessage,
	cancelSubscription chan interface{},
	err error,
) {
	// Subscribe to channel
	natsChan := make(chan *nats.Msg, 64)
	sub, err := c.connection.QueueSubscribeSyncWithChan(channel, c.queueGroup, natsChan)
	if err != nil {
		return nil, nil, err
	}

	// Handle events
	receivedBrokerMessages = make(chan extensions.BrokerMessage, 64)
	cancelSubscription = make(chan interface{}, 1)
	go func() {
		for stop := false; !stop; {
			select {
			// If its a new message, then handle it
			case msg := <-natsChan:
				transferMessageToBroker(receivedBrokerMessages, msg)
			// If its a closure request from function caller, then stop everything
			case <-cancelSubscription:
				stop = true
			}
		}

		// Unsubscribe from channel
		if err := sub.Unsubscribe(); err != nil && !errors.Is(err, nats.ErrConnectionClosed) && c.logger != nil {
			c.logger.Error(ctx, err.Error())
		}

		close(receivedBrokerMessages)
	}()

	return receivedBrokerMessages, cancelSubscription, nil
}

func transferMessageToBroker(receivedBrokerMessages chan extensions.BrokerMessage, msg *nats.Msg) {
	// Get headers
	headers := make(map[string][]byte, len(msg.Header))
	for k, v := range msg.Header {
		if len(v) > 0 {
			headers[k] = []byte(v[0])
		}
	}

	// Create and transmit message to user
	receivedBrokerMessages <- extensions.BrokerMessage{
		Headers: headers,
		Payload: msg.Data,
	}
}
