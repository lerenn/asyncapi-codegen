package nats

import (
	"crypto/tls"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

//nolint:funlen // this is only for testing
func TestValidateAckMechanism(t *testing.T) {
	subj := "CoreNatsValidateAckMechanism"
	nb, err := NewController("nats://nats:4222", WithQueueGroup(subj))
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

func TestSecureConnectionToNATSCore(t *testing.T) {

	t.Run("test connection is not successfully to TLS secured core NATS broker without TLS config", func(t *testing.T) {
		_, err := NewController("nats://nats-tls:4222", WithQueueGroup("secureConnectTest"))
		assert.Error(t, err, "new connection to TLS secured NATS broker without TLS config should return a error")
	})

	t.Run("test connection is successfully to TLS secured core NATS broker with TLS config", func(t *testing.T) {
		nb, err := NewController("nats://nats-tls:4222", WithQueueGroup("secureConnectTest"),
			// just for testing use tls.Config with InsecureSkipVerify: true to skip server certificate validation for our self signed certificate
			WithConnectionOpts(nats.Secure(&tls.Config{InsecureSkipVerify: true})))
		defer nb.Close()
		assert.NoError(t, err, "new connection to TLS secured NATS broker with TLS config should return no error")
	})

	t.Run("test connection is not successfully to TLS secured core NATS broker with TLS config and missing credentials", func(t *testing.T) {
		_, err := NewController("nats://nats-tls-basic-auth:4222", WithQueueGroup("secureConnectTest"),
			// just for testing use tls.Config with InsecureSkipVerify: true to skip server certificate validation for our self signed certificate
			WithConnectionOpts(nats.Secure(&tls.Config{InsecureSkipVerify: true})))
		assert.Error(t, err, "new connection to TLS secured NATS broker with TLS config and missing credentials should return a error")
	})

	t.Run("test connection is successfully to TLS secured core NATS broker with TLS config and credentials", func(t *testing.T) {
		nb, err := NewController("nats://nats-tls-basic-auth:4222", WithQueueGroup("secureConnectTest"),
			// just for testing use tls.Config with InsecureSkipVerify: true to skip server certificate validation for our self signed certificate
			WithConnectionOpts(nats.Secure(&tls.Config{InsecureSkipVerify: true}), nats.UserInfo("user", "password")))
		defer nb.Close()
		assert.NoError(t, err, "new connection to TLS secured NATS broker with TLS config and basic credentials should return no error")
	})
}
