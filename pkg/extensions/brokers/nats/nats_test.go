package nats

import (
	"sync"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestValidateAckMechanism(t *testing.T) {
	nb, err := NewController("nats://nats:4222", WithQueueGroup("ValidateAckMechanism"))
	assert.NoError(t, err, "new controller should not return error")

	subj := "ValidateAckMechanism"

	t.Run("validate ack is not supported in core NATS", func(t *testing.T) {
		wg := sync.WaitGroup{}

		sub, err := nb.connection.Subscribe(subj, func(msg *nats.Msg) {
			defer wg.Done()
			err := msg.Ack()
			assert.ErrorIs(t, err, nats.ErrMsgNoReply, "core NATS should not support acks")
		})
		defer assert.NoError(t, sub.Drain())
		assert.NoError(t, err, "subscribe should not return error")

		wg.Add(1)
		err = nb.connection.Publish(subj, []byte("testmessage"))
		assert.NoError(t, err, "publish should not return error")

		wg.Wait()
	})

	t.Run("validate nak is not supported in core NATS", func(t *testing.T) {
		wg := sync.WaitGroup{}

		sub, err := nb.connection.Subscribe(subj, func(msg *nats.Msg) {
			defer wg.Done()
			err := msg.Nak()
			assert.ErrorIs(t, err, nats.ErrMsgNoReply, "core NATS should not support naks")
		})
		defer assert.NoError(t, sub.Drain())
		assert.NoError(t, err, "subscribe should not return error")

		wg.Add(1)
		err = nb.connection.Publish(subj, []byte("testmessage"))
		assert.NoError(t, err, "publish should not return error")

		wg.Wait()
	})
}
