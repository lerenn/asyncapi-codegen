package kafka

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

// scriptedReader returns a predefined sequence of (message, error) results and
// then io.EOF, letting us drive the consumption loop without a live broker.
type scriptedReader struct {
	steps []scriptedStep
	index int
}

type scriptedStep struct {
	msg kafkago.Message
	err error
}

func (r *scriptedReader) next() (kafkago.Message, error) {
	if r.index >= len(r.steps) {
		return kafkago.Message{}, io.EOF
	}
	step := r.steps[r.index]
	r.index++

	return step.msg, step.err
}

func (r *scriptedReader) ReadMessage(_ context.Context) (kafkago.Message, error)  { return r.next() }
func (r *scriptedReader) FetchMessage(_ context.Context) (kafkago.Message, error) { return r.next() }
func (r *scriptedReader) CommitMessages(_ context.Context, _ ...kafkago.Message) error {
	return nil
}

// TestConsumerContinuesAfterTransientError ensures that a transient read error
// does not kill the subscription: the handler must keep consuming and still
// deliver the next message (issue #285). Before the fix, any non-EOF error made
// the handler return, so the message after the error was never delivered.
func TestConsumerContinuesAfterTransientError(t *testing.T) {
	logger := extensions.Logger(extensions.DummyLogger{})

	reader := &scriptedReader{steps: []scriptedStep{
		{err: errors.New("transient: connection reset by peer")},
		{msg: kafkago.Message{Value: []byte("hello")}},
		// Exhausted afterwards -> io.EOF -> handler returns.
	}}

	sub := extensions.NewBrokerChannelSubscription(
		make(chan extensions.AcknowledgeableBrokerMessage, 8),
		make(chan any, 1),
	)

	done := make(chan struct{})
	go func() {
		autoCommitMessagesHandler(&logger)(context.Background(), reader, sub)
		close(done)
	}()

	// The message published after the transient error must still arrive.
	select {
	case msg := <-sub.MessagesChannel():
		require.Equal(t, "hello", string(msg.Payload))
	case <-time.After(2 * time.Second):
		t.Fatal("message after transient error was not delivered (handler exited early)")
	}

	// After exhausting the script, the reader returns io.EOF and the handler stops.
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("handler did not stop on io.EOF")
	}
}

// TestConsumerStopsOnContextCancellation ensures the handler exits when the
// context is cancelled instead of busy-looping on the resulting error.
func TestConsumerStopsOnContextCancellation(t *testing.T) {
	logger := extensions.Logger(extensions.DummyLogger{})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// A reader that always returns the context error once cancelled.
	reader := &scriptedReader{steps: []scriptedStep{
		{err: context.Canceled},
	}}

	sub := extensions.NewBrokerChannelSubscription(
		make(chan extensions.AcknowledgeableBrokerMessage, 1),
		make(chan any, 1),
	)

	done := make(chan struct{})
	go func() {
		autoCommitMessagesHandler(&logger)(ctx, reader, sub)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("handler did not stop after context cancellation")
	}
}
