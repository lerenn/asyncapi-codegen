package brokers

import (
	"context"
	"errors"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/nats-io/nats.go"
)

// NATSController is the NATSController implementation for asyncapi-codegen
type NATSController struct {
	connection *nats.Conn
	logger     extensions.Logger
	queueGroup string
}

// NewNATSController creates a new NATS that fulfill the BrokerLinker interface
func NewNATSController(connection *nats.Conn) *NATSController {
	return &NATSController{
		connection: connection,
		queueGroup: DefaultQueueGroupID,
		logger:     extensions.DummyLogger{},
	}
}

// SetQueueGroup sets a custom queue group name for channel subscription
//
// It can be used for multiple applications listening one the same channel but
// wants to listen on different queues.
func (c *NATSController) SetQueueGroup(name string) {
	c.queueGroup = name
}

// SetLogger set a custom logger that will log operations on broker controller
func (c *NATSController) SetLogger(logger extensions.Logger) {
	c.logger = logger
}

// Publish a message to the broker
func (c *NATSController) Publish(_ context.Context, channel string, bm extensions.BrokerMessage) error {
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

// Subscribe to messages from the broker
func (c *NATSController) Subscribe(ctx context.Context, channel string) (msgs chan extensions.BrokerMessage, stop chan interface{}, err error) {
	// Subscribe to channel
	natsMsgs := make(chan *nats.Msg, 64)
	sub, err := c.connection.QueueSubscribeSyncWithChan(channel, c.queueGroup, natsMsgs)
	if err != nil {
		return nil, nil, err
	}

	// Handle events
	msgs = make(chan extensions.BrokerMessage, 64)
	stop = make(chan interface{}, 1)
	go func() {
		for {
			select {
			// Handle new message
			case msg := <-natsMsgs:
				// Get headers
				headers := make(map[string][]byte, len(msg.Header))
				for k, v := range msg.Header {
					if len(v) > 0 {
						headers[k] = []byte(v[0])
					}
				}

				// Create message
				msgs <- extensions.BrokerMessage{
					Headers: headers,
					Payload: msg.Data,
				}
			// Handle closure request from function caller
			case <-stop:
				if err := sub.Unsubscribe(); err != nil && !errors.Is(err, nats.ErrConnectionClosed) && c.logger != nil {
					c.logger.Error(ctx, err.Error())
				}

				if err := sub.Drain(); err != nil && !errors.Is(err, nats.ErrConnectionClosed) && c.logger != nil {
					c.logger.Error(ctx, err.Error())
				}

				close(msgs)
			}
		}
	}()

	return msgs, stop, nil
}
