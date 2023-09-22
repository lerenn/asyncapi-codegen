package extensions

import (
	"context"
)

// BrokerChannelSubscription is a struct that contains every returned structures
// when subscribing a channel.
type BrokerChannelSubscription struct {
	messages chan BrokerMessage
	cancel   chan any
}

// NewBrokerChannelSubscription creates a new broker channel subscription based
// on the channels used to receive message and cancel the subscription.
func NewBrokerChannelSubscription(messages chan BrokerMessage, cancel chan any) BrokerChannelSubscription {
	return BrokerChannelSubscription{
		messages: messages,
		cancel:   cancel,
	}
}

// TransmitReceivedMessage should only be used by the broker to transmit the
// new received messages to the user.
func (bcs BrokerChannelSubscription) TransmitReceivedMessage(msg BrokerMessage) {
	bcs.messages <- msg
}

// MessagesChannel returns the channel that will get the received messages from
// broker and from which the user should listen.
func (bcs BrokerChannelSubscription) MessagesChannel() <-chan BrokerMessage {
	return bcs.messages
}

// WaitForCancellationAsync should be used by the broker only to wait for user request
// for cancellation. As it is asynchronous, it will return immediately after the call.
func (bcs BrokerChannelSubscription) WaitForCancellationAsync(cleanup func()) {
	go func() {
		// Wait for cancel request
		<-bcs.cancel

		// Execute cleanup function
		cleanup()

		// Close messages in order to avoid new messages
		close(bcs.messages)

		// Close cancel to let listeners know that the cancellation is complete
		close(bcs.cancel)
	}()
}

// Cancel cancels the subscription from user perspective. It will ask for clean
// up on broker, which will return when finished to avoid dangling resources, such
// as non-existent queue listeners on (broker) server side.
func (bcs BrokerChannelSubscription) Cancel() {
	// Send a cancellation request
	bcs.cancel <- true

	// Wait for the cancellation to be effective
	for open := true; open; {
		_, open = <-bcs.cancel
	}
}

// BrokerMessage is a wrapper that will contain all information regarding a message.
type BrokerMessage struct {
	Headers map[string][]byte
	Payload []byte
}

// IsUninitialized check if the BrokerMessage is at zero value, i.e. the
// uninitialized structure. It can be used to check that a channel is closed.
func (bm BrokerMessage) IsUninitialized() bool {
	return bm.Headers == nil && bm.Payload == nil
}

// BrokerController represents the functions that should be implemented to connect
// the broker to the application or the user.
type BrokerController interface {
	// Publish a message to the broker
	Publish(ctx context.Context, channel string, mw BrokerMessage) error

	// Subscribe to messages from the broker
	Subscribe(ctx context.Context, channel string) (BrokerChannelSubscription, error)
}
