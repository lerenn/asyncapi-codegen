package natsjetstream

import (
	"context"
	"sync"
	"testing"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen // this is only for testing
func TestValidateAckMechanism(t *testing.T) {
	subj := "ValidateAckMechanism"

	broker, err := NewController(
		"nats://nats-jetstream:4222",
		WithStreamConfig(jetstream.StreamConfig{
			Name:     subj,
			Subjects: []string{subj},
		}),
		WithConsumerConfig(jetstream.ConsumerConfig{Name: "ValidateAckMechanism"}),
	)
	assert.NoError(t, err, "new controller should not return error")

	t.Run("validate ack is supported in NATS jetstream", func(t *testing.T) {
		wg := sync.WaitGroup{}
		stream, err := broker.jetStream.Stream(context.Background(), subj)
		assert.NoError(t, err, "stream should not return error")

		cons, err := stream.Consumer(context.Background(), broker.consumerName)
		assert.NoError(t, err, "consumer should not return error")

		cc, err := cons.Consume(func(msg jetstream.Msg) {
			defer wg.Done()
			err := msg.Ack()
			assert.NoError(t, err, "NATS jetstream should support acks")
		})
		assert.NoError(t, err, "consume should not return error")
		defer cc.Stop()

		wg.Add(1)
		_, err = broker.jetStream.Publish(context.Background(), subj, []byte("testmessage"))
		assert.NoError(t, err, "publish should not return error")

		wg.Wait()
	})

	t.Run("validate nak is supported in NATS jetstream", func(t *testing.T) {
		wg := sync.WaitGroup{}

		stream, err := broker.jetStream.Stream(context.Background(), subj)
		assert.NoError(t, err, "stream should not return error")

		cons, err := stream.Consumer(context.Background(), broker.consumerName)
		assert.NoError(t, err, "consumer should not return error")

		cc, err := cons.Consume(func(msg jetstream.Msg) {
			defer wg.Done()

			// use term instead of nak to tell the broker to not redeliver this message again - nak is then supported as well
			// otherwise the redelivery break the tests because it will be consumed again and wg is empty
			err := msg.Term()
			assert.NoError(t, err, "NATS jetstream should support naks")
		})
		assert.NoError(t, err, "consume should not return error")
		defer cc.Stop()

		wg.Add(1)
		_, err = broker.jetStream.Publish(context.Background(), subj, []byte("testmessage"))
		assert.NoError(t, err, "publish should not return error")

		wg.Wait()
	})
}
