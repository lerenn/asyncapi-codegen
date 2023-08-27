// Package "issues" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package issues

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	apiContext "github.com/lerenn/asyncapi-codegen/pkg/context"
	"github.com/lerenn/asyncapi-codegen/pkg/log"
	"github.com/lerenn/asyncapi-codegen/pkg/middleware"
)

// ClientSubscriber represents all handlers that are expecting messages for Client
type ClientSubscriber interface {
	// Chat
	Chat(ctx context.Context, msg ChatMessage, done bool)

	// Status
	Status(ctx context.Context, msg StatusMessage, done bool)
}

// ClientController is the structure that provides publishing capabilities to the
// developer and and connect the broker with the Client
type ClientController struct {
	brokerController BrokerController
	stopSubscribers  map[string]chan interface{}
	logger           log.Interface
	middlewares      []middleware.Middleware
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
		middlewares:      make([]middleware.Middleware, 0),
	}, nil
}

// SetLogger attaches a logger that will log operations on controller
func (c *ClientController) SetLogger(logger log.Interface) {
	c.logger = logger
	c.brokerController.SetLogger(logger)
}

// AddMiddlewares attaches middlewares that will be executed when sending or
// receiving messages
func (c *ClientController) AddMiddlewares(middleware ...middleware.Middleware) {
	c.middlewares = append(c.middlewares, middleware...)
}

func (c ClientController) wrapMiddlewares(middlewares []middleware.Middleware, last middleware.Next) func(ctx context.Context) {
	var called bool

	// If there is no more middleware
	if len(middlewares) == 0 {
		return func(ctx context.Context) {
			if !called {
				called = true
				last(ctx)
			}
		}
	}

	// Wrap middleware into a check function that will call execute the middleware
	// and call the next wrapped middleware if the returned function has not been
	// called already
	next := c.wrapMiddlewares(middlewares[1:], last)
	return func(ctx context.Context) {
		// Call the middleware and the following if it has not been done already
		if !called {
			called = true
			ctx = middlewares[0](ctx, next)

			// If next has already been called in middleware, it should not be
			// executed again
			next(ctx)
		}
	}
}

func (c ClientController) executeMiddlewares(ctx context.Context, callback func(ctx context.Context)) {
	// Wrap middleware to have 'next' function when calling them
	wrapped := c.wrapMiddlewares(c.middlewares, callback)

	// Execute wrapped middlewares
	wrapped(ctx)
}

func addClientContextValues(ctx context.Context, path, operation string) context.Context {
	ctx = context.WithValue(ctx, apiContext.KeyIsModule, "asyncapi")
	ctx = context.WithValue(ctx, apiContext.KeyIsProvider, "client")
	ctx = context.WithValue(ctx, apiContext.KeyIsChannel, path)
	return context.WithValue(ctx, apiContext.KeyIsOperation, operation)
}

// Close will clean up any existing resources on the controller
func (c *ClientController) Close(ctx context.Context) {
	// Unsubscribing remaining channels
	c.logger.Info(ctx, "Closing Client controller")
	c.UnsubscribeAll(ctx)
}

// SubscribeAll will subscribe to channels without parameters on which the app is expecting messages.
// For channels with parameters, they should be subscribed independently.
func (c *ClientController) SubscribeAll(ctx context.Context, as ClientSubscriber) error {
	if as == nil {
		return ErrNilClientSubscriber
	}

	if err := c.SubscribeStatus(ctx, as.Status); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (c *ClientController) UnsubscribeAll(ctx context.Context) {
	// Unsubscribe channels with no parameters (if any)
	c.UnsubscribeStatus(ctx)

	// Unsubscribe remaining channels
	for n, stopChan := range c.stopSubscribers {
		stopChan <- true
		delete(c.stopSubscribers, n)
	}
}

// SubscribeStatus will subscribe to new messages from '/status' channel.
//
// Callback function 'fn' will be called each time a new message is received.
// The 'done' argument indicates when the subscription is canceled and can be
// used to clean up resources.
func (c *ClientController) SubscribeStatus(ctx context.Context, fn func(ctx context.Context, msg StatusMessage, done bool)) error {
	// Get channel path
	path := "/status"

	// Set context
	ctx = addClientContextValues(ctx, path, "subscribe")
	ctx = context.WithValue(ctx, apiContext.KeyIsDirection, "reception")

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

			// Process message
			msg, err := newStatusMessageFromUniversalMessage(um)
			if err != nil {
				ctx = context.WithValue(ctx, apiContext.KeyIsMessage, um)
				c.logger.Error(ctx, err.Error())
			}
			ctx = context.WithValue(ctx, apiContext.KeyIsMessage, msg)

			// Send info if message is correct or susbcription is closed
			if err == nil || !open {
				c.logger.Info(ctx, "Received new message")

				// Execute middlewares with the callback
				c.executeMiddlewares(ctx, func(ctx context.Context) {
					fn(ctx, msg, !open)
				})
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

// UnsubscribeStatus will unsubscribe messages from '/status' channel
func (c *ClientController) UnsubscribeStatus(ctx context.Context) {
	// Get channel path
	path := "/status"

	// Set context
	ctx = addClientContextValues(ctx, path, "unsubscribe")

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

// PublishChat will publish messages to '/chat' channel
func (c *ClientController) PublishChat(ctx context.Context, msg ChatMessage) error {
	// Get channel path
	path := "/chat"

	// Set context
	ctx = addClientContextValues(ctx, path, "publish")
	ctx = context.WithValue(ctx, apiContext.KeyIsMessage, msg)
	ctx = context.WithValue(ctx, apiContext.KeyIsDirection, "publication")

	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Publish the message in middlewares
	c.executeMiddlewares(ctx, func(ctx context.Context) {
		// Publish on event broker
		c.logger.Info(ctx, "Publishing to channel")
		err = c.brokerController.Publish(ctx, path, um)
	})

	// Return error from publication on broker
	return err
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

// ChatMessage is the message expected for 'Chat' channel
type ChatMessage struct {
	// Payload will be inserted in the message payload
	Payload string
}

func NewChatMessage() ChatMessage {
	var msg ChatMessage

	return msg
}

// newChatMessageFromUniversalMessage will fill a new ChatMessage with data from UniversalMessage
func newChatMessageFromUniversalMessage(um UniversalMessage) (ChatMessage, error) {
	var msg ChatMessage

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(um.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toUniversalMessage will generate an UniversalMessage from ChatMessage data
func (msg ChatMessage) toUniversalMessage() (UniversalMessage, error) {
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

// Message is the message expected for ” channel
type Message struct {
	// Payload will be inserted in the message payload
	Payload string
}

func NewMessage() Message {
	var msg Message

	return msg
}

// newMessageFromUniversalMessage will fill a new Message with data from UniversalMessage
func newMessageFromUniversalMessage(um UniversalMessage) (Message, error) {
	var msg Message

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(um.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toUniversalMessage will generate an UniversalMessage from Message data
func (msg Message) toUniversalMessage() (UniversalMessage, error) {
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

// StatusMessage is the message expected for 'Status' channel
type StatusMessage struct {
	// Payload will be inserted in the message payload
	Payload string
}

func NewStatusMessage() StatusMessage {
	var msg StatusMessage

	return msg
}

// newStatusMessageFromUniversalMessage will fill a new StatusMessage with data from UniversalMessage
func newStatusMessageFromUniversalMessage(um UniversalMessage) (StatusMessage, error) {
	var msg StatusMessage

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(um.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toUniversalMessage will generate an UniversalMessage from StatusMessage data
func (msg StatusMessage) toUniversalMessage() (UniversalMessage, error) {
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
