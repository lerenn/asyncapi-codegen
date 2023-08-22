// Package "simple" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package simple

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	aapiContext "github.com/lerenn/asyncapi-codegen/pkg/context"
	"github.com/lerenn/asyncapi-codegen/pkg/log"
)

// AppController is the structure that provides publishing capabilities to the
// developer and and connect the broker with the App
type AppController struct {
	brokerController BrokerController
	stopSubscribers  map[string]chan interface{}
	logger           log.Logger
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
	}, nil
}

// SetLogger attaches a logger that will log operations on controller
func (c *AppController) SetLogger(logger log.Logger) {
	c.logger = logger
	c.brokerController.SetLogger(logger)
}

func addAppContextValues(ctx context.Context, path, operation string) context.Context {
	ctx = context.WithValue(ctx, aapiContext.KeyIsModule, "asyncapi")
	ctx = context.WithValue(ctx, aapiContext.KeyIsProvider, "app")
	ctx = context.WithValue(ctx, aapiContext.KeyIsAction, path)
	return context.WithValue(ctx, aapiContext.KeyIsOperation, operation)
}

// Close will clean up any existing resources on the controller
func (c *AppController) Close(ctx context.Context) {
	// Unsubscribing remaining channels
}

// PublishUserSignedup will publish messages to 'user/signedup' channel
func (c *AppController) PublishUserSignedup(ctx context.Context, msg UserSignedUpMessage) error {
	// Get channel path
	path := "user/signedup"

	// Set context
	ctx = addAppContextValues(ctx, path, "publish")
	ctx = context.WithValue(ctx, aapiContext.KeyIsMessage, msg)

	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Publish on event broker
	c.logger.Info(ctx, "Publishing to channel")
	return c.brokerController.Publish(ctx, path, um)
}

// ClientSubscriber represents all handlers that are expecting messages for Client
type ClientSubscriber interface {
	// UserSignedup
	UserSignedup(ctx context.Context, msg UserSignedUpMessage, done bool)
}

// ClientController is the structure that provides publishing capabilities to the
// developer and and connect the broker with the Client
type ClientController struct {
	brokerController BrokerController
	stopSubscribers  map[string]chan interface{}
	logger           log.Logger
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
	}, nil
}

// SetLogger attaches a logger that will log operations on controller
func (c *ClientController) SetLogger(logger log.Logger) {
	c.logger = logger
	c.brokerController.SetLogger(logger)
}

func addClientContextValues(ctx context.Context, path, operation string) context.Context {
	ctx = context.WithValue(ctx, aapiContext.KeyIsModule, "asyncapi")
	ctx = context.WithValue(ctx, aapiContext.KeyIsProvider, "client")
	ctx = context.WithValue(ctx, aapiContext.KeyIsAction, path)
	return context.WithValue(ctx, aapiContext.KeyIsOperation, operation)
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

	if err := c.SubscribeUserSignedup(ctx, as.UserSignedup); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (c *ClientController) UnsubscribeAll(ctx context.Context) {
	// Unsubscribe channels with no parameters (if any)
	c.UnsubscribeUserSignedup(ctx)

	// Unsubscribe remaining channels
	for n, stopChan := range c.stopSubscribers {
		stopChan <- true
		delete(c.stopSubscribers, n)
	}
}

// SubscribeUserSignedup will subscribe to new messages from 'user/signedup' channel.
//
// Callback function 'fn' will be called each time a new message is received.
// The 'done' argument indicates when the subscription is canceled and can be
// used to clean up resources.
func (c *ClientController) SubscribeUserSignedup(ctx context.Context, fn func(ctx context.Context, msg UserSignedUpMessage, done bool)) error {
	// Get channel path
	path := "user/signedup"

	// Set context
	ctx = addClientContextValues(ctx, path, "subscribe")

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
			msg, err := newUserSignedUpMessageFromUniversalMessage(um)
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

// UnsubscribeUserSignedup will unsubscribe messages from 'user/signedup' channel
func (c *ClientController) UnsubscribeUserSignedup(ctx context.Context) {
	// Get channel path
	path := "user/signedup"

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
	SetLogger(logger log.Logger)

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

// UserSignedUpMessage is the message expected for 'UserSignedUp' channel
type UserSignedUpMessage struct {
	// Payload will be inserted in the message payload
	Payload struct {
		// Description: Name of the user
		DisplayName *string `json:"display_name"`

		// Description: Email of the user
		Email *string `json:"email"`
	}
}

func NewUserSignedUpMessage() UserSignedUpMessage {
	var msg UserSignedUpMessage

	return msg
}

// newUserSignedUpMessageFromUniversalMessage will fill a new UserSignedUpMessage with data from UniversalMessage
func newUserSignedUpMessageFromUniversalMessage(um UniversalMessage) (UserSignedUpMessage, error) {
	var msg UserSignedUpMessage

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(um.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toUniversalMessage will generate an UniversalMessage from UserSignedUpMessage data
func (msg UserSignedUpMessage) toUniversalMessage() (UniversalMessage, error) {
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
