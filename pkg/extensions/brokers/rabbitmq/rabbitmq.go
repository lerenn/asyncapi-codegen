package rabbitmq

import (
    "context"
    "fmt"
    "sync"

    "github.com/lerenn/asyncapi-codegen/pkg/extensions"
    "github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers"
    amqp "github.com/rabbitmq/amqp091-go"
)

// Check interface implementation at compile time
var _ extensions.BrokerController = (*Controller)(nil)

// ExchangeDeclare represents RabbitMQ exchange configuration
type ExchangeDeclare struct {
    Type       string     // Exchange type (direct, fanout, topic, headers)
    Passive    bool       // If true, won't declare exchange, just check if exists
    Durable    bool       // Survives broker restart
    AutoDelete bool       // Deleted when last binding is removed
    Internal   bool       // If true, exchange cannot be directly published to by clients
    NoWait     bool       // If true, doesn't wait for server confirmation
    Arguments  amqp.Table // Additional exchange arguments
}

// QueueDeclare represents RabbitMQ queue configuration
type QueueDeclare struct {
    Durable    bool       // Survives broker restart
    Exclusive  bool       // Restricted to this connection
    AutoDelete bool       // Deleted when last consumer unsubscribes
    NoWait     bool       // If true, doesn't wait for server confirmation
    Arguments  amqp.Table // Additional queue arguments
}

// Controller manages RabbitMQ connections and operations
type Controller struct {
    url             string
    connection      *amqp.Connection
    logger          extensions.Logger
    queueGroup      string
    exchangeOptions ExchangeDeclare
    queueOptions    QueueDeclare
    mu              sync.Mutex // Protects connection state
    closed          bool
}

// ControllerOption configures the Controller during creation
type ControllerOption func(*Controller) error

// Default configuration constants
const (
    DefaultExchangeType = "direct"
    DefaultQueueGroup   = brokers.DefaultQueueGroupID
)

// NewController creates and initializes a new RabbitMQ controller
func NewController(url string, options ...ControllerOption) (*Controller, error) {
    c := &Controller{
        url:        url,
        queueGroup: DefaultQueueGroup,
        logger:     extensions.DummyLogger{},
        exchangeOptions: ExchangeDeclare{
            Type:      DefaultExchangeType,
            Arguments: make(amqp.Table),
        },
        queueOptions: QueueDeclare{
            Arguments: make(amqp.Table),
        },
    }

    for _, opt := range options {
        if err := opt(c); err != nil {
            return nil, fmt.Errorf("failed to apply option: %w", err)
        }
    }

    if c.connection == nil {
        if err := c.connect(); err != nil {
            return nil, fmt.Errorf("failed to establish initial connection: %w", err)
        }
    }

    return c, nil
}

// connect establishes a connection to RabbitMQ
func (c *Controller) connect() error {
    conn, err := amqp.Dial(c.url)
    if err != nil {
        return err
    }
    c.connection = conn
    return nil
}

// Controller options
func WithLogger(logger extensions.Logger) ControllerOption {
    return func(c *Controller) error {
        c.logger = logger
        return nil
    }
}

func WithConnectionOpts(config amqp.Config) ControllerOption {
    return func(c *Controller) error {
        conn, err := amqp.DialConfig(c.url, config)
        if err != nil {
            return fmt.Errorf("failed to connect with config: %w", err)
        }
        c.connection = conn
        return nil
    }
}

func WithQueueGroup(queueGroup string) ControllerOption {
    return func(c *Controller) error {
        if queueGroup == "" {
            return fmt.Errorf("queue group cannot be empty")
        }
        c.queueGroup = queueGroup
        return nil
    }
}

func WithQueueOptions(options QueueDeclare) ControllerOption {
    return func(c *Controller) error {
        c.queueOptions = options
        return nil
    }
}

func WithExchangeOptions(options ExchangeDeclare) ControllerOption {
    return func(c *Controller) error {
        if !isValidExchangeType(options.Type) {
            return fmt.Errorf("invalid exchange type: %s", options.Type)
        }
        c.exchangeOptions = options
        return nil
    }
}

// isValidExchangeType validates exchange type
func isValidExchangeType(exchangeType string) bool {
    switch exchangeType {
    case "direct", "fanout", "topic", "headers":
        return true
    default:
        return false
    }
}

// Publish sends a message to the specified queue
func (c *Controller) Publish(ctx context.Context, queueName string, bm extensions.BrokerMessage) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    if c.closed {
        return fmt.Errorf("controller is closed")
    }

    ch, err := c.connection.Channel()
    if err != nil {
        return fmt.Errorf("failed to open channel: %w", err)
    }
    defer ch.Close()

    if err := c.declareExchange(ch); err != nil {
        return err
    }

    if err := c.declareQueue(ch, queueName); err != nil {
        return err
    }

    return c.publishMessage(ch, queueName, bm)
}

func (c *Controller) declareExchange(ch *amqp.Channel) error {
    return ch.ExchangeDeclare(
        c.queueGroup,
        c.exchangeOptions.Type,
        c.exchangeOptions.Durable,
        c.exchangeOptions.AutoDelete,
        c.exchangeOptions.Internal,
        c.exchangeOptions.NoWait,
        c.exchangeOptions.Arguments,
    )
}

func (c *Controller) declareQueue(ch *amqp.Channel, queueName string) error {
    _, err := ch.QueueDeclare(
        queueName,
        c.queueOptions.Durable,
        c.queueOptions.AutoDelete,
        c.queueOptions.Exclusive,
        c.queueOptions.NoWait,
        c.queueOptions.Arguments,
    )
    return err
}

func (c *Controller) publishMessage(ch *amqp.Channel, queueName string, bm extensions.BrokerMessage) error {
    headers := amqp.Table{}
    for k, v := range bm.Headers {
        headers[k] = v
    }

    return ch.Publish(
        c.queueGroup,
        queueName,
        false,
        false,
        amqp.Publishing{
            Body:            bm.Payload,
            Headers:         headers,
            ContentType:     "application/octet-stream",
            ContentEncoding: "binary",
            Timestamp:       time.Now(),
        },
    )
}

// Subscribe creates a subscription to the specified queue
func (c *Controller) Subscribe(ctx context.Context, queueName string) (extensions.BrokerChannelSubscription, error) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if c.closed {
        return extensions.BrokerChannelSubscription{}, fmt.Errorf("controller is closed")
    }

    ch, err := c.connection.Channel()
    if err != nil {
        return extensions.BrokerChannelSubscription{}, fmt.Errorf("failed to open channel: %w", err)
    }

    if err := c.declareQueue(ch, queueName); err != nil {
        ch.Close()
        return extensions.BrokerChannelSubscription{}, err
    }

    return c.setupConsumer(ctx, ch, queueName)
}

func (c *Controller) setupConsumer(ctx context.Context, ch *amqp.Channel, queueName string) (extensions.BrokerChannelSubscription, error) {
    sub := extensions.NewBrokerChannelSubscription(
        make(chan extensions.AcknowledgeableBrokerMessage, brokers.BrokerMessagesQueueSize),
        make(chan any, 1),
    )

    msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
    if err != nil {
        ch.Close()
        return extensions.BrokerChannelSubscription{}, fmt.Errorf("failed to start consumer: %w", err)
    }

    go c.handleMessages(ctx, ch, sub, msgs)
    return sub, nil
}

func (c *Controller) handleMessages(ctx context.Context, ch *amqp.Channel, sub extensions.BrokerChannelSubscription, msgs <-chan amqp.Delivery) {
    defer ch.Close()
    for {
        select {
        case <-ctx.Done():
            return
        case <-sub.StopChan():
            return
        case d, ok := <-msgs:
            if !ok {
                return
            }
            sub.TransmitReceivedMessage(extensions.NewAcknowledgeableBrokerMessage(
                extensions.BrokerMessage{
                    Headers: convertHeaders(d.Headers),
                    Payload: d.Body,
                },
                &AcknowledgementHandler{Delivery: &d},
            ))
        }
    }
}

func convertHeaders(headers amqp.Table) map[string][]byte {
    result := make(map[string][]byte)
    for k, v := range headers {
        switch val := v.(type) {
        case []byte:
            result[k] = val
        case string:
            result[k] = []byte(val)
        default:
            result[k] = []byte(fmt.Sprintf("%v", v))
        }
    }
    return result
}

// Close cleanly shuts down the controller
func (c *Controller) Close() {
    c.mu.Lock()
    defer c.mu.Unlock()

    if c.closed {
        return
    }

    if c.connection != nil {
        if err := c.connection.Close(); err != nil {
            c.logger.Error(context.Background(), fmt.Sprintf("failed to close connection: %v", err))
        }
    }
    c.closed = true
}

// AcknowledgementHandler implements message acknowledgment
type AcknowledgementHandler struct {
    Delivery *amqp.Delivery
}

func (h *AcknowledgementHandler) AckMessage() {
    if h.Delivery != nil {
        _ = h.Delivery.Ack(false)
    }
}

func (h *AcknowledgementHandler) NakMessage() {
    if h.Delivery != nil {
        _ = h.Delivery.Nack(false, true) // Requeue the message
    }
}
