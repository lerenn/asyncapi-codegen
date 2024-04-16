package kafka

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
)

// Check that it still fills the interface.
var _ extensions.BrokerController = (*Controller)(nil)

// Controller is the Kafka implementation for asyncapi-codegen.
type Controller struct {
	hosts []string
	// Reception only
	groupID string

	dialer *kafka.Dialer

	partition  int
	maxBytes   int
	autoCommit bool

	connectionTest bool

	logger extensions.Logger
}

// MessagesHandler is a function that can be used to process messages from the broker.
type MessagesHandler func(
	ctx context.Context,
	r *kafka.Reader,
	sub extensions.BrokerChannelSubscription,
)

// ControllerOption is a function that can be used to configure a Kafka controller
// Examples: WithGroupID(), WithPartition(), WithMaxBytes(), WithLogger().
type ControllerOption func(controller *Controller)

// NewController creates a new KafkaController that fulfill the BrokerLinker interface.
func NewController(hosts []string, options ...ControllerOption) (*Controller, error) {
	// Create default controller
	controller := &Controller{
		hosts:          hosts,
		logger:         extensions.DummyLogger{},
		groupID:        brokers.DefaultQueueGroupID,
		dialer:         kafka.DefaultDialer,
		partition:      0,
		maxBytes:       10e6, // 10MB
		autoCommit:     true,
		connectionTest: true,
	}

	// Execute options
	for _, option := range options {
		option(controller)
	}

	// kafka has no ping or something like this to test if dialer can create a successful connection to kafka
	// so if connectionTest is enabled create a connection and try to list brokers from kafka and validate
	// we can make a connection to kafka
	if controller.connectionTest {
		conn, err := controller.dialer.Dial("tcp", controller.hosts[0])
		if err != nil {
			return nil, fmt.Errorf("failed to create dialer for kafka: %w", err)
		}

		_, err = conn.Brokers()
		if err != nil {
			return nil, fmt.Errorf("failed to connect to kafka: %w", err)
		}
		conn.Close()
	}

	return controller, nil
}

// WithGroupID set a custom group ID for channel subscription.
func WithGroupID(groupID string) ControllerOption {
	return func(controller *Controller) {
		controller.groupID = groupID
	}
}

// WithPartition set the partition to use for the topic.
func WithPartition(partition int) ControllerOption {
	return func(controller *Controller) {
		controller.partition = partition
	}
}

// WithMaxBytes set the maximum size of a message.
func WithMaxBytes(maxBytes int) ControllerOption {
	return func(controller *Controller) {
		controller.maxBytes = maxBytes
	}
}

// WithLogger set a custom logger that will log operations on broker controller.
func WithLogger(logger extensions.Logger) ControllerOption {
	return func(controller *Controller) {
		controller.logger = logger
	}
}

// WithAutoCommit set if a AutoCommitMessagesHandler or ManualCommitMessagesHandler
// should be used for processing the messages.
func WithAutoCommit(enabled bool) ControllerOption {
	return func(controller *Controller) {
		controller.autoCommit = enabled
	}
}

// WithTLS set the tls.Config that will be used for kafka.Dial, kafka.Reader and kafka.Writer.
func WithTLS(tls *tls.Config) ControllerOption {
	return func(controller *Controller) {
		controller.dialer.TLS = tls
	}
}

// WithSasl set the sasl.Mechanism that will be used for kafka.Dial, kafka.Reader and kafka.Writer.
func WithSasl(sasl sasl.Mechanism) ControllerOption {
	return func(controller *Controller) {
		controller.dialer.SASLMechanism = sasl
	}
}

// WithConnectionTest set the connectionTest feature toggle to configure if NewController
// should validate the connection on creation.
func WithConnectionTest(enabled bool) ControllerOption {
	return func(controller *Controller) {
		controller.connectionTest = enabled
	}
}

// Publish a message to the broker.
func (c *Controller) Publish(ctx context.Context, channel string, um extensions.BrokerMessage) error {
	// Create new writer
	w := kafka.Writer{
		Addr:     kafka.TCP(c.hosts...),
		Topic:    channel,
		Balancer: &kafka.LeastBytes{},
		Transport: &kafka.Transport{
			// reuse the optionally TLS and SASLMechanism from dialer provided by the user to pass it to the writer
			// it can be nil
			TLS:  c.dialer.TLS.Clone(),
			SASL: c.dialer.SASLMechanism,
		},
	}

	defer w.Close()

	// Create the message
	msg := kafka.Message{
		Headers: make([]kafka.Header, 0),
	}

	// Set message content and headers
	msg.Value = um.Payload
	for k, v := range um.Headers {
		msg.Headers = append(msg.Headers, kafka.Header{Key: k, Value: v})
	}

	for {
		// Publish message
		err := w.WriteMessages(ctx, msg)

		// If there is no error then return
		if err == nil {
			return nil
		}

		// Create topic if not exists, then it means that the topic is being
		// created, so let's retry
		if errors.Is(err, kafka.UnknownTopicOrPartition) {
			c.logger.Warning(ctx, fmt.Sprintf("Topic %s does not exists: request creation and retry", channel))
			if err := c.checkTopicExistOrCreateIt(ctx, channel); err != nil {
				return err
			}

			continue
		}

		// Unexpected error
		return err
	}
}

// Subscribe to messages from the broker.
func (c *Controller) Subscribe(ctx context.Context, channel string) (extensions.BrokerChannelSubscription, error) {
	// Check that topic exists before
	if err := c.checkTopicExistOrCreateIt(ctx, channel); err != nil {
		return extensions.BrokerChannelSubscription{}, err
	}

	// Create reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   c.hosts,
		Topic:     channel,
		Partition: c.partition,
		MaxBytes:  c.maxBytes,
		GroupID:   c.groupID,
		Dialer:    c.dialer,
	})

	// Create subscription
	sub := extensions.NewBrokerChannelSubscription(
		make(chan extensions.AcknowledgeableBrokerMessage, brokers.BrokerMessagesQueueSize),
		make(chan any, 1),
	)

	// Handle events
	if c.autoCommit {
		go autoCommitMessagesHandler(&c.logger)(ctx, r, sub)
	} else {
		go manualCommitMessagesHandler(&c.logger)(ctx, r, sub)
	}

	// Wait for cancellation and stop the kafka listener when it happens
	sub.WaitForCancellationAsync(func() {
		if err := r.Close(); err != nil {
			c.logger.Error(ctx, err.Error())
		}
	})

	return sub, nil
}

func (c *Controller) checkTopicExistOrCreateIt(ctx context.Context, topic string) error {
	// Get connection to first host
	conn, err := c.dialer.Dial("tcp", c.hosts[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	for i := 0; ; i++ {
		// Create topic
		topicConfigs := []kafka.TopicConfig{{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}}
		err = conn.CreateTopics(topicConfigs...)
		if err != nil {
			return err
		}

		// Read partitions
		partitions, err := conn.ReadPartitions()
		if err != nil {
			return err
		}

		// Get topic from partitions
		for _, p := range partitions {
			if topic == p.Topic {
				if i > 0 {
					c.logger.Warning(ctx, fmt.Sprintf("Topic %s has been created.", topic))
				}
				return nil
			}
		}

		c.logger.Warning(ctx, fmt.Sprintf("Topic %s doesn't exists yet, retrying (#%d)", topic, i))
	}
}

// autoCommitMessagesHandler provides a MessagesHandler with auto commit.
//
// Using auto commit could result in an offset being committed before the
// message is fully processed and handled but allow more throughput.
//
// Maybe consider to use the manualCommitMessagesHandler.
func autoCommitMessagesHandler(
	logger *extensions.Logger,
) func(ctx context.Context, r *kafka.Reader, sub extensions.BrokerChannelSubscription) {
	return func(ctx context.Context, r *kafka.Reader, sub extensions.BrokerChannelSubscription) {
		for {
			msg, err := r.ReadMessage(ctx)
			if err != nil {
				// If the error is not io.EOF, then it is a real error
				if !errors.Is(err, io.EOF) {
					(*logger).Warning(ctx, fmt.Sprintf("Error when reading message: %q", err.Error()))
				}

				return
			}

			// Get headers
			headers := make(map[string][]byte, len(msg.Headers))
			for _, header := range msg.Headers {
				headers[header.Key] = header.Value
			}

			// Send received message
			sub.TransmitReceivedMessage(extensions.NewAcknowledgeableBrokerMessage(
				extensions.BrokerMessage{
					Headers: headers,
					Payload: msg.Value,
				},
				BrokerAcknowledgment{NoopCommit}))
		}
	}
}

// manualCommitMessagesHandler provides a MessagesHandler with manual commit
// the message is committed by user via the AcknowledgementHandler.
func manualCommitMessagesHandler(
	logger *extensions.Logger,
) func(ctx context.Context, r *kafka.Reader, sub extensions.BrokerChannelSubscription) {
	return func(ctx context.Context, r *kafka.Reader, sub extensions.BrokerChannelSubscription) {
		for {
			msg, err := r.FetchMessage(ctx)
			if err != nil {
				// If the error is not io.EOF, then it is a real error
				if !errors.Is(err, io.EOF) {
					(*logger).Warning(ctx, fmt.Sprintf("Error when reading message: %q", err.Error()))
				}

				return
			}

			// Get headers
			headers := make(map[string][]byte, len(msg.Headers))
			for _, header := range msg.Headers {
				headers[header.Key] = header.Value
			}

			// Send received message
			sub.TransmitReceivedMessage(extensions.NewAcknowledgeableBrokerMessage(
				extensions.BrokerMessage{
					Headers: headers,
					Payload: msg.Value,
				},
				BrokerAcknowledgment{doCommit: func() {
					if err := r.CommitMessages(ctx, msg); err != nil {
						(*logger).Error(ctx, fmt.Sprintf("error on committing message: %q", err.Error()))
					}
				}},
			))
		}
	}
}

var _ extensions.BrokerAcknowledgment = (*BrokerAcknowledgment)(nil)

// BrokerAcknowledgment for kafka broker.
// Naks are not supported on kafka side. Committing the message is the only way to handling the message.
// Proper errorhandling needs to be done by the subscriber.
type BrokerAcknowledgment struct {
	doCommit func()
}

// AckMessage acknowledges the message.
func (k BrokerAcknowledgment) AckMessage() {
	k.doCommit()
}

// NakMessage negatively acknowledges the message.
func (k BrokerAcknowledgment) NakMessage() {
	k.doCommit()
}

// NoopCommit is a no operation commit function.
func NoopCommit() {}
