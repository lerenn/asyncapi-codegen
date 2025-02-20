package rabbitmq

import (
	"context"
	"sync"
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	testutil "github.com/lerenn/asyncapi-codegen/pkg/utils/test"
	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen // this is only for testing
func TestValidateAckMechanism(t *testing.T) {
	// Establish a connection to the AMQP broker
	rmqb, err := NewController(
		testutil.BrokerAddress(testutil.BrokerAddressParams{
			Schema:         "amqp",
			DockerizedAddr: "rabbitmq",
			Port:           "5672",
		}),
	)
	assert.NoError(t, err, "new controller should not return error")
	defer rmqb.Close()

	t.Run("validate ack is supported in AMQP", func(t *testing.T) {
		queueName := "TestQueueAck"

		// Open a channel
		ch, err := rmqb.connection.Channel()
		assert.NoError(t, err, "should be able to open a channel")
		defer ch.Close()

		// Declare a queue
		_, err = ch.QueueDeclare(
			queueName,
			false, // durable
			true,  // auto-delete
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		assert.NoError(t, err, "should be able to declare queue")

		wg := sync.WaitGroup{}

		msgs, err := ch.Consume(
			queueName,
			"",    // consumer tag
			false, // auto-ack
			false, // exclusive
			false, // no-local (deprecated)
			false, // no-wait
			nil,   // arguments
		)
		assert.NoError(t, err, "should be able to start consuming messages")

		wg.Add(1)
		go func() {
			defer wg.Done()
			for d := range msgs {
				err := d.Ack(false)
				assert.NoError(t, err, "should be able to ack the message")
				break
			}
		}()

		err = ch.Publish(
			"",        // exchange
			queueName, // routing key (queue name)
			false,     // mandatory
			false,     // immediate
			amqp091.Publishing{
				ContentType: "text/plain",
				Body:        []byte("testmessage"),
			},
		)
		assert.NoError(t, err, "should be able to publish message")

		wg.Wait()
	})

	t.Run("validate nack is supported in AMQP", func(t *testing.T) {
		queueName := "TestQueueNack"

		// Open a channel
		ch, err := rmqb.connection.Channel()
		assert.NoError(t, err, "should be able to open a channel")
		defer ch.Close()

		// Declare a queue
		_, err = ch.QueueDeclare(
			queueName,
			false, // durable
			true,  // auto-delete
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		assert.NoError(t, err, "should be able to declare queue")

		wg := sync.WaitGroup{}

		msgs, err := ch.Consume(
			queueName,
			"",    // consumer tag
			false, // auto-ack
			false, // exclusive
			false, // no-local (deprecated)
			false, // no-wait
			nil,   // arguments
		)
		assert.NoError(t, err, "should be able to start consuming messages")

		wg.Add(1)
		go func() {
			defer wg.Done()
			for d := range msgs {
				err := d.Nack(false, false)
				assert.NoError(t, err, "should be able to nack the message")
				break
			}
		}()

		err = ch.Publish(
			"",        // exchange
			queueName, // routing key (queue name)
			false,     // mandatory
			false,     // immediate
			amqp091.Publishing{
				ContentType: "text/plain",
				Body:        []byte("testmessage"),
			},
		)
		assert.NoError(t, err, "should be able to publish message")

		wg.Wait()
	})
}

func TestRabbitMQController_WithQueueGroup(t *testing.T) {
	subj := "test-queue-group"
	controller, err := NewController(
		testutil.BrokerAddress(testutil.BrokerAddressParams{
			Schema:         "amqp",
			DockerizedAddr: "rabbitmq",
			Port:           "5672",
		}),
		WithQueueGroup(subj),
	)
	assert.NoError(t, err, "should be able to create RabbitMQ controller")
	defer controller.Close()
	ch, err := controller.connection.Channel()
	assert.NoError(t, err, "should be able to get channel")
	defer ch.Close()
}

func TestRabbitMQController_WithExchangeOptionsAndQueueOptions(t *testing.T) {
	controller, err := NewController(
		testutil.BrokerAddress(testutil.BrokerAddressParams{
			Schema:         "amqp",
			DockerizedAddr: "rabbitmq",
			Port:           "5672",
		}),
		WithQueueGroup("test-queue-group-for-all-options"),
		WithExchangeOptions(ExchangeDeclare{
			Type:       "fanout",
			Durable:    false,
			AutoDelete: true,
			Internal:   false,
			NoWait:     false,
			Arguments:  amqp091.Table{},
		}),
		WithQueueOptions(QueueDeclare{
			Durable:    false,
			AutoDelete: true,
			Exclusive:  false,
			NoWait:     false,
			Arguments:  amqp091.Table{},
		}),
	)
	assert.NoError(t, err, "should be able to create RabbitMQ controller")
	defer controller.Close()
	ch, err := controller.connection.Channel()
	assert.NoError(t, err, "should be able to get channel")
	defer ch.Close()

	_, err = controller.Subscribe(context.Background(), "test-queue2")
	assert.NoError(t, err, "should be able to subscribe to queue")

	err = controller.Publish(context.Background(), "test-queue-for-all-options", extensions.BrokerMessage{
		Payload: []byte("test-payload"),
	})
	assert.NoError(t, err, "should be able to publish to queue")
}

func TestValidExchangeType(t *testing.T) {
	assert.True(t, isValidExchangeType("direct"))
	assert.True(t, isValidExchangeType("fanout"))
	assert.True(t, isValidExchangeType("topic"))
	assert.True(t, isValidExchangeType("headers"))
	assert.False(t, isValidExchangeType("invalid"))
	assert.False(t, isValidExchangeType(""))
	assert.False(t, isValidExchangeType(" "))
	assert.False(t, isValidExchangeType("direct "))
}
