// Package "main" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"

	"github.com/google/uuid"
)

// UserSubscriber represents all handlers that are expecting messages for User
type UserSubscriber interface {
	// Pong subscribes to messages placed on the 'pong' channel
	Pong(ctx context.Context, msg PongMessage)
}

// UserController is the structure that provides publishing capabilities to the
// developer and and connect the broker with the User
type UserController struct {
	controller
}

// NewUserController links the User to the broker
func NewUserController(bc extensions.BrokerController, options ...ControllerOption) (*UserController, error) {
	// Check if broker controller has been provided
	if bc == nil {
		return nil, extensions.ErrNilBrokerController
	}

	// Create default controller
	controller := controller{
		broker:        bc,
		subscriptions: make(map[string]extensions.BrokerChannelSubscription),
		logger:        extensions.DummyLogger{},
		middlewares:   make([]extensions.Middleware, 0),
	}

	// Apply options
	for _, option := range options {
		option(&controller)
	}

	return &UserController{controller: controller}, nil
}

func (c UserController) wrapMiddlewares(middlewares []extensions.Middleware, last extensions.NextMiddleware) func(ctx context.Context) {
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

func (c UserController) executeMiddlewares(ctx context.Context, callback func(ctx context.Context)) {
	// Wrap middleware to have 'next' function when calling them
	wrapped := c.wrapMiddlewares(c.middlewares, callback)

	// Execute wrapped middlewares
	wrapped(ctx)
}

func addUserContextValues(ctx context.Context, path string) context.Context {
	ctx = context.WithValue(ctx, extensions.ContextKeyIsVersion, "1.0.0")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsProvider, "user")
	return context.WithValue(ctx, extensions.ContextKeyIsChannel, path)
}

// Close will clean up any existing resources on the controller
func (c *UserController) Close(ctx context.Context) {
	// Unsubscribing remaining channels
	c.UnsubscribeAll(ctx)

	c.logger.Info(ctx, "Closed user controller")
}

// SubscribeAll will subscribe to channels without parameters on which the app is expecting messages.
// For channels with parameters, they should be subscribed independently.
func (c *UserController) SubscribeAll(ctx context.Context, as UserSubscriber) error {
	if as == nil {
		return extensions.ErrNilUserSubscriber
	}

	if err := c.SubscribePong(ctx, as.Pong); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (c *UserController) UnsubscribeAll(ctx context.Context) {
	c.UnsubscribePong(ctx)
}

// SubscribePong will subscribe to new messages from 'pong' channel.
//
// Callback function 'fn' will be called each time a new message is received.
func (c *UserController) SubscribePong(ctx context.Context, fn func(ctx context.Context, msg PongMessage)) error {
	// Get channel path
	path := "pong"

	// Set context
	ctx = addUserContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsMessageDirection, "reception")

	// Check if there is already a subscription
	_, exists := c.subscriptions[path]
	if exists {
		err := fmt.Errorf("%w: %q channel is already subscribed", extensions.ErrAlreadySubscribedChannel, path)
		c.logger.Error(ctx, err.Error())
		return err
	}

	// Subscribe to broker channel
	sub, err := c.broker.Subscribe(ctx, path)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return err
	}
	c.logger.Info(ctx, "Subscribed to channel")

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for {
			// Wait for next message
			bMsg, open := <-sub.MessagesChannel()

			// If subscription is closed and there is no more message
			// (i.e. uninitialized message), then exit the function
			if !open && bMsg.IsUninitialized() {
				return
			}

			// Set broker message to context
			ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, bMsg)

			// Process message
			msg, err := newPongMessageFromBrokerMessage(bMsg)
			if err != nil {
				c.logger.Error(ctx, err.Error())
			}
			msgCtx := context.WithValue(ctx, extensions.ContextKeyIsMessage, msg)

			// Add correlation ID to context if it exists
			if id := msg.CorrelationID(); id != "" {
				ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, id)
			}

			// Execute middlewares with the callback
			c.executeMiddlewares(msgCtx, func(ctx context.Context) {
				fn(ctx, msg)
			})
		}
	}()

	// Add the cancel channel to the inside map
	c.subscriptions[path] = sub

	return nil
}

// UnsubscribePong will unsubscribe messages from 'pong' channel
func (c *UserController) UnsubscribePong(ctx context.Context) {
	// Get channel path
	path := "pong"

	// Check if there subscribers for this channel
	sub, exists := c.subscriptions[path]
	if !exists {
		return
	}

	// Set context
	ctx = addUserContextValues(ctx, path)

	// Stop the subscription
	sub.Cancel()

	// Remove if from the subscribers
	delete(c.subscriptions, path)

	c.logger.Info(ctx, "Unsubscribed from channel")
}

// PublishPing will publish messages to 'ping' channel
func (c *UserController) PublishPing(ctx context.Context, msg PingMessage) error {
	// Get channel path
	path := "ping"

	// Set correlation ID if it does not exist
	if id := msg.CorrelationID(); id == "" {
		msg.SetCorrelationID(uuid.New().String())
	}

	// Set context
	ctx = addUserContextValues(ctx, path)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsMessage, msg)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsMessageDirection, "publication")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, msg.CorrelationID())

	// Convert to BrokerMessage
	bMsg, err := msg.toBrokerMessage()
	if err != nil {
		return err
	}

	// Set broker message to context
	ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, bMsg)

	// Publish the message on event-broker through middlewares
	c.executeMiddlewares(ctx, func(ctx context.Context) {
		err = c.broker.Publish(ctx, path, bMsg)
	})

	// Return error from publication on broker
	return err
}

// WaitForPong will wait for a specific message by its correlation ID
//
// The pub function is the publication function that should be used to send the message
// It will be called after subscribing to the channel to avoid race condition, and potentially loose the message
func (cc *UserController) WaitForPong(ctx context.Context, publishMsg MessageWithCorrelationID, pub func(ctx context.Context) error) (PongMessage, error) {
	// Get channel path
	path := "pong"

	// Set context
	ctx = addUserContextValues(ctx, path)

	// Subscribe to broker channel
	sub, err := cc.broker.Subscribe(ctx, path)
	if err != nil {
		cc.logger.Error(ctx, err.Error())
		return PongMessage{}, err
	}
	cc.logger.Info(ctx, "Subscribed to channel")

	// Close subscriber on leave
	defer func() {
		// Stop the subscription
		sub.Cancel()

		// Logging unsubscribing
		cc.logger.Info(ctx, "Unsubscribed from channel")
	}()

	// Execute callback for publication
	if err = pub(ctx); err != nil {
		return PongMessage{}, err
	}

	// Wait for corresponding response
	for {
		select {
		case bMsg, open := <-sub.MessagesChannel():
			// If subscription is closed and there is no more message
			// (i.e. uninitialized message), then the subscription ended before
			// receiving the expected message
			if !open && bMsg.IsUninitialized() {
				cc.logger.Error(ctx, "Channel closed before getting message")
				return PongMessage{}, extensions.ErrSubscriptionCanceled
			}

			// Get new message
			msg, err := newPongMessageFromBrokerMessage(bMsg)
			if err != nil {
				cc.logger.Error(ctx, err.Error())
			}

			// If message doesn't have corresponding correlation ID, then continue
			if publishMsg.CorrelationID() != msg.CorrelationID() {
				continue
			}

			// Set context with received values as it is the expected message
			msgCtx := context.WithValue(ctx, extensions.ContextKeyIsMessage, msg)
			msgCtx = context.WithValue(msgCtx, extensions.ContextKeyIsBrokerMessage, bMsg)
			msgCtx = context.WithValue(msgCtx, extensions.ContextKeyIsMessageDirection, "reception")
			msgCtx = context.WithValue(msgCtx, extensions.ContextKeyIsCorrelationID, publishMsg.CorrelationID())

			// Execute middlewares before returning
			cc.executeMiddlewares(msgCtx, func(_ context.Context) {
				/* Nothing to do more */
			})

			// return the message to the caller
			return msg, nil
		case <-ctx.Done(): // Set corrsponding error if context is done
			cc.logger.Error(ctx, "Context done before getting message")
			return PongMessage{}, extensions.ErrContextCanceled
		}
	}
}

// controller is the controller that will be used to communicate with the broker
// It will be used internally by AppController and UserController
type controller struct {
	// broker is the broker controller that will be used to communicate
	broker extensions.BrokerController
	// subscriptions is a map of all subscriptions
	subscriptions map[string]extensions.BrokerChannelSubscription
	// logger is the logger that will be used² to log operations on controller
	logger extensions.Logger
	// middlewares are the middlewares that will be executed when sending or
	// receiving messages
	middlewares []extensions.Middleware
}

// ControllerOption is the type of the options that can be passed
// when creating a new Controller
type ControllerOption func(controller *controller)

// WithLogger attaches a logger to the controller
func WithLogger(logger extensions.Logger) ControllerOption {
	return func(controller *controller) {
		controller.logger = logger
	}
}

// WithMiddlewares attaches middlewares that will be executed when sending or receiving messages
func WithMiddlewares(middlewares ...extensions.Middleware) ControllerOption {
	return func(controller *controller) {
		controller.middlewares = middlewares
	}
}

type MessageWithCorrelationID interface {
	CorrelationID() string
	SetCorrelationID(id string)
}

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
		// Description: Correlation ID set by user
		CorrelationID *string `json:"correlation_id"`
	}

	// Payload will be inserted in the message payload
	Payload string
}

func NewPingMessage() PingMessage {
	var msg PingMessage

	// Set correlation ID
	u := uuid.New().String()
	msg.Headers.CorrelationID = &u

	return msg
}

// newPingMessageFromBrokerMessage will fill a new PingMessage with data from generic broker message
func newPingMessageFromBrokerMessage(bMsg extensions.BrokerMessage) (PingMessage, error) {
	var msg PingMessage

	// Convert to string
	msg.Payload = string(bMsg.Payload)

	// Get each headers from broker message
	for k, v := range bMsg.Headers {
		switch {
		case k == "correlationId": // Retrieving CorrelationID header
			h := string(v)
			msg.Headers.CorrelationID = &h
		default:
			// TODO: log unknown error
		}
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toBrokerMessage will generate a generic broker message from PingMessage data
func (msg PingMessage) toBrokerMessage() (extensions.BrokerMessage, error) {
	// TODO: implement checks on message

	// Convert to []byte
	payload := []byte(msg.Payload)

	// Add each headers to broker message
	headers := make(map[string][]byte, 1)

	// Adding CorrelationID header
	if msg.Headers.CorrelationID != nil {
		headers["correlationId"] = []byte(*msg.Headers.CorrelationID)
	}

	return extensions.BrokerMessage{
		Headers: headers,
		Payload: payload,
	}, nil
}

// CorrelationID will give the correlation ID of the message, based on AsyncAPI spec
func (msg PingMessage) CorrelationID() string {
	if msg.Headers.CorrelationID != nil {
		return *msg.Headers.CorrelationID
	}

	return ""
}

// SetCorrelationID will set the correlation ID of the message, based on AsyncAPI spec
func (msg *PingMessage) SetCorrelationID(id string) {
	msg.Headers.CorrelationID = &id
}

// SetAsResponseFrom will correlate the message with the one passed in parameter.
// It will assign the 'req' message correlation ID to the message correlation ID,
// both specified in AsyncAPI spec.
func (msg *PingMessage) SetAsResponseFrom(req MessageWithCorrelationID) {
	id := req.CorrelationID()
	msg.Headers.CorrelationID = &id
}

// PongMessage is the message expected for 'Pong' channel
type PongMessage struct {
	// Headers will be used to fill the message headers
	Headers struct {
		// Description: Correlation ID set by user on corresponding request
		CorrelationID *string `json:"correlation_id"`
	}

	// Payload will be inserted in the message payload
	Payload struct {
		// Description: Pong message
		Message string `json:"message"`

		// Description: Pong creation time
		Time time.Time `json:"time"`
	}
}

func NewPongMessage() PongMessage {
	var msg PongMessage

	// Set correlation ID
	u := uuid.New().String()
	msg.Headers.CorrelationID = &u

	return msg
}

// newPongMessageFromBrokerMessage will fill a new PongMessage with data from generic broker message
func newPongMessageFromBrokerMessage(bMsg extensions.BrokerMessage) (PongMessage, error) {
	var msg PongMessage

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(bMsg.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// Get each headers from broker message
	for k, v := range bMsg.Headers {
		switch {
		case k == "correlationId": // Retrieving CorrelationID header
			h := string(v)
			msg.Headers.CorrelationID = &h
		default:
			// TODO: log unknown error
		}
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toBrokerMessage will generate a generic broker message from PongMessage data
func (msg PongMessage) toBrokerMessage() (extensions.BrokerMessage, error) {
	// TODO: implement checks on message

	// Marshal payload to JSON
	payload, err := json.Marshal(msg.Payload)
	if err != nil {
		return extensions.BrokerMessage{}, err
	}

	// Add each headers to broker message
	headers := make(map[string][]byte, 1)

	// Adding CorrelationID header
	if msg.Headers.CorrelationID != nil {
		headers["correlationId"] = []byte(*msg.Headers.CorrelationID)
	}

	return extensions.BrokerMessage{
		Headers: headers,
		Payload: payload,
	}, nil
}

// CorrelationID will give the correlation ID of the message, based on AsyncAPI spec
func (msg PongMessage) CorrelationID() string {
	if msg.Headers.CorrelationID != nil {
		return *msg.Headers.CorrelationID
	}

	return ""
}

// SetCorrelationID will set the correlation ID of the message, based on AsyncAPI spec
func (msg *PongMessage) SetCorrelationID(id string) {
	msg.Headers.CorrelationID = &id
}

// SetAsResponseFrom will correlate the message with the one passed in parameter.
// It will assign the 'req' message correlation ID to the message correlation ID,
// both specified in AsyncAPI spec.
func (msg *PongMessage) SetAsResponseFrom(req MessageWithCorrelationID) {
	id := req.CorrelationID()
	msg.Headers.CorrelationID = &id
}
