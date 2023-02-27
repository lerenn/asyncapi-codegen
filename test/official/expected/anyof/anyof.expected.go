// Package "anyof" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package anyof

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AppSubscriber represents all handlers that are expecting messages for App
type AppSubscriber interface {
	// Test
	Test(msg TestMessagesMessage)
}

// AppController is the structure that provides publishing capabilities to the
// developer and and connect the broker with the App
type AppController struct {
	brokerController BrokerController
	stopSubscribers  map[string]chan interface{}
	errChan          chan Error
}

// NewAppController links the App to the broker
func NewAppController(bs BrokerController) (*AppController, error) {
	if bs == nil {
		return nil, ErrNilBrokerController
	}

	return &AppController{
		brokerController: bs,
		stopSubscribers:  make(map[string]chan interface{}),
		errChan:          make(chan Error, 256),
	}, nil
}

// Errors will give back the channel that contains errors and that you can listen to handle errors
// Please take a look at Error struct form information on error
func (c AppController) Errors() <-chan Error {
	return c.errChan
}

// Close will clean up any existing resources on the controller
func (c *AppController) Close() {
	c.UnsubscribeAll()
	close(c.errChan)
}

// SubscribeAll will subscribe to channels on which the app is expecting messages
func (c *AppController) SubscribeAll(as AppSubscriber) error {
	if as == nil {
		return ErrNilAppSubscriber
	}

	if err := c.SubscribeTest(as.Test); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (c *AppController) UnsubscribeAll() {
	c.UnsubscribeTest()
}

// SubscribeTest will subscribe to new messages from 'test' channel
func (c *AppController) SubscribeTest(fn func(msg TestMessagesMessage)) error {
	// Check if there is already a subscription
	_, exists := c.stopSubscribers["test"]
	if exists {
		return fmt.Errorf("%w: test channel is already subscribed", ErrAlreadySubscribedChannel)
	}

	// Subscribe to broker channel
	msgs, stop, err := c.brokerController.Subscribe("test")
	if err != nil {
		return err
	}

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for um, open := <-msgs; open; um, open = <-msgs {
			msg, err := newTestMessagesMessageFromUniversalMessage(um)
			if err != nil {
				c.handleError("test", err)
			} else {
				fn(msg)
			}
		}
	}()

	// Add the stop channel to the inside map
	c.stopSubscribers["test"] = stop

	return nil
}

// UnsubscribeTest will unsubscribe messages from 'test' channel
func (c *AppController) UnsubscribeTest() {
	stopChan, exists := c.stopSubscribers["test"]
	if !exists {
		return
	}

	stopChan <- true
	delete(c.stopSubscribers, "test")
}

func (c *AppController) handleError(channelName string, err error) {
	// Wrap error with the channel name
	errWrapped := Error{
		Channel: channelName,
		Err:     err,
	}

	// Send it to the error channel
	select {
	case c.errChan <- errWrapped:
	default:
		// Drop error if it's full or closed
	}
}

// ClientController is the structure that provides publishing capabilities to the
// developer and and connect the broker with the Client
type ClientController struct {
	brokerController BrokerController
	stopSubscribers  map[string]chan interface{}
	errChan          chan Error
}

// NewClientController links the Client to the broker
func NewClientController(bs BrokerController) (*ClientController, error) {
	if bs == nil {
		return nil, ErrNilBrokerController
	}

	return &ClientController{
		brokerController: bs,
		stopSubscribers:  make(map[string]chan interface{}),
		errChan:          make(chan Error, 256),
	}, nil
}

// Errors will give back the channel that contains errors and that you can listen to handle errors
// Please take a look at Error struct form information on error
func (c ClientController) Errors() <-chan Error {
	return c.errChan
}

// Close will clean up any existing resources on the controller
func (c *ClientController) Close() {
	close(c.errChan)
}

// PublishTest will publish messages to 'test' channel
func (c *ClientController) PublishTest(msg TestMessagesMessage) error {
	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Publish on event broker
	return c.brokerController.Publish("test", um)
}

func (c *ClientController) handleError(channelName string, err error) {
	// Wrap error with the channel name
	errWrapped := Error{
		Channel: channelName,
		Err:     err,
	}

	// Send it to the error channel
	select {
	case c.errChan <- errWrapped:
	default:
		// Drop error if it's full or closed
	}
}

const (
	// CorrelationIDField is the name of the field that will contain the correlation ID
	CorrelationIDField = "correlation_id"
)

// UniversalMessage is a wrapper that will contain all information regarding a message
type UniversalMessage struct {
	CorrelationID *string
	Payload       []byte
}

// BrokerController represents the functions that should be implemented to connect
// the broker to the application or the client
type BrokerController interface {
	// Publish a message to the broker
	Publish(channel string, mw UniversalMessage) error

	// Subscribe to messages from the broker
	Subscribe(channel string) (msgs chan UniversalMessage, stop chan interface{}, err error)
}

var (
	// Generic error for AsyncAPI generated code
	ErrAsyncAPI = errors.New("error when using AsyncAPI")

	// ErrContextCancelled is given when a given context is cancelled
	ErrContextCancelled = fmt.Errorf("%w: context cancelled", ErrAsyncAPI)

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

type MessageWithCorrelationID interface {
	CorrelationID() string
}

type Error struct {
	Channel string
	Err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("channel %q: err %v", e.Channel, e.Err)
}

// TestMessagesMessage is the message expected for 'TestMessages' channel
type TestMessagesMessage struct {
	// Payload will be inserted in the message payload
	Payload struct {
		Key  *string `json:"key"`
		Key2 *string `json:"key2"`
	}
}

func NewTestMessagesMessage() TestMessagesMessage {
	var msg TestMessagesMessage

	return msg
}

// newTestMessagesMessageFromUniversalMessage will fill a new TestMessagesMessage with data from UniversalMessage
func newTestMessagesMessageFromUniversalMessage(um UniversalMessage) (TestMessagesMessage, error) {
	var msg TestMessagesMessage

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(um.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toUniversalMessage will generate an UniversalMessage from TestMessagesMessage data
func (msg TestMessagesMessage) toUniversalMessage() (UniversalMessage, error) {
	// TODO: implement checks on message

	// Marshal payload to JSON
	payload, err := json.Marshal(msg.Payload)
	if err != nil {
		return UniversalMessage{}, err
	}

	return UniversalMessage{
		Payload: payload,
	}, nil
}

// ObjectWithKeySchema is a schema from the AsyncAPI specification required in messages
type ObjectWithKeySchema struct {
	Key *string `json:"key"`
}

// ObjectWithKey2Schema is a schema from the AsyncAPI specification required in messages
type ObjectWithKey2Schema struct {
	Key2 *string `json:"key2"`
}
