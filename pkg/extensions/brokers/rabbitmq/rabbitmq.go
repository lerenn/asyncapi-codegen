package rabbitmq

import (
	"context"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Compile-time check to ensure that Controller implements the extensions.BrokerController interface.
var _ extensions.BrokerController = (*Controller)(nil)

// ExchangeDeclare is a subset of amqp.ExchangeDeclare, used for configuring exchange options.
type ExchangeDeclare struct {
	Type       string
	Passive    bool
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Arguments  amqp.Table
}

// QueueDeclare is a subset of amqp.QueueDeclare, used for configuring queue options.
type QueueDeclare struct {
	Durable    bool
	Exclusive  bool
	AutoDelete bool
	NoWait     bool
	Arguments  amqp.Table
}

// Controller implements the extensions.BrokerController interface for RabbitMQ.
type Controller struct {
	url        string
	connection *amqp.Connection
	logger     extensions.Logger
	queueGroup string
	// Options are now pointers so we can check if they are set (nil or not nil)
	exchangeOptions *ExchangeDeclare
	queueOptions    *QueueDeclare
}

// ControllerOption is a functional option type to configure the Controller.
type ControllerOption func(controller *Controller) error

// NewController creates a new RabbitMQ controller.
func NewController(url string, options ...ControllerOption) (*Controller, error) {
	// Default values
	controller := &Controller{
		url:        url,
		queueGroup: brokers.DefaultQueueGroupID,
		logger:     extensions.DummyLogger{},
	}

	// Apply options
	for _, option := range options {
		if err := option(controller); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	// Establish connection if not provided via WithConnectionOpts
	if controller.connection == nil {
		conn, err := amqp.Dial(url)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
		}
		controller.connection = conn
	}

	return controller, nil
}

// WithLogger sets a custom logger.
func WithLogger(logger extensions.Logger) ControllerOption {
	return func(controller *Controller) error {
		controller.logger = logger
		return nil
	}
}

// WithConnectionOpts uses the provided amqp.Config for the RabbitMQ connection.
func WithConnectionOpts(config amqp.Config) ControllerOption {
	return func(controller *Controller) error {
		conn, err := amqp.DialConfig(controller.url, config)
		if err != nil {
			return fmt.Errorf("failed to connect to RabbitMQ with custom config: %w", err)
		}
		controller.connection = conn
		return nil
	}
}

// WithQueueGroup sets the queue group (exchange name in RabbitMQ) for the controller.
func WithQueueGroup(queueGroup string) ControllerOption {
	return func(controller *Controller) error {
		controller.queueGroup = queueGroup
		return nil
	}
}

// WithQueueOptions sets the queue options.
func WithQueueOptions(options QueueDeclare) ControllerOption {
	return func(controller *Controller) error {
		controller.queueOptions = &options // store a pointer
		return nil
	}
}

// mergeQueueOptions merges user-provided options with defaults, prioritizing user options.
func mergeQueueOptions(defaultOptions, userOptions *QueueDeclare) QueueDeclare {
	merged := defaultOptions // Start with default
	if userOptions != nil {
		merged = *userOptions // If options are provided, overwrite defaults
	}

	// Handle nested nil check for arguments
	if merged.Arguments == nil && defaultOptions.Arguments != nil {
		merged.Arguments = defaultOptions.Arguments
	}

	return merged
}

// WithExchangeOptions sets the exchange options.
func WithExchangeOptions(options ExchangeDeclare) ControllerOption {
	return func(controller *Controller) error {
		controller.exchangeOptions = &options // store a pointer
		return nil
	}
}

// isValidExchangeType checks if the provided exchange type is valid.
func isValidExchangeType(exchangeType string) bool {
	switch exchangeType {
	case "direct", "fanout", "topic", "headers":
		return true
	default:
		return false
	}
}

// mergeExchangeOptions merges user-provided exchange options with defaults.
func mergeExchangeOptions(defaultOptions, userOptions *ExchangeDeclare) ExchangeDeclare {
	merged := defaultOptions // Start with default
	if userOptions != nil {
		merged = *userOptions // If options provided, overwrite defaults
	}

	if !isValidExchangeType(merged.Type) {
		merged.Type = defaultOptions.Type // Use default if invalid
	}

	// Handle nested nil checks for arguments
	if merged.Arguments == nil && defaultOptions.Arguments != nil {
		merged.Arguments = defaultOptions.Arguments
	}

	return merged
}

// Publish publishes a message to the specified queue in RabbitMQ.
func (c *Controller) Publish(ctx context.Context, queueName string, bm extensions.BrokerMessage) error {
	channel, err := c.connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer channel.Close()

	// Merge options with defaults, doing it on each publish allows per-publish customization if needed.
	c.mergeDefaultOptions()

	// Declare Exchange
	if err := channel.ExchangeDeclare(
		c.queueGroup,             // name
		c.exchangeOptions.Type,   // type
		c.exchangeOptions.Durable,    // durable
		c.exchangeOptions.AutoDelete, // auto-deleted
		c.exchangeOptions.Internal,   // internal
		c.exchangeOptions.NoWait,     // no-wait
		c.exchangeOptions.Arguments,  // arguments
	); err != nil {
		return fmt.Errorf("failed to declare an exchange: %w", err)
	}

	// Declare Queue
	_, err = channel.QueueDeclare(
		queueName,                // name
		c.queueOptions.Durable,    // durable
		c.queueOptions.AutoDelete, // delete when unused
		c.queueOptions.Exclusive,  // exclusive
		c.queueOptions.NoWait,     // no-wait
		c.queueOptions.Arguments,  // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Convert headers to amqp.Table
	headers := amqp.Table{}
	for k, v := range bm.Headers {
		headers[k] = v
	}

	// Publish the message
	if err := channel.PublishWithContext(
		ctx,
		c.queueGroup, // exchange
		queueName,    // routing key (queue name)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Headers:         headers,
			ContentType:     "application/octet-stream", // Or determine dynamically
			ContentEncoding: "binary",                   // Or determine dynamically
			Body:            bm.Payload,
		},
	); err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	c.logger.Info(ctx, fmt.Sprintf("Published message to queue %s", queueName))
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

	c.exchangeOptions = &mergeExchangeOptions(defaultExchangeOptions, c.exchangeOptions)

	defaultQueueOptions := QueueDeclare{
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		NoWait:     false,
		Arguments:  amqp.Table{},
	}

	c.queueOptions = &mergeQueueOptions(defaultQueueOptions, c.queueOptions)
}

// Subscribe subscribes to messages from the specified queue in RabbitMQ.
func (c *Controller) Subscribe(ctx context.Context, queueName string) (extensions.BrokerChannelSubscription, error) {
	// Setup subscription channels
	sub := extensions.NewBrokerChannelSubscription(
		make(chan extensions.AcknowledgeableBrokerMessage, brokers.BrokerMessagesQueueSize),
		make(chan any, 1),
	)

	channel, err := c.connection.Channel()
	if err != nil {
		return extensions.BrokerChannelSubscription{}, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Merge options (for queue declaration, could be different from publish options)
	c.mergeDefaultOptions()

	// Declare Queue
	_, err = channel.QueueDeclare(
		queueName,                // name
		c.queueOptions.Durable,    // durable
		c.queueOptions.AutoDelete, // delete when unused
		c.queueOptions.Exclusive,  // exclusive
		c.queueOptions.NoWait,     // no-wait
		c.queueOptions.Arguments,  // arguments
	)
	if err != nil {
		return extensions.BrokerChannelSubscription{}, fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Start consuming messages
	msgs, err := channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack.  Set to false for manual acknowledgements.
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return extensions.BrokerChannelSubscription{}, fmt.Errorf("failed to register a consumer: %w", err)
	}

	// Asynchronously handle cancellation and cleanup
	sub.WaitForCancellationAsync(func() {
		if err := channel.Cancel("", false); err != nil {
			c.logger.Error(ctx, fmt.Sprintf("failed to cancel consumer: %v", err))
		}
		channel.Close()
		c.logger.Info(ctx, fmt.Sprintf("Unsubscribed from queue %s", queueName))
	})

	// Goroutine to process received messages
	go func() {
		for delivery := range msgs {
			// Convert amqp.Table headers to map[string][]byte
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

			// Create and transmit AcknowledgeableBrokerMessage
			sub.TransmitReceivedMessage(extensions.NewAcknowledgeableBrokerMessage(
				extensions.BrokerMessage{
					Headers: headers,
					Payload: delivery.Body,
				},
				&AcknowledgementHandler{
					Delivery: &delivery,
					ctx:      ctx,
					logger:   c.logger,
				},
			))
		}
		// The loop breaks when the msgs channel is closed (e.g., on connection loss or explicit close).
		// The subscription should also be stopped in that case.
		sub.Cancel()
	}()

	c.logger.Info(ctx, fmt.Sprintf("Subscribed to queue %s", queueName))
	return sub, nil
}

// Close closes the RabbitMQ connection.
func (c *Controller) Close() {
	if c.connection != nil {
		if err := c.connection.Close(); err != nil {
			// Use background context as the original context could be done
			c.logger.Error(context.Background(), fmt.Sprintf("error closing connection: %v", err))
		}
	}
	c.logger.Info(context.Background(), "Closed connection") // Using Background() since the Controller doesn't have a context
}

// Ensure AcknowledgementHandler implements the BrokerAcknowledgment interface.
var _ extensions.BrokerAcknowledgment = (*AcknowledgementHandler)(nil)

// AcknowledgementHandler handles message acknowledgements.
type AcknowledgementHandler struct {
	Delivery *amqp.Delivery
	ctx      context.Context // add context for logging
	logger   extensions.Logger
}

// AckMessage acknowledges the message.
func (h *AcknowledgementHandler) AckMessage() {
	if err := h.Delivery.Ack(false); err != nil {
		h.logger.Error(h.ctx, fmt.Sprintf("failed to ack message: %v", err))
	}
}

// NakMessage negatively acknowledges the message.
func (h *AcknowledgementHandler) NakMessage() {
	if err := h.Delivery.Nack(false, false); err != nil {
		h.logger.Error(h.ctx, fmt.Sprintf("failed to nack message: %v", err))
	}
}