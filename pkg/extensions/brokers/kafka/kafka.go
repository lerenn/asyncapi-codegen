package kafka

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers"
	"github.com/segmentio/kafka-go"
)

// Controller is the Kafka implementation for asyncapi-codegen.
type Controller struct {
	logger    extensions.Logger
	groupID   string
	hosts     []string
	partition int
	maxBytes  int
}

// ControllerOption is a function that can be used to configure a Kafka controller
// Examples: WithGroupID(), WithPartition(), WithMaxBytes(), WithLogger().
type ControllerOption func(controller *Controller)

// NewController creates a new KafkaController that fulfill the BrokerLinker interface.
func NewController(hosts []string, options ...ControllerOption) *Controller {
	// Create default controller
	controller := &Controller{
		logger:    extensions.DummyLogger{},
		groupID:   brokers.DefaultQueueGroupID,
		hosts:     hosts,
		partition: 0,
		maxBytes:  10e6, // 10MB
	}

	// Execute options
	for _, option := range options {
		option(controller)
	}

	return controller
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

// Publish a message to the broker.
func (c *Controller) Publish(ctx context.Context, channel string, um extensions.BrokerMessage) error {
	// Create new writer
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  c.hosts,
		Topic:    channel,
		Balancer: &kafka.LeastBytes{},
	})

	// Allow topic creation
	w.AllowAutoTopicCreation = true

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
			c.logger.Warning(ctx, fmt.Sprintf("Topic %s does not exists: waiting for creation and retry", channel))
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Unexpected error
		return err
	}
}

// Subscribe to messages from the broker.
func (c *Controller) Subscribe(ctx context.Context, channel string) (msgs chan extensions.BrokerMessage, stop chan interface{}, err error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   c.hosts,
		Topic:     channel,
		Partition: c.partition,
		MaxBytes:  c.maxBytes,
		GroupID:   c.groupID,
	})

	// Handle events
	msgs = make(chan extensions.BrokerMessage, 64)
	stop = make(chan interface{}, 1)
	go func() {
		for {
			msg, err := r.ReadMessage(ctx)
			if err != nil {
				break
			}

			// Get headers
			headers := make(map[string][]byte, len(msg.Headers))
			for _, header := range msg.Headers {
				headers[header.Key] = header.Value
			}

			// Create message
			msgs <- extensions.BrokerMessage{
				Headers: headers,
				Payload: msg.Value,
			}
		}
	}()

	go func() {
		// Handle closure request from function caller
		for range stop {
			c.logger.Info(ctx, "Stopping subscriber")
			if err := r.Close(); err != nil && c.logger != nil {
				c.logger.Error(ctx, err.Error())
			}
			close(msgs)
		}
	}()

	return msgs, stop, nil
}
