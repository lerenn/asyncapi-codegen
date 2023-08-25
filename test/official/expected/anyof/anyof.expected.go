// Package "anyof" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package anyof

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	aapiContext "github.com/lerenn/asyncapi-codegen/pkg/context"
	"github.com/lerenn/asyncapi-codegen/pkg/log"
	"github.com/lerenn/asyncapi-codegen/pkg/middleware"
)

// AppSubscriber represents all handlers that are expecting messages for App
type AppSubscriber interface {
	// Test
	Test(ctx context.Context, msg TestMessagesMessage, done bool)
}

// AppController is the structure that provides publishing capabilities to the
// developer and and connect the broker with the App
type AppController struct {
	brokerController BrokerController
	stopSubscribers  map[string]chan interface{}
	logger           log.Interface
	middlewares      []middleware.Interface
}

// NewAppController links the App to the broker
func NewAppController(bs BrokerController) (*AppController, error) {
	if bs == nil {
		return nil, ErrNilBrokerController
	}

	return &AppController{
		brokerController: bs,
		stopSubscribers:  make(map[string]chan interface{}),
		logger:           log.Silent{},
		middlewares:      make([]middleware.Interface, 0),
	}, nil
}

// SetLogger attaches a logger that will log operations on controller
func (c *AppController) SetLogger(logger log.Interface) {
	c.logger = logger
	c.brokerController.SetLogger(logger)
}

// AddMiddlewares attaches middlewares that will be executed when sending or
// receiving messages
func (c *AppController) AddMiddlewares(middleware ...middleware.Interface) {
	c.middlewares = append(c.middlewares, middleware...)
}

func (c AppController) executeMiddlewares(ctx context.Context, um UniversalMessage) {
	for _, m := range c.middlewares {
		m(ctx, um.Payload)
	}
}

func addAppContextValues(ctx context.Context, path, operation string) context.Context {
	ctx = context.WithValue(ctx, aapiContext.KeyIsModule, "asyncapi")
	ctx = context.WithValue(ctx, aapiContext.KeyIsProvider, "app")
	ctx = context.WithValue(ctx, aapiContext.KeyIsChannel, path)
	return context.WithValue(ctx, aapiContext.KeyIsOperation, operation)
}

// Close will clean up any existing resources on the controller
func (c *AppController) Close(ctx context.Context) {
	// Unsubscribing remaining channels
	c.logger.Info(ctx, "Closing App controller")
	c.UnsubscribeAll(ctx)
}

// SubscribeAll will subscribe to channels without parameters on which the app is expecting messages.
// For channels with parameters, they should be subscribed independently.
func (c *AppController) SubscribeAll(ctx context.Context, as AppSubscriber) error {
	if as == nil {
		return ErrNilAppSubscriber
	}

	if err := c.SubscribeTest(ctx, as.Test); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (c *AppController) UnsubscribeAll(ctx context.Context) {
	// Unsubscribe channels with no parameters (if any)
	c.UnsubscribeTest(ctx)

	// Unsubscribe remaining channels
	for n, stopChan := range c.stopSubscribers {
		stopChan <- true
		delete(c.stopSubscribers, n)
	}
}

// SubscribeTest will subscribe to new messages from 'test' channel.
//
// Callback function 'fn' will be called each time a new message is received.
// The 'done' argument indicates when the subscription is canceled and can be
// used to clean up resources.
func (c *AppController) SubscribeTest(ctx context.Context, fn func(ctx context.Context, msg TestMessagesMessage, done bool)) error {
	// Get channel path
	path := "test"

	// Set context
	ctx = addAppContextValues(ctx, path, "subscribe")
	ctx = context.WithValue(ctx, aapiContext.KeyIsDirection, "reception")

	// Check if there is already a subscription
	_, exists := c.stopSubscribers[path]
	if exists {
		err := fmt.Errorf("%w: %q channel is already subscribed", ErrAlreadySubscribedChannel, path)
		c.logger.Error(ctx, err.Error())
		return err
	}

	// Subscribe to broker channel
	c.logger.Info(ctx, "Subscribing to channel")
	msgs, stop, err := c.brokerController.Subscribe(ctx, path)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return err
	}

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for {
			// Wait for next message
			um, open := <-msgs

			// Execute middlewares
			c.executeMiddlewares(ctx, um)

			// Process message
			msg, err := newTestMessagesMessageFromUniversalMessage(um)
			if err != nil {
				ctx = context.WithValue(ctx, aapiContext.KeyIsMessage, um)
				c.logger.Error(ctx, err.Error())
			}
			ctx = context.WithValue(ctx, aapiContext.KeyIsMessage, msg)

			// Send info if message is correct or susbcription is closed
			if err == nil || !open {
				c.logger.Info(ctx, "Received new message")
				fn(ctx, msg, !open)
			}

			// If subscription is closed, then exit the function
			if !open {
				return
			}
		}
	}()

	// Add the stop channel to the inside map
	c.stopSubscribers[path] = stop

	return nil
}

// UnsubscribeTest will unsubscribe messages from 'test' channel
func (c *AppController) UnsubscribeTest(ctx context.Context) {
	// Get channel path
	path := "test"

	// Set context
	ctx = addAppContextValues(ctx, path, "unsubscribe")

	// Get stop channel
	stopChan, exists := c.stopSubscribers[path]
	if !exists {
		return
	}

	// Stop the channel and remove the entry
	c.logger.Info(ctx, "Unsubscribing from channel")
	stopChan <- true
	delete(c.stopSubscribers, path)
}

// ClientController is the structure that provides publishing capabilities to the
// developer and and connect the broker with the Client
type ClientController struct {
	brokerController BrokerController
	stopSubscribers  map[string]chan interface{}
	logger           log.Interface
	middlewares      []middleware.Interface
}

// NewClientController links the Client to the broker
func NewClientController(bs BrokerController) (*ClientController, error) {
	if bs == nil {
		return nil, ErrNilBrokerController
	}

	return &ClientController{
		brokerController: bs,
		stopSubscribers:  make(map[string]chan interface{}),
		logger:           log.Silent{},
		middlewares:      make([]middleware.Interface, 0),
	}, nil
}

// SetLogger attaches a logger that will log operations on controller
func (c *ClientController) SetLogger(logger log.Interface) {
	c.logger = logger
	c.brokerController.SetLogger(logger)
}

// AddMiddlewares attaches middlewares that will be executed when sending or
// receiving messages
func (c *ClientController) AddMiddlewares(middleware ...middleware.Interface) {
	c.middlewares = append(c.middlewares, middleware...)
}

func (c ClientController) executeMiddlewares(ctx context.Context, um UniversalMessage) {
	for _, m := range c.middlewares {
		m(ctx, um.Payload)
	}
}

func addClientContextValues(ctx context.Context, path, operation string) context.Context {
	ctx = context.WithValue(ctx, aapiContext.KeyIsModule, "asyncapi")
	ctx = context.WithValue(ctx, aapiContext.KeyIsProvider, "client")
	ctx = context.WithValue(ctx, aapiContext.KeyIsChannel, path)
	return context.WithValue(ctx, aapiContext.KeyIsOperation, operation)
}

// Close will clean up any existing resources on the controller
func (c *ClientController) Close(ctx context.Context) {
	// Unsubscribing remaining channels
}

// PublishTest will publish messages to 'test' channel
func (c *ClientController) PublishTest(ctx context.Context, msg TestMessagesMessage) error {
	// Get channel path
	path := "test"

	// Set context
	ctx = addClientContextValues(ctx, path, "publish")
	ctx = context.WithValue(ctx, aapiContext.KeyIsMessage, msg)
	ctx = context.WithValue(ctx, aapiContext.KeyIsDirection, "publication")

	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Execute middlewares
	c.executeMiddlewares(ctx, um)

	// Publish on event broker
	c.logger.Info(ctx, "Publishing to channel")

	return c.brokerController.Publish(ctx, path, um)
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
	// SetLogger set a logger that will log operations on broker controller
	SetLogger(logger log.Interface)

	// Publish a message to the broker
	Publish(ctx context.Context, channel string, mw UniversalMessage) error

	// Subscribe to messages from the broker
	Subscribe(ctx context.Context, channel string) (msgs chan UniversalMessage, stop chan interface{}, err error)

	// SetQueueName sets the name of the queue that will be used by the broker
	SetQueueName(name string)
}

var (
	// Generic error for AsyncAPI generated code
	ErrAsyncAPI = errors.New("error when using AsyncAPI")

	// ErrContextCanceled is given when a given context is canceled
	ErrContextCanceled = fmt.Errorf("%w: context canceled", ErrAsyncAPI)

	// ErrNilBrokerController is raised when a nil broker controller is user
	ErrNilBrokerController = fmt.Errorf("%w: nil broker controller has been used", ErrAsyncAPI)

	// ErrNilAppSubscriber is raised when a nil app subscriber is user
	ErrNilAppSubscriber = fmt.Errorf("%w: nil app subscriber has been used", ErrAsyncAPI)

	// ErrNilClientSubscriber is raised when a nil client subscriber is user
	ErrNilClientSubscriber = fmt.Errorf("%w: nil client subscriber has been used", ErrAsyncAPI)

	// ErrAlreadySubscribedChannel is raised when a subscription is done twice
	// or more without unsubscribing
	ErrAlreadySubscribedChannel = fmt.Errorf("%w: the channel has already been subscribed", ErrAsyncAPI)

	// ErrSubscriptionCanceled is raised when expecting something and the subscription has been canceled before it happens
	ErrSubscriptionCanceled = fmt.Errorf("%w: the subscription has been canceled", ErrAsyncAPI)
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
