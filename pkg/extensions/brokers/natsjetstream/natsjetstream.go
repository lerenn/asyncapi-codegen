package natsjetstream

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Check that it still fills the interface.
var _ extensions.BrokerController = (*Controller)(nil)

// Controller is the Controller implementation for asyncapi-codegen.
type Controller struct {
	natsConn       *nats.Conn
	jetStream      jetstream.JetStream
	logger         extensions.Logger
	streamName     string
	consumerName   string
	consumeContext jetstream.ConsumeContext
	channels       map[string]chan jetstream.Msg

	nakDelay time.Duration
}

// ControllerOption is a function that can be used to configure a NATS controller.
type ControllerOption func(controller *Controller) error

// WithLogger set a custom logger that will log operations on broker controller.
func WithLogger(logger extensions.Logger) ControllerOption {
	return func(controller *Controller) error {
		controller.logger = logger
		return nil
	}
}

// WithStreamConfig creates or updates a stream based on the given stream configuration.
func WithStreamConfig(config jetstream.StreamConfig) ControllerOption {
	return func(controller *Controller) error {
		if config.Name == "" {
			return fmt.Errorf("stream name is required")
		}
		controller.streamName = config.Name

		if _, err := controller.jetStream.CreateStream(context.Background(), config); err != nil {
			if !errors.Is(err, jetstream.ErrStreamNameAlreadyInUse) {
				return fmt.Errorf("could not create stream: %w", err)
			}
			if _, err = controller.jetStream.UpdateStream(context.Background(), config); err != nil {
				return fmt.Errorf("could not create or update stream: %w", err)
			}
		}

		return nil
	}
}

// WithStream uses the given stream name (the stream has to be created before initializing).
func WithStream(name string) ControllerOption {
	return func(controller *Controller) error {
		controller.streamName = name
		return nil
	}
}

// WithConsumerConfig creates or updates a consumer based on the given consumer configuration.
func WithConsumerConfig(config jetstream.ConsumerConfig) ControllerOption {
	return func(controller *Controller) error {
		if config.Name == "" {
			return fmt.Errorf("consumer name is required")
		}
		controller.consumerName = config.Name

		_, err := controller.jetStream.CreateOrUpdateConsumer(context.Background(), controller.streamName, config)
		if err != nil {
			return fmt.Errorf("could not create or update consumer: %w", err)
		}

		return nil
	}
}

// WithConsumer uses the given consumer name (the consumer has to be created before initializing).
func WithConsumer(name string) ControllerOption {
	return func(controller *Controller) error {
		controller.consumerName = name
		return nil
	}
}

// WithNakDelay set the delay when redeliver messages via nak
func WithNakDelay(duration time.Duration) ControllerOption {
	return func(controller *Controller) error {
		controller.nakDelay = duration
		return nil
	}
}

// NewController creates a new NATS JetStream controller.
func NewController(url string, options ...ControllerOption) (*Controller, error) {
	// Connect to NATS
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("could not connect to nats: %w", err)
	}

	// Create a JetStream management interface
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("could not connect to jetstream: %w", err)
	}

	// Creates default controller
	controller := &Controller{
		natsConn:       nc,
		jetStream:      js,
		logger:         extensions.DummyLogger{},
		channels:       make(map[string]chan jetstream.Msg),
		consumeContext: nil,
		nakDelay:       time.Second * 5,
	}

	// Execute options
	for _, option := range options {
		if err := option(controller); err != nil {
			return nil, fmt.Errorf("could not apply option to controller: %w", err)
		}
	}

	return controller, nil
}

// Publish a message to the broker.
func (c *Controller) Publish(ctx context.Context, channel string, bm extensions.BrokerMessage) error {
	msg := nats.NewMsg(channel)

	// Set message headers and content
	for k, v := range bm.Headers {
		msg.Header.Set(k, string(v))
	}
	msg.Data = bm.Payload

	// Publish message
	if _, err := c.jetStream.PublishMsg(ctx, msg); err != nil {
		return err
	}

	return nil
}

// Subscribe to messages from the broker.
func (c *Controller) Subscribe(ctx context.Context, channel string) (extensions.BrokerChannelSubscription, error) {
	// Create a new subscription
	sub := extensions.NewBrokerChannelSubscription(
		make(chan extensions.AcknowledgeableBrokerMessage, brokers.BrokerMessagesQueueSize),
		make(chan any, 1),
	)

	if c.channels[channel] == nil {
		c.channels[channel] = make(chan jetstream.Msg)
	}
	if err := c.ConsumeIfNeeded(ctx); err != nil {
		return extensions.BrokerChannelSubscription{}, err
	}

	go func() {
		for message := range c.channels[channel] {
			c.logger.Info(ctx, fmt.Sprintf("Received message for %s", channel), extensions.LogInfo{
				Key:   "message",
				Value: message,
			})
			c.HandleMessage(ctx, message, sub)
		}
	}()

	// Wait for cancellation and drain the NATS subscription
	sub.WaitForCancellationAsync(func() {
		close(c.channels[channel])
		delete(c.channels, channel)
		c.StopConsumeIfNeeded()
	})

	return sub, nil
}

// HandleMessage handles a message received from a stream.
func (c *Controller) HandleMessage(ctx context.Context, msg jetstream.Msg, sub extensions.BrokerChannelSubscription) {
	// Get headers
	headers := make(map[string][]byte, len(msg.Headers()))
	for k, v := range msg.Headers() {
		if len(v) > 0 {
			headers[k] = []byte(v[0])
		}
	}

	// Create and transmit message to user
	sub.TransmitReceivedMessage(extensions.NewAcknowledgeableBrokerMessage(
		extensions.BrokerMessage{
			Headers: headers,
			Payload: msg.Data(),
		},
		AcknowledgementHandler{
			doAck: func() {
				if err := msg.Ack(); err != nil {
					c.logger.Error(ctx, fmt.Sprintf("error on ack message: %q", err.Error()))
				}
			},
			doNak: func() {
				if err := msg.NakWithDelay(c.nakDelay); err != nil {
					c.logger.Error(ctx, fmt.Sprintf("error on nak message: %q", err.Error()))
				}
			},
		}))
}

// Close closes everything related to the broker.
func (c *Controller) Close() {
	c.natsConn.Close()
}

// ConsumeIfNeeded starts consuming messages if needed.
func (c *Controller) ConsumeIfNeeded(ctx context.Context) error {
	if c.consumeContext == nil {
		consumer, err := c.jetStream.Consumer(ctx, c.streamName, c.consumerName)
		if err != nil {
			return err
		}
		consumeContext, err := consumer.Consume(c.ConsumeMessage(ctx))
		if err != nil {
			return err
		}
		c.consumeContext = consumeContext
	}

	return nil
}

// StopConsumeIfNeeded stops consuming messages if needed (there is no other subscription).
func (c *Controller) StopConsumeIfNeeded() {
	if len(c.channels) == 0 && c.consumeContext != nil {
		c.consumeContext.Stop()
		c.consumeContext = nil
	}
}

// ConsumeMessage writes the message to the correct channel of the subject or in case
// there is no subscription the message will be acknowledged.
func (c *Controller) ConsumeMessage(ctx context.Context) jetstream.MessageHandler {
	return func(msg jetstream.Msg) {
		if c.channels[msg.Subject()] == nil {
			c.logger.Warning(
				ctx,
				fmt.Sprintf(
					"Received message for not subscribed channel '%s'. Message will be ack'd.",
					msg.Subject(),
				),
				extensions.LogInfo{Key: "msg.subject", Value: msg.Subject()},
				extensions.LogInfo{Key: "msg.headers", Value: msg.Headers()},
				extensions.LogInfo{Key: "msg.data", Value: msg.Data()},
			)
			_ = msg.Ack()
		}
		c.channels[msg.Subject()] <- msg
	}
}

var _ extensions.BrokerAcknowledgment = (*AcknowledgementHandler)(nil)

// AcknowledgementHandler for nats jetstream broker.
type AcknowledgementHandler struct {
	doAck func()
	doNak func()
}

func (k AcknowledgementHandler) AckMessage() {
	k.doAck()
}

func (k AcknowledgementHandler) NakMessage() {
	k.doNak()
}
