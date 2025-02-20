package rabbitmq

import (
	"context"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Check that it still fills the interface.
var _ extensions.BrokerController = (*Controller)(nil)

// ExchangeDeclare is a partial copy of amqp.ExchangeDeclare .
type ExchangeDeclare struct {
	Type       string
	Passive    bool
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Arguments  amqp.Table
}

// QueueDeclare is a partial copy of amqp.QueueDeclare .
type QueueDeclare struct {
	Durable    bool
	Exclusive  bool
	AutoDelete bool
	NoWait     bool
	Arguments  amqp.Table
}

// Controller is the Controller implementation for asyncapi-codegen.
type Controller struct {
	url             string
	connection      *amqp.Connection
	logger          extensions.Logger
	queueGroup      string
	exchangeOptions ExchangeDeclare
	queueOptions    QueueDeclare
}

// ControllerOption is a function that can be used to configure a RabbitMQ controller.
// Examples: WithQueueGroup(), WithLogger().
type ControllerOption func(controller *Controller) error

// NewController creates a new RabbitMQ controller.
func NewController(url string, options ...ControllerOption) (*Controller, error) {
	// Create default controller
	controller := &Controller{
		url:        url,
		queueGroup: brokers.DefaultQueueGroupID,
		logger:     extensions.DummyLogger{},
	}

	// Execute options
	for _, option := range options {
		if err := option(controller); err != nil {
			return nil, fmt.Errorf("could not apply option to controller: %w", err)
		}
	}

	// If connection not already created with WithConnectionOpts, connect to RabbitMQ
	if controller.connection == nil {
		conn, err := amqp.Dial(url)
		if err != nil {
			return nil, fmt.Errorf("could not connect to RabbitMQ: %w", err)
		}
		controller.connection = conn
	}

	return controller, nil
}

// WithLogger sets a custom logger that will log operations on the broker controller.
func WithLogger(logger extensions.Logger) ControllerOption {
	return func(controller *Controller) error {
		controller.logger = logger
		return nil
	}
}

// WithConnectionOpts sets the RabbitMQ.Config to connect to RabbitMQ.
func WithConnectionOpts(config amqp.Config) ControllerOption {
	return func(controller *Controller) error {
		conn, err := amqp.DialConfig(controller.url, config)
		if err != nil {
			return fmt.Errorf("could not connect to RabbitMQ: %w", err)
		}
		controller.connection = conn
		return nil
	}
}

// WithQueueGroup sets the queue group to use for the broker controller.
func WithQueueGroup(queueGroup string) ControllerOption {
	return func(controller *Controller) error {
		controller.queueGroup = queueGroup
		return nil
	}
}

// WithQueueOptions sets the queue options to use for the broker controller.
func WithQueueOptions(options QueueDeclare) ControllerOption {
	return func(controller *Controller) error {
		controller.queueOptions = options
		return nil
	}
}

func mergeQueueOptions(defaultOptions, options QueueDeclare) QueueDeclare {
	if options.Arguments == nil {
		options.Arguments = defaultOptions.Arguments
	}
	return options
}

// WithExchangeOptions sets the exchange options to use for the broker controller.
func WithExchangeOptions(options ExchangeDeclare) ControllerOption {
	return func(controller *Controller) error {
		controller.exchangeOptions = options
		return nil
	}
}

func isValidExchangeType(exchangeType string) bool {
	switch exchangeType {
	case "direct", "fanout", "topic", "headers":
		return true
	default:
		return false
	}
}

func mergeExchangeOptions(defaultOptions, options ExchangeDeclare) ExchangeDeclare {
	if !isValidExchangeType(options.Type) {
		options.Type = defaultOptions.Type
	}
	if options.Arguments == nil {
		options.Arguments = defaultOptions.Arguments
	}
	return options
}

// Publish a message to the broker.
func (c *Controller) Publish(_ context.Context, queueName string, bm extensions.BrokerMessage) error {
	channel, err := c.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	c.mergeDefaultOptions()

	err = channel.ExchangeDeclare(
		c.queueGroup,
		c.exchangeOptions.Type,
		c.exchangeOptions.Durable,
		c.exchangeOptions.AutoDelete,
		c.exchangeOptions.Internal,
		c.exchangeOptions.NoWait,
		c.exchangeOptions.Arguments,
	)
	if err != nil {
		return fmt.Errorf("could not declare exchange: %w", err)
	}

	_, err = channel.QueueDeclare(
		queueName,
		c.queueOptions.Durable,
		c.queueOptions.AutoDelete,
		c.queueOptions.Exclusive,
		c.queueOptions.NoWait,
		c.queueOptions.Arguments,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Convert headers
	headers := amqp.Table{}
	for k, v := range bm.Headers {
		headers[k] = v
	}

	// Publish message
	err = channel.Publish(
		c.queueGroup, // exchange
		queueName,    // routing key (queue name)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Body:            bm.Payload,
			Headers:         headers,
			ContentType:     "application/octet-stream",
			ContentEncoding: "binary",
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) mergeDefaultOptions() {
	defaultExchangeOptions := ExchangeDeclare{
		Type:       "direct",
		Durable:    false,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Arguments:  amqp.Table{},
	}

	c.exchangeOptions = mergeExchangeOptions(defaultExchangeOptions, c.exchangeOptions)

	defaultQueueOptions := QueueDeclare{
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		NoWait:     false,
		Arguments:  amqp.Table{},
	}

	c.queueOptions = mergeQueueOptions(defaultQueueOptions, c.queueOptions)
}

// Subscribe to messages from the broker.
//
//nolint:funlen
func (c *Controller) Subscribe(ctx context.Context, queueName string) (
	extensions.BrokerChannelSubscription, error) {
	// Create a new subscription
	sub := extensions.NewBrokerChannelSubscription(
		make(chan extensions.AcknowledgeableBrokerMessage, brokers.BrokerMessagesQueueSize),
		make(chan any, 1),
	)

	// Create a new channel
	channel, err := c.connection.Channel()
	if err != nil {
		return extensions.BrokerChannelSubscription{}, err
	}

	// Ensure the queue exists
	_, err = channel.QueueDeclare(
		queueName,
		false, // durable
		false, // auto-delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return extensions.BrokerChannelSubscription{}, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Start consuming
	msgs, err := channel.Consume(
		queueName,
		"",    // consumer tag
		false, // auto-ack
		false, // exclusive
		false, // no-local (deprecated)
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return extensions.BrokerChannelSubscription{}, fmt.Errorf("failed to start consuming from queue: %w", err)
	}

	// Wait for cancellation and clean up
	sub.WaitForCancellationAsync(func() {
		if err := channel.Cancel("", false); err != nil {
			c.logger.Error(ctx, fmt.Sprintf("failed to cancel consumer: %v", err))
		}
		channel.Close()
	})

	// Start a goroutine to receive messages and pass them to sub
	go func() {
		// No need to defer channel.Close() here as it will be closed in the cancellation handler
		for delivery := range msgs {
			// Get headers
			headers := make(map[string][]byte)
			for key, value := range delivery.Headers {
				switch v := value.(type) {
				case []byte:
					headers[key] = v
				case string:
					headers[key] = []byte(v)
				default:
					headers[key] = []byte(fmt.Sprintf("%v", v))
				}
			}

			// Create and transmit message to user
			sub.TransmitReceivedMessage(extensions.NewAcknowledgeableBrokerMessage(
				extensions.BrokerMessage{
					Headers: headers,
					Payload: delivery.Body,
				},
				&AcknowledgementHandler{
					Delivery: &delivery,
				},
			))
		}
	}()

	return sub, nil
}

// Close closes everything related to the broker.
func (c *Controller) Close() {
	if c.connection != nil {
		c.connection.Close()
	}
}

var _ extensions.BrokerAcknowledgment = (*AcknowledgementHandler)(nil)

// AcknowledgementHandler handles message acknowledgements.
type AcknowledgementHandler struct {
	Delivery *amqp.Delivery
}

// AckMessage acknowledges the message.
func (h *AcknowledgementHandler) AckMessage() {
	_ = h.Delivery.Ack(false)
}

// NakMessage negatively acknowledges the message.
func (h *AcknowledgementHandler) NakMessage() {
	_ = h.Delivery.Nack(false, false)
}
