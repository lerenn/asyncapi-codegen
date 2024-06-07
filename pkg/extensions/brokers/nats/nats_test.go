package nats

import (
	"crypto/tls"
	"fmt"
	"sync"
	"testing"

	testutil "github.com/TheSadlig/asyncapi-codegen/pkg/utils/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestValidateAckMechanism(t *testing.T) {
	subj := "CoreNatsValidateAckMechanism"
	nb, err := NewController(
		testutil.BrokerAddress(testutil.BrokerAddressParams{
			Schema:         "nats",
			DockerizedAddr: "nats",
			Port:           "4222",
		}),
		WithQueueGroup(subj))
	assert.NoError(t, err, "new controller should not return error")

	t.Run("validate ack is not supported in core NATS", func(t *testing.T) {
		wg := sync.WaitGroup{}
		subj = fmt.Sprintf("%s/%s", subj, "ack")

		sub, err := nb.connection.Subscribe(subj, func(msg *nats.Msg) {
			defer wg.Done()
			err := msg.Ack()
			assert.ErrorIs(t, err, nats.ErrMsgNoReply, "core NATS should not support acks")
		})

		// for some reason calling drain with defer assert.NoError(t, sub.Drain()) breaks the test so wrapping in a closure
		defer func(sub *nats.Subscription) {
			err := sub.Drain()
			assert.NoError(t, err, "Drain should not return a error")
		}(sub)

		assert.NoError(t, err, "subscribe should not return error")

		wg.Add(1)
		err = nb.connection.Publish(subj, []byte("testmessage"))
		assert.NoError(t, err, "publish should not return error")

		wg.Wait()
	})

	t.Run("validate nak is not supported in core NATS", func(t *testing.T) {
		wg := sync.WaitGroup{}
		subj = fmt.Sprintf("%s/%s", subj, "nak")

		sub, err := nb.connection.Subscribe(subj, func(msg *nats.Msg) {
			defer wg.Done()
			err := msg.Nak()
			assert.ErrorIs(t, err, nats.ErrMsgNoReply, "core NATS should not support naks")
		})

		// for some reason calling drain with defer assert.NoError(t, sub.Drain()) breaks the test so wrapping in a closure
		defer func(sub *nats.Subscription) {
			err := sub.Drain()
			assert.NoError(t, err, "Drain should not return a error")
		}(sub)
		assert.NoError(t, err, "subscribe should not return error")

		wg.Add(1)
		err = nb.connection.Publish(subj, []byte("testmessage"))
		assert.NoError(t, err, "publish should not return error")

		wg.Wait()
	})
}

//nolint:funlen
func TestSecureConnectionToNATSCore(t *testing.T) {
	// for testing with InsecureSkipVerify to skip server certificate validation for our self-signed certificate
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	t.Run("test connection is not successfully to TLS secured core NATS broker without TLS config", func(t *testing.T) {
		_, err := NewController(
			testutil.BrokerAddress(testutil.BrokerAddressParams{
				Schema:         "nats",
				DockerizedAddr: "nats-tls",
				DockerizedPort: "4222",
				LocalPort:      "4223",
			}),
			WithQueueGroup("secureConnectTest"))
		assert.Error(t, err, "new connection to TLS secured NATS broker without TLS config should return a error")
	})

	t.Run("test connection is successfully to TLS secured core NATS broker with TLS config", func(t *testing.T) {
		nb, err := NewController(
			testutil.BrokerAddress(testutil.BrokerAddressParams{
				Schema:         "nats",
				DockerizedAddr: "nats-tls",
				DockerizedPort: "4222",
				LocalPort:      "4223",
			}),
			WithQueueGroup("secureConnectTest"),
			WithConnectionOpts(nats.Secure(tlsConfig)))
		assert.NoError(t, err, "new connection to TLS secured NATS broker with TLS config should return no error")
		defer nb.Close()
	})

	t.Run("test connection is not successfully to TLS secured core NATS broker with TLS config and missing credentials",
		func(t *testing.T) {
			_, err := NewController(
				testutil.BrokerAddress(testutil.BrokerAddressParams{
					Schema:         "nats",
					DockerizedAddr: "nats-tls-basic-auth",
					DockerizedPort: "4222",
					LocalPort:      "4224",
				}),
				WithQueueGroup("secureConnectTest"),
				WithConnectionOpts(nats.Secure(tlsConfig)),
			)
			assert.Error(t, err, "new connection to TLS secured NATS broker with TLS config and missing credentials should return a error") //nolint:lll
		})

	t.Run("test connection is successfully to TLS secured core NATS broker with TLS config and credentials",
		func(t *testing.T) {
			nb, err := NewController(
				testutil.BrokerAddress(testutil.BrokerAddressParams{
					Schema:         "nats",
					DockerizedAddr: "nats-tls-basic-auth",
					DockerizedPort: "4222",
					LocalPort:      "4224",
				}),
				WithQueueGroup("secureConnectTest"),
				WithConnectionOpts(
					nats.Secure(tlsConfig),
					nats.UserInfo("user", "password"),
				),
			)
			assert.NoError(t, err,
				"new connection to TLS secured NATS broker with TLS config and basic credentials should return no error")
			defer nb.Close()
		})
}
