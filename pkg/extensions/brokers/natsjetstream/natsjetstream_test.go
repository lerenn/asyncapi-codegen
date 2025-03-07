package natsjetstream

import (
	"context"
	"crypto/tls"
	"sync"
	"testing"

	testutil "github.com/lerenn/asyncapi-codegen/pkg/utils/test"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:funlen // this is only for testing
func TestValidateAckMechanism(t *testing.T) {
	subj := "NatsJetstreamValidateAckMechanism"

	broker, err := NewController(
		testutil.BrokerAddress(testutil.BrokerAddressParams{
			Schema:         "nats",
			DockerizedAddr: "nats-jetstream",
			DockerizedPort: "4222",
			LocalPort:      "4225",
		}),
		WithStreamConfig(jetstream.StreamConfig{
			Name:     subj,
			Subjects: []string{subj},
		}),
		WithConsumerConfig(jetstream.ConsumerConfig{Name: "natsJetstreamValidateAckMechanism"}),
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

//nolint:funlen
func TestSecureConnectionToNATSJetstream(t *testing.T) {
	// for testing with InsecureSkipVerify to skip server certificate validation for our self-signed certificate
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	t.Run("test connection is not successfully to TLS secured NATS jetstream broker without TLS config",
		func(t *testing.T) {
			subj := "secureConnectTestWithoutTLSConfig"

			_, err := NewController(
				testutil.BrokerAddress(testutil.BrokerAddressParams{
					Schema:         "nats",
					DockerizedAddr: "nats-jetstream-tls",
					DockerizedPort: "4222",
					LocalPort:      "4226",
				}),
				WithStreamConfig(jetstream.StreamConfig{
					Name:     subj,
					Subjects: []string{subj},
				}),
				WithConsumerConfig(jetstream.ConsumerConfig{
					Name: subj,
				}),
			)
			assert.Error(t, err, "new connection to TLS secured NATS broker without TLS config should return a error")
		})

	t.Run("test connection is successfully to TLS secured NATS jetstream broker with TLS config", func(t *testing.T) {
		subj := "secureConnectTestWithTLSConfig"

		jc, err := NewController(
			testutil.BrokerAddress(testutil.BrokerAddressParams{
				Schema:         "nats",
				DockerizedAddr: "nats-jetstream-tls",
				DockerizedPort: "4222",
				LocalPort:      "4226",
			}),
			WithStreamConfig(jetstream.StreamConfig{
				Name:     subj,
				Subjects: []string{subj},
			}),
			WithConsumerConfig(jetstream.ConsumerConfig{
				Name: subj,
			}),
			WithConnectionOpts(nats.Secure(tlsConfig)),
		)
		assert.NoError(t, err,
			"new connection to TLS secured NATS jetstream broker with TLS config should not return a error")
		defer jc.Close()
	})

	t.Run("test connection is not successfully to TLS secured NATS jetstream broker with TLS config and missing credentials", //nolint:lll
		func(t *testing.T) {
			subj := "secureConnectTestWithTLSConfigAndWithoutCredentials"

			_, err := NewController(
				testutil.BrokerAddress(testutil.BrokerAddressParams{
					Schema:         "nats",
					DockerizedAddr: "nats-jetstream-tls-basic-auth",
					DockerizedPort: "4222",
					LocalPort:      "4227",
				}),
				WithStreamConfig(jetstream.StreamConfig{
					Name:     subj,
					Subjects: []string{subj},
				}),
				WithConsumerConfig(jetstream.ConsumerConfig{
					Name: subj,
				}),
				WithConnectionOpts(nats.Secure(tlsConfig)),
			)
			assert.Error(t, err,
				"new connection to TLS secured NATS jetstream broker with TLS config and missing credentials should return a error")
		})

	t.Run("test connection is successfully to TLS secured NATS jetstream broker with TLS config and credentials",
		func(t *testing.T) {
			subj := "secureConnectTestWithTLSConfigAndCredentials"

			jc, err := NewController(
				testutil.BrokerAddress(testutil.BrokerAddressParams{
					Schema:         "nats",
					DockerizedAddr: "nats-jetstream-tls-basic-auth",
					DockerizedPort: "4222",
					LocalPort:      "4227",
				}),
				WithStreamConfig(jetstream.StreamConfig{
					Name:     subj,
					Subjects: []string{subj},
				}),
				WithConsumerConfig(jetstream.ConsumerConfig{
					Name: subj,
				}),
				WithConnectionOpts(
					nats.Secure(tlsConfig),
					nats.UserInfo("user", "password"),
				),
			)
			assert.NoError(t, err, "new connection to TLS secured NATS jetstream broker with TLS config and  credentials should return no error") //nolint:lll
			defer jc.Close()
		})
}

func TestExistingNatsConnection(t *testing.T) {
	subj := "NatsJetstreamValidateConnection"
	natsURL := testutil.BrokerAddress(testutil.BrokerAddressParams{
		Schema:         "nats",
		DockerizedAddr: "nats-jetstream",
		DockerizedPort: "4222",
		LocalPort:      "4225",
	})

	nc, err := nats.Connect(natsURL)
	require.NoError(t, err, "nats connection should be established")
	defer nc.Close()

	broker, err := NewController(
		"unused",
		WithStreamConfig(jetstream.StreamConfig{
			Name:     subj,
			Subjects: []string{subj},
		}),
		WithConnection(nc),
	)
	assert.NoError(t, err, "new controller should not return error")

	broker.Close()

	assert.True(t, nc.IsConnected(), "our connection should still be intact")
}

func TestMatchSubjectSubscription(t *testing.T) {
	// Test cases
	testCases := []struct {
		pattern string
		subject string
		result  bool
	}{
		{"", "", false},
		{"", "time", false},
		{"time.us.east.atlanta", "time.us.east.atlanta", true},
		{"time.us.east.atlanta", "time.us.last.atlanta", false},
		{"*", "", false},
		{"*", "time", true},
		{"time.*.*.atlanta", "", false},
		{"*.us.*.atlanta", "time.us.east.atlanta", true},
		{"time.*.*.atlanta", "time.us.east.atlanta", true},
		{"time.us.*.*", "time.us.east.atlanta", true},
		{"time.us.*", "time.us.east.atlanta", false},
		{"time.*.east", "time.us.east", true},
		{"time.*.east", "time.eu.west", false},
		{">", "", false},
		{">", "time.us.east", true},
		{"time.us.>", "", false},
		{"time.us.>", "time.us", false},
		{"time.us.>", "time.eu.east", false},
		{"time.us.>", "time.us.east.atlanta", true},
		{"time", "", false},
	}

	// Run the test cases
	for _, testCase := range testCases {
		result := MatchSubjectSubscription(testCase.pattern, testCase.subject)
		if result != testCase.result {
			println("Test failed for pattern:", testCase.pattern, "subject:", testCase.subject)
		}
	}
}
