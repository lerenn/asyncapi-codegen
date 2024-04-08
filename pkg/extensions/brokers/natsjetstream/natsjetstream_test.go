package natsjetstream

import (
	"context"
	"crypto/tls"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

//nolint:funlen // this is only for testing
func TestValidateAckMechanism(t *testing.T) {
	subj := "NatsJetstreamValidateAckMechanism"

	broker, err := NewController(
		"nats://nats-jetstream:4222",
		WithStreamConfig(jetstream.StreamConfig{
			Name:     subj,
			Subjects: []string{subj},
		}),
		WithConsumerConfig(jetstream.ConsumerConfig{Name: "natsJetstreamValidateAckMechanism"}),
	)
	assert.NoError(t, err, "new controller should not return error")
	//defer broker.Close()

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

func TestSecureConnectionToNATSJetstream(t *testing.T) {

	t.Run("test connection is not successfully to TLS secured NATS jetstream broker without TLS config", func(t *testing.T) {
		subj := "secureConnectTestWithoutTLSConfig"

		_, err := NewController(
			"nats://nats-jetstream-tls:4222",
			WithStreamConfig(jetstream.StreamConfig{
				Name:     subj,
				Subjects: []string{subj},
			}),
			WithConsumerConfig(jetstream.ConsumerConfig{Name: "secureConnectTestWithoutTLSConfig"}),
		)
		assert.Error(t, err, "new connection to TLS secured NATS broker without TLS config should return a error")
	})

	t.Run("test connection is successfully to TLS secured NATS jetstream broker with TLS config", func(t *testing.T) {
		subj := "secureConnectTestWithTLSConfig"

		jc, err := NewController(
			"nats://nats-jetstream-tls:4222",
			WithStreamConfig(jetstream.StreamConfig{
				Name:     subj,
				Subjects: []string{subj},
			}),
			WithConsumerConfig(jetstream.ConsumerConfig{Name: "secureConnectTestWithoutTLSConfig"}),
			// just for testing use tls.Config with InsecureSkipVerify: true to skip server certificate validation for our self signed certificate
			WithConnectionOpts(nats.Secure(&tls.Config{InsecureSkipVerify: true})),
		)
		defer jc.Close()
		assert.NoError(t, err, "new connection to TLS secured NATS jetstream broker with TLS config should not return a error")
	})

	t.Run("test connection is not successfully to TLS secured NATS jetstream broker with TLS config and missing credentials", func(t *testing.T) {
		subj := "secureConnectTestWithTLSConfigAndWithoutCredentials"

		_, err := NewController(
			"nats://nats-jetstream-tls-basic-auth:4222",
			WithStreamConfig(jetstream.StreamConfig{
				Name:     subj,
				Subjects: []string{subj},
			}),
			WithConsumerConfig(jetstream.ConsumerConfig{Name: "secureConnectTestWithTLSConfigAndMissingBasicAuth"}),
			// just for testing use tls.Config with InsecureSkipVerify: true to skip server certificate validation for our self signed certificate
			WithConnectionOpts(nats.Secure(&tls.Config{InsecureSkipVerify: true})),
		)
		assert.Error(t, err, "new connection to TLS secured NATS jetstream broker with TLS config and missing credentials should return a error")
	})

	t.Run("test connection is successfully to TLS secured NATS jetstream broker with TLS config and credentials", func(t *testing.T) {
		subj := "secureConnectTestWithTLSConfigAndCredentials"

		_, err := NewController(
			"nats://nats-jetstream-tls-basic-auth:4222",
			WithStreamConfig(jetstream.StreamConfig{
				Name:     subj,
				Subjects: []string{subj},
			}),
			WithConsumerConfig(jetstream.ConsumerConfig{Name: "secureConnectTestWithTLSConfigAndMissingBasicAuth"}),
			// just for testing use tls.Config with InsecureSkipVerify: true to skip server certificate validation for our self signed certificate
			WithConnectionOpts(nats.Secure(&tls.Config{InsecureSkipVerify: true}), nats.UserInfo("user", "password")),
		)
		assert.NoError(t, err, "new connection to TLS secured NATS jetstream broker with TLS config and  credentials should return no error")
	})
}
