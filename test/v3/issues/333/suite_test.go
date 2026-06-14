package issue333

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/stretchr/testify/suite"
)

const eventsChannel = "v3.issue333.events"

func TestSuite(t *testing.T) {
	suite.Run(t, NewSuite())
}

type Suite struct {
	suite.Suite
}

func NewSuite() *Suite {
	return &Suite{}
}

// ackRecorder records whether a delivered message ended up acked or naked.
type ackRecorder struct {
	acked chan struct{}
	naked chan struct{}
}

func newAckRecorder() *ackRecorder {
	return &ackRecorder{
		acked: make(chan struct{}, 1),
		naked: make(chan struct{}, 1),
	}
}

func (a *ackRecorder) AckMessage() { a.acked <- struct{}{} }
func (a *ackRecorder) NakMessage() { a.naked <- struct{}{} }

// memBroker is a minimal in-memory broker that lets the test inject raw broker
// messages on a channel and observe how they are acknowledged.
type memBroker struct {
	mu   sync.Mutex
	subs map[string]extensions.BrokerChannelSubscription
}

func newMemBroker() *memBroker {
	return &memBroker{subs: make(map[string]extensions.BrokerChannelSubscription)}
}

func (b *memBroker) Publish(_ context.Context, _ string, _ extensions.BrokerMessage) error {
	return nil
}

func (b *memBroker) Subscribe(_ context.Context, channel string) (extensions.BrokerChannelSubscription, error) {
	sub := extensions.NewBrokerChannelSubscription(
		make(chan extensions.AcknowledgeableBrokerMessage, 8),
		make(chan any, 1),
	)
	sub.WaitForCancellationAsync(func() {})

	b.mu.Lock()
	b.subs[channel] = sub
	b.mu.Unlock()

	return sub, nil
}

// inject delivers a raw payload to the subscribers of the events channel and
// returns the recorder that observes its acknowledgement.
func (b *memBroker) inject(payload []byte) *ackRecorder {
	b.mu.Lock()
	sub := b.subs[eventsChannel]
	b.mu.Unlock()

	rec := newAckRecorder()
	sub.TransmitReceivedMessage(extensions.NewAcknowledgeableBrokerMessage(
		extensions.BrokerMessage{Payload: payload},
		rec,
	))

	return rec
}

// receiver captures the messages dispatched to each per-message callback.
type receiver struct {
	signedUp chan UserSignedUpMessage
	deleted  chan UserDeletedMessage
}

func newReceiver() *receiver {
	return &receiver{
		signedUp: make(chan UserSignedUpMessage, 1),
		deleted:  make(chan UserDeletedMessage, 1),
	}
}

func (r *receiver) ReceiveEventsOperationForUserSignedUpReceived(_ context.Context, msg UserSignedUpMessage) error {
	r.signedUp <- msg
	return nil
}

func (r *receiver) ReceiveEventsOperationForUserDeletedReceived(_ context.Context, msg UserDeletedMessage) error {
	r.deleted <- msg
	return nil
}

// TestDispatchesEachMessageTypeToItsOwnCallback ensures that two distinct
// message types received on the same channel are each dispatched to their own
// per-message callback using try-each-until-valid discrimination (issue #333).
func (suite *Suite) TestDispatchesEachMessageTypeToItsOwnCallback() {
	broker := newMemBroker()
	ctrl, err := NewAppController(broker)
	suite.Require().NoError(err)
	defer ctrl.Close(context.Background())

	rcv := newReceiver()
	suite.Require().NoError(ctrl.SubscribeToAllChannels(context.Background(), rcv))

	// A user_signed_up message must reach the UserSignedUp callback only.
	rec := broker.inject([]byte(`{"event":"user_signed_up","id":"signed-1"}`))
	select {
	case msg := <-rcv.signedUp:
		suite.Require().Equal("signed-1", msg.Payload.Id)
	case <-time.After(time.Second):
		suite.Fail("UserSignedUp callback was not called")
	}
	suite.requireAcked(rec)
	suite.Require().Empty(rcv.deleted, "UserDeleted callback must not be called for a user_signed_up message")

	// A user_deleted message must reach the UserDeleted callback only, even
	// though UserDeleted is tried after UserSignedUp in alphabetical order.
	rec = broker.inject([]byte(`{"event":"user_deleted","id":"deleted-1"}`))
	select {
	case msg := <-rcv.deleted:
		suite.Require().Equal("deleted-1", msg.Payload.Id)
	case <-time.After(time.Second):
		suite.Fail("UserDeleted callback was not called")
	}
	suite.requireAcked(rec)
	suite.Require().Empty(rcv.signedUp, "UserSignedUp callback must not be called for a user_deleted message")
}

// TestUnmatchedMessageIsNakedAndReported ensures a message that matches none of
// the expected types is naked and surfaced through the error handler with
// ErrNoMatchingMessage (issue #333).
func (suite *Suite) TestUnmatchedMessageIsNakedAndReported() {
	handlerErr := make(chan error, 1)
	broker := newMemBroker()
	ctrl, err := NewAppController(broker, WithErrorHandler(
		func(_ context.Context, _ string, _ *extensions.AcknowledgeableBrokerMessage, opErr error) {
			handlerErr <- opErr
		},
	))
	suite.Require().NoError(err)
	defer ctrl.Close(context.Background())

	rcv := newReceiver()
	suite.Require().NoError(ctrl.SubscribeToAllChannels(context.Background(), rcv))

	rec := broker.inject([]byte(`{"event":"unknown","id":"x"}`))

	select {
	case opErr := <-handlerErr:
		suite.Require().ErrorIs(opErr, extensions.ErrNoMatchingMessage)
	case <-time.After(time.Second):
		suite.Fail("error handler was not called for an unmatched message")
	}

	suite.requireNaked(rec)
	suite.Require().Empty(rcv.signedUp)
	suite.Require().Empty(rcv.deleted)
}

// TestUnsubscribeStopsDispatch ensures that unsubscribing a single message type
// removes only its handler while the others keep receiving (issue #333).
func (suite *Suite) TestUnsubscribeStopsDispatch() {
	broker := newMemBroker()
	ctrl, err := NewAppController(broker)
	suite.Require().NoError(err)
	defer ctrl.Close(context.Background())

	rcv := newReceiver()
	suite.Require().NoError(ctrl.SubscribeToAllChannels(context.Background(), rcv))

	// Unsubscribe only the UserSignedUp handler.
	ctrl.UnsubscribeFromReceiveEventsOperationForUserSignedUp(context.Background())

	// UserDeleted must still be dispatched.
	rec := broker.inject([]byte(`{"event":"user_deleted","id":"deleted-2"}`))
	select {
	case msg := <-rcv.deleted:
		suite.Require().Equal("deleted-2", msg.Payload.Id)
	case <-time.After(time.Second):
		suite.Fail("UserDeleted callback was not called after unsubscribing UserSignedUp")
	}
	suite.requireAcked(rec)
}

func (suite *Suite) requireAcked(rec *ackRecorder) {
	select {
	case <-rec.acked:
	case <-rec.naked:
		suite.Fail("message was naked but should have been acked")
	case <-time.After(time.Second):
		suite.Fail("message was neither acked nor naked")
	}
}

func (suite *Suite) requireNaked(rec *ackRecorder) {
	select {
	case <-rec.naked:
	case <-rec.acked:
		suite.Fail("message was acked but should have been naked")
	case <-time.After(time.Second):
		suite.Fail("message was neither acked nor naked")
	}
}
