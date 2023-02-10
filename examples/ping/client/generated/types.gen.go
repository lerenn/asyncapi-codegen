// Package "generated" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package generated

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	// Generic error for AsyncAPI generated code
	ErrAsyncAPI = errors.New("error when using AsyncAPI")

	// ErrTimedOut is given when any timeout happen
	ErrTimedOut = fmt.Errorf("%w: time out", ErrAsyncAPI)

	// ErrNilBrokerController is raised when a nil broker controller is user
	ErrNilBrokerController = fmt.Errorf("%w: nil broker controller has been used", ErrAsyncAPI)

	// ErrNilAppSubscriber is raised when a nil app subscriber is user
	ErrNilAppSubscriber = fmt.Errorf("%w: nil app subscriber has been used", ErrAsyncAPI)

	// ErrNilClientSubscriber is raised when a nil client subscriber is user
	ErrNilClientSubscriber = fmt.Errorf("%w: nil client subscriber has been used", ErrAsyncAPI)

	// ErrAlreadySubscribedChannel is raised when a subscription is done twice
	// or more without unsubscribing
	ErrAlreadySubscribedChannel = fmt.Errorf("%w: the channel has already been subscribed", ErrAsyncAPI)
)

type Error struct {
	Channel string
	Err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("channel %q: err %v", e.Channel, e.Err)
}

// PingMessage is the message expected for 'Ping' channel
type PingMessage struct {
	// Headers will be used to fill the message headers
	Headers struct {
		// CorrelationID is Correlation ID set by client
		CorrelationID string `json:"correlation_id"`
	}

	// Payload will be inserted in the message payload
	Payload string
}

func NewPingMessage() PingMessage {
	var msg PingMessage

	// Set correlation ID
	msg.Headers.CorrelationID = uuid.New().String()

	return msg
}

// newPingMessageFromUniversalMessage will fill a new PingMessage with data from UniversalMessage
func newPingMessageFromUniversalMessage(um UniversalMessage) (PingMessage, error) {
	var msg PingMessage

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(um.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// Get correlation ID
	msg.Headers.CorrelationID = um.CorrelationID

	// TODO: run checks on msg type

	return msg, nil
}

// toUniversalMessage will generate an UniversalMessage from PingMessage data
func (msg PingMessage) toUniversalMessage() (UniversalMessage, error) {
	// TODO: implement checks on message

	// Marshal payload to JSON
	payload, err := json.Marshal(msg.Payload)
	if err != nil {
		return UniversalMessage{}, err
	}

	// Set correlation ID if it does not exist
	var correlationID string
	if msg.Headers.CorrelationID != "" {
		correlationID = msg.Headers.CorrelationID
	} else {
		correlationID = uuid.New().String()
	}

	return UniversalMessage{
		Payload:       payload,
		CorrelationID: correlationID,
	}, nil
}

// PongMessage is the message expected for 'Pong' channel
type PongMessage struct {
	// Headers will be used to fill the message headers
	Headers struct {
		// CorrelationID is Correlation ID set by client on corresponding request
		CorrelationID string `json:"correlation_id"`
	}

	// Payload will be inserted in the message payload
	Payload string
}

func NewPongMessage() PongMessage {
	var msg PongMessage

	// Set correlation ID
	msg.Headers.CorrelationID = uuid.New().String()

	return msg
}

// newPongMessageFromUniversalMessage will fill a new PongMessage with data from UniversalMessage
func newPongMessageFromUniversalMessage(um UniversalMessage) (PongMessage, error) {
	var msg PongMessage

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(um.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// Get correlation ID
	msg.Headers.CorrelationID = um.CorrelationID

	// TODO: run checks on msg type

	return msg, nil
}

// toUniversalMessage will generate an UniversalMessage from PongMessage data
func (msg PongMessage) toUniversalMessage() (UniversalMessage, error) {
	// TODO: implement checks on message

	// Marshal payload to JSON
	payload, err := json.Marshal(msg.Payload)
	if err != nil {
		return UniversalMessage{}, err
	}

	// Set correlation ID if it does not exist
	var correlationID string
	if msg.Headers.CorrelationID != "" {
		correlationID = msg.Headers.CorrelationID
	} else {
		correlationID = uuid.New().String()
	}

	return UniversalMessage{
		Payload:       payload,
		CorrelationID: correlationID,
	}, nil
}
