package rabbitmq

import (
	"sync"
	"testing"

	testutil "github.com/lerenn/asyncapi-codegen/pkg/utils/test"
	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen // this is only for testing
func TestValidateAckMechanism(t *testing.T) {
	// Establish a connection to the AMQP broker
	subj := "CoreRabbitmqValidateAckMechanism"
	rmqb, err := NewController(
		testutil.BrokerAddress(testutil.BrokerAddressParams{
			Schema:         "amqp",
			DockerizedAddr: "rabbitmq",
			Port:           "5672",
		}),
		WithQueueGroup(subj, "direct", false, true, false, false, nil))
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
	queueGroupID := "test-queue-group"
	controller, err := NewController(
		testutil.BrokerAddress(testutil.BrokerAddressParams{
			Schema:         "amqp",
			DockerizedAddr: "rabbitmq",
			Port:           "5672",
		}),
		WithQueueGroup(queueGroupID, "topic", false, true, false, false, nil),
	)
	assert.NoError(t, err, "should be able to create RabbitMQ controller")
	defer controller.Close()
	ch, err := controller.connection.Channel()
	assert.NoError(t, err, "should be able to get channel")
	defer ch.Close()
}
