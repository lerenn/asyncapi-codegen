package brokers

import (
	"context"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/segmentio/kafka-go"
)

// KafkaController is the Kafka implementation for asyncapi-codegen
type KafkaController struct {
	logger    extensions.Logger
	queueName string
	hosts     []string
	groupID   string
	partition int
	maxBytes  int
}

type KafkaControllerOption func(controller *KafkaController)

// NewKafkaController creates a new KafkaController that fulfill the BrokerLinker interface
func NewKafkaController(hosts []string, options ...KafkaControllerOption) *KafkaController {
	controller := &KafkaController{
		queueName: "asyncapi",
		hosts:     hosts,
		partition: 0,
		maxBytes:  10e6, // 10MB
	}
	for _, option := range options {
		option(controller)
	}
	return controller
}

func WithGroupID(groupID string) KafkaControllerOption {
	return func(controller *KafkaController) {
		controller.groupID = groupID
	}
}

func WithPartition(partition int) KafkaControllerOption {
	return func(controller *KafkaController) {
		controller.partition = partition
	}
}

func WithMaxBytes(maxBytes int) KafkaControllerOption {
	return func(controller *KafkaController) {
		controller.maxBytes = maxBytes
	}
}

// SetQueueName sets a custom queue name for channel subscription
//
// It can be used for multiple applications listening one the same channel but
// wants to listen on different queues.
func (c *KafkaController) SetQueueName(name string) {
	c.queueName = name
}

// SetLogger set a custom logger that will log operations on broker controller
func (c *KafkaController) SetLogger(logger extensions.Logger) {
	c.logger = logger
}

// Publish a message to the broker
func (c *KafkaController) Publish(ctx context.Context, channel string, um extensions.BrokerMessage) error {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(c.hosts...),
		Topic:                  channel,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	msg := kafka.Message{
		Headers: make([]kafka.Header, 0),
	}

	// Set message content
	msg.Value = um.Payload
	if um.CorrelationID != nil {
		msg.Headers = append(msg.Headers, kafka.Header{
			Key:   correlationIDField,
			Value: []byte(*um.CorrelationID),
		})
	}

	// Publish message
	if err := w.WriteMessages(ctx, msg); err != nil {
		return err
	}

	return nil
}

// Subscribe to messages from the broker
func (c *KafkaController) Subscribe(ctx context.Context, channel string) (msgs chan extensions.BrokerMessage, stop chan interface{}, err error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   c.hosts,
		Topic:     channel,
		Partition: c.partition,
		MaxBytes:  c.maxBytes,
	})

	getHeaders := func(msg kafka.Message, key string) string {
		for _, header := range msg.Headers {
			if header.Key == key {
				return string(header.Value)
			}
		}
		return ""
	}

	// Handle events
	msgs = make(chan extensions.BrokerMessage, 64)
	stop = make(chan interface{}, 1)
	go func() {
		for {
			msg, err := r.ReadMessage(ctx)
			if err != nil {
				break
			}
			var correlationID *string

			// Add correlation ID if not empty
			str := getHeaders(msg, correlationIDField)
			if str != "" {
				correlationID = &str
			}

			// Create message
			msgs <- extensions.BrokerMessage{
				Payload:       msg.Value,
				CorrelationID: correlationID,
			}
		}
	}()

	go func() {
		// Handle closure request from function caller
		for range stop {
			fmt.Print("Stopping subscriber")
			if err := r.Close(); err != nil && c.logger != nil {
				c.logger.Error(ctx, err.Error())
			}
			close(msgs)
		}
	}()

	return msgs, stop, nil
}
