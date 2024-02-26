// Package "requestreply" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package requestreply

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"

	"github.com/google/uuid"
)

// AppController is the structure that provides sending capabilities to the
// developer and and connect the broker with the App
type AppController struct {
	controller
}

// NewAppController links the App to the broker
func NewAppController(bc extensions.BrokerController, options ...ControllerOption) (*AppController, error) {
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

	return &AppController{controller: controller}, nil
}

func (c AppController) wrapMiddlewares(
	middlewares []extensions.Middleware,
	callback extensions.NextMiddleware,
) func(ctx context.Context, msg *extensions.BrokerMessage) error {
	var called bool

	// If there is no more middleware
	if len(middlewares) == 0 {
		return func(ctx context.Context, msg *extensions.BrokerMessage) error {
			// Call the callback if it exists and it has not been called already
			if callback != nil && !called {
				called = true
				return callback(ctx)
			}

			// Nil can be returned, as the callback has already been called
			return nil
		}
	}

	// Get the next function to call from next middlewares or callback
	next := c.wrapMiddlewares(middlewares[1:], callback)

	// Wrap middleware into a check function that will call execute the middleware
	// and call the next wrapped middleware if the returned function has not been
	// called already
	return func(ctx context.Context, msg *extensions.BrokerMessage) error {
		// Call the middleware and the following if it has not been done already
		if !called {
			// Create the next call with the context and the message
			nextWithArgs := func(ctx context.Context) error {
				return next(ctx, msg)
			}

			// Call the middleware and register it as already called
			called = true
			if err := middlewares[0](ctx, msg, nextWithArgs); err != nil {
				return err
			}

			// If next has already been called in middleware, it should not be executed again
			return nextWithArgs(ctx)
		}

		// Nil can be returned, as the next middleware has already been called
		return nil
	}
}

func (c AppController) executeMiddlewares(ctx context.Context, msg *extensions.BrokerMessage, callback extensions.NextMiddleware) error {
	// Wrap middleware to have 'next' function when calling them
	wrapped := c.wrapMiddlewares(c.middlewares, callback)

	// Execute wrapped middlewares
	return wrapped(ctx, msg)
}

func addAppContextValues(ctx context.Context, addr string) context.Context {
	ctx = context.WithValue(ctx, extensions.ContextKeyIsVersion, "1.0.0")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsProvider, "app")
	return context.WithValue(ctx, extensions.ContextKeyIsChannel, addr)
}

// Close will clean up any existing resources on the controller
func (c *AppController) Close(ctx context.Context) {
	// Unsubscribing remaining channels
}

// ListenPingRequest will receive 'Ping' messages from 'issue130.ping' channel
//
// Callback function 'fn' will be called each time a new message is received.
//
// NOTE: for now, this only support the first message from AsyncAPI list.
//
// NOTE: for now, this only support the first message from AsyncAPI list.
// If you need support for other messages, please raise an issue.
func (c *AppController) ListenPingRequest(ctx context.Context, fn func(ctx context.Context, msg Ping)) error {
	// Get channel address
	addr := "issue130.ping"

	// Set context
	ctx = addAppContextValues(ctx, addr)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "reception")

	// Check if the controller is already listening
	_, exists := c.subscriptions[addr]
	if exists {
		err := fmt.Errorf("%w: controller is already listening on channel %q", extensions.ErrAlreadySubscribedChannel, addr)
		c.logger.Error(ctx, err.Error())
		return err
	}

	// Subscribe to broker channel
	sub, err := c.broker.Subscribe(ctx, addr)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return err
	}
	c.logger.Info(ctx, "Listening to channel")

	// Asynchronously listen to new messages and pass them to app receiver
	go func() {
		for {
			// Wait for next message
			brokerMsg, open := <-sub.MessagesChannel()

			// If subscription is closed and there is no more message
			// (i.e. uninitialized message), then exit the function
			if !open && brokerMsg.IsUninitialized() {
				return
			}

			// Set broker message to context
			ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

			// Execute middlewares before handling the message
			if err := c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
				// Process message
				msg, err := newPingFromBrokerMessage(brokerMsg)
				if err != nil {
					return err
				}

				// Execute the subscription function
				fn(ctx, msg)

				return nil
			}); err != nil {
				c.logger.Error(ctx, err.Error())
			}
		}
	}()

	// Add the cancel channel to the inside map
	c.subscriptions[addr] = sub

	return nil
}

// UnlistenPingRequest will stop the reception of messages from 'issue130.ping' channel.
// A timeout can be set in context to avoid blocking operation, if needed.
func (c *AppController) UnlistenPingRequest(ctx context.Context) {
	// Get channel address
	addr := "issue130.ping"

	// Check if there receivers for this channel
	sub, exists := c.subscriptions[addr]
	if !exists {
		return
	}

	// Set context
	ctx = addAppContextValues(ctx, addr)

	// Stop the subscription
	sub.Cancel(ctx)

	// Remove if from the receivers
	delete(c.subscriptions, addr)

	c.logger.Info(ctx, "Unlistend from channel")
} // ListenPingRequestWithCorrelationID will receive 'PingWithCorrelationID' messages from 'issue130.pingWithCorrelationID' channel
// Callback function 'fn' will be called each time a new message is received.
//
// NOTE: for now, this only support the first message from AsyncAPI list.
//
// NOTE: for now, this only support the first message from AsyncAPI list.
// If you need support for other messages, please raise an issue.
func (c *AppController) ListenPingRequestWithCorrelationID(ctx context.Context, fn func(ctx context.Context, msg PingWithCorrelationID)) error {
	// Get channel address
	addr := "issue130.pingWithCorrelationID"

	// Set context
	ctx = addAppContextValues(ctx, addr)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "reception")

	// Check if the controller is already listening
	_, exists := c.subscriptions[addr]
	if exists {
		err := fmt.Errorf("%w: controller is already listening on channel %q", extensions.ErrAlreadySubscribedChannel, addr)
		c.logger.Error(ctx, err.Error())
		return err
	}

	// Subscribe to broker channel
	sub, err := c.broker.Subscribe(ctx, addr)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return err
	}
	c.logger.Info(ctx, "Listening to channel")

	// Asynchronously listen to new messages and pass them to app receiver
	go func() {
		for {
			// Wait for next message
			brokerMsg, open := <-sub.MessagesChannel()

			// If subscription is closed and there is no more message
			// (i.e. uninitialized message), then exit the function
			if !open && brokerMsg.IsUninitialized() {
				return
			}

			// Set broker message to context
			ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

			// Execute middlewares before handling the message
			if err := c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
				// Process message
				msg, err := newPingWithCorrelationIDFromBrokerMessage(brokerMsg)
				if err != nil {
					return err
				}

				// Add correlation ID to context if it exists
				if id := msg.CorrelationID(); id != "" {
					ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, id)
				}

				// Execute the subscription function
				fn(ctx, msg)

				return nil
			}); err != nil {
				c.logger.Error(ctx, err.Error())
			}
		}
	}()

	// Add the cancel channel to the inside map
	c.subscriptions[addr] = sub

	return nil
}

// UnlistenPingRequestWithCorrelationID will stop the reception of messages from 'issue130.pingWithCorrelationID' channel.
// A timeout can be set in context to avoid blocking operation, if needed.
func (c *AppController) UnlistenPingRequestWithCorrelationID(ctx context.Context) {
	// Get channel address
	addr := "issue130.pingWithCorrelationID"

	// Check if there receivers for this channel
	sub, exists := c.subscriptions[addr]
	if !exists {
		return
	}

	// Set context
	ctx = addAppContextValues(ctx, addr)

	// Stop the subscription
	sub.Cancel(ctx)

	// Remove if from the receivers
	delete(c.subscriptions, addr)

	c.logger.Info(ctx, "Unlistend from channel")
}

// ReplyToPingRequest should be used to reply to PingRequest messages by sending
// 'Pong' messages to 'issue130.pong' channel.
//
// NOTE: for now, this only support the first message from AsyncAPI list.
// If you need support for other messages, please raise an issue.
func (c *AppController) ReplyToPingRequest(ctx context.Context, msg Pong) error {
	// Get channel address
	addr := "issue130.pong"

	// Set context
	ctx = addAppContextValues(ctx, addr)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "publication")

	// Convert to BrokerMessage
	brokerMsg, err := msg.toBrokerMessage()
	if err != nil {
		return err
	}

	// Set broker message to context
	ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

	// Send the message on event-broker through middlewares
	return c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
		return c.broker.Publish(ctx, addr, brokerMsg)
	})
}

// ReplyToPingRequestWithCorrelationID should be used to reply to PingRequestWithCorrelationID messages by sending
// 'PongWithCorrelationID' messages to 'issue130.pongWithCorrelationID' channel.
//
// NOTE: for now, this only support the first message from AsyncAPI list.
// If you need support for other messages, please raise an issue.
func (c *AppController) ReplyToPingRequestWithCorrelationID(ctx context.Context, msg PongWithCorrelationID) error {
	// Get channel address
	addr := "issue130.pongWithCorrelationID"

	// Set correlation ID if it does not exist
	if id := msg.CorrelationID(); id == "" {
		c.logger.Error(ctx, extensions.ErrNoCorrelationIDSet.Error())
		return extensions.ErrNoCorrelationIDSet

	}

	// Set context
	ctx = addAppContextValues(ctx, addr)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "publication")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, msg.CorrelationID())

	// Convert to BrokerMessage
	brokerMsg, err := msg.toBrokerMessage()
	if err != nil {
		return err
	}

	// Set broker message to context
	ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

	// Send the message on event-broker through middlewares
	return c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
		return c.broker.Publish(ctx, addr, brokerMsg)
	})
}

// UserListener represents all handlers that are listening messages for User
type UserListener interface {
	// PingRequest listen to messages placed on the 'pingRequest' channel
	PingRequest(ctx context.Context, msg Ping)

	// PingRequestWithCorrelationID listen to messages placed on the 'pingRequestWithCorrelationID' channel
	PingRequestWithCorrelationID(ctx context.Context, msg PingWithCorrelationID)
}

// UserController is the structure that provides sending capabilities to the
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

func (c UserController) wrapMiddlewares(
	middlewares []extensions.Middleware,
	callback extensions.NextMiddleware,
) func(ctx context.Context, msg *extensions.BrokerMessage) error {
	var called bool

	// If there is no more middleware
	if len(middlewares) == 0 {
		return func(ctx context.Context, msg *extensions.BrokerMessage) error {
			// Call the callback if it exists and it has not been called already
			if callback != nil && !called {
				called = true
				return callback(ctx)
			}

			// Nil can be returned, as the callback has already been called
			return nil
		}
	}

	// Get the next function to call from next middlewares or callback
	next := c.wrapMiddlewares(middlewares[1:], callback)

	// Wrap middleware into a check function that will call execute the middleware
	// and call the next wrapped middleware if the returned function has not been
	// called already
	return func(ctx context.Context, msg *extensions.BrokerMessage) error {
		// Call the middleware and the following if it has not been done already
		if !called {
			// Create the next call with the context and the message
			nextWithArgs := func(ctx context.Context) error {
				return next(ctx, msg)
			}

			// Call the middleware and register it as already called
			called = true
			if err := middlewares[0](ctx, msg, nextWithArgs); err != nil {
				return err
			}

			// If next has already been called in middleware, it should not be executed again
			return nextWithArgs(ctx)
		}

		// Nil can be returned, as the next middleware has already been called
		return nil
	}
}

func (c UserController) executeMiddlewares(ctx context.Context, msg *extensions.BrokerMessage, callback extensions.NextMiddleware) error {
	// Wrap middleware to have 'next' function when calling them
	wrapped := c.wrapMiddlewares(c.middlewares, callback)

	// Execute wrapped middlewares
	return wrapped(ctx, msg)
}

func addUserContextValues(ctx context.Context, addr string) context.Context {
	ctx = context.WithValue(ctx, extensions.ContextKeyIsVersion, "1.0.0")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsProvider, "user")
	return context.WithValue(ctx, extensions.ContextKeyIsChannel, addr)
}

// Close will clean up any existing resources on the controller
func (c *UserController) Close(ctx context.Context) {
	// Unsubscribing remaining channels
	c.UnlistenAll(ctx)

	c.logger.Info(ctx, "Closed user controller")
}

// ListenAll will receive to channels without parameters on which the app is expecting messages.
// For channels with parameters, they should be received independently.
func (c *UserController) ListenAll(ctx context.Context, as UserListener) error {
	if as == nil {
		return extensions.ErrNilUserListener
	}

	return nil
}

// UnlistenAll will stop the listening of all remaining listening channels
func (c *UserController) UnlistenAll(ctx context.Context) {
}

// AsyncPingRequest will send 'Ping' messages to 'issue130.ping' channel.
// NOTE: this won't wait for reply, use the normal version to get the reply or do the catching reply manually.
//
// NOTE: for now, this only support the first message from AsyncAPI list.
// If you need support for other messages, please raise an issue.
func (c *UserController) AsyncPingRequest(ctx context.Context, msg Ping) error {
	// Get channel address
	addr := "issue130.ping"

	// Set context
	ctx = addUserContextValues(ctx, addr)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "publication")

	// Convert to BrokerMessage
	brokerMsg, err := msg.toBrokerMessage()
	if err != nil {
		return err
	}

	// Set broker message to context
	ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

	// Send the message on event-broker through middlewares
	return c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
		return c.broker.Publish(ctx, addr, brokerMsg)
	})
}

// PingRequest will send a message and wait for the reply message
// on channel 'issue130.pong'.
//
// If a correlation ID is set in the AsyncAPI, then this will wait for the
// reply with the same correlation ID. Otherwise, it will returns the first
// message on the reply channel.
//
// A timeout can be set in context to avoid blocking operation, if needed.
func (c *UserController) PingRequest(ctx context.Context, msg Ping) (Pong, error) {
	// Get receiving channel address
	addr := "issue130.pong"

	// Set context
	ctx = addUserContextValues(ctx, addr)

	// Listening to broker channel
	sub, err := c.broker.Subscribe(ctx, addr)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return Pong{}, err
	}
	c.logger.Info(ctx, "Listening to channel")

	// Close receiver on leave
	defer func() {
		// Stop the subscription
		sub.Cancel(ctx)

		// Logging unsubscribing
		c.logger.Info(ctx, "Unlistend from channel")
	}()

	// Send the message
	if err := c.AsyncPingRequest(ctx, msg); err != nil {
		c.logger.Error(ctx, "error happened when sending message", extensions.LogInfo{Key: "error", Value: err.Error()})
		return Pong{}, fmt.Errorf("error happened when sending message: %w", err)
	}

	// Wait for corresponding response
	for {
		select {
		case brokerMsg, open := <-sub.MessagesChannel():
			// If subscription is closed and there is no more message
			// (i.e. uninitialized message), then the subscription ended before
			// receiving the expected message
			if !open && brokerMsg.IsUninitialized() {
				c.logger.Error(ctx, "Channel closed before getting message")
				return Pong{}, extensions.ErrSubscriptionCanceled
			}

			// There is correlation no ID, so it will automatically return at
			// the first received message.

			// Set context with received values as it is the expected message
			msgCtx := context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())
			msgCtx = context.WithValue(msgCtx, extensions.ContextKeyIsDirection, "reception")

			// Execute middlewares before returning
			if err := c.executeMiddlewares(msgCtx, &brokerMsg, nil); err != nil {
				return Pong{}, err
			}

			// Return the message to the caller
			//
			// NOTE: it is transformed from the broker again, as it could have
			// been modified by middlewares
			return newPongFromBrokerMessage(brokerMsg)
		case <-ctx.Done(): // Set corrsponding error if context is done
			c.logger.Error(ctx, "Context done before getting message")
			return Pong{}, extensions.ErrContextCanceled
		}
	}
}

// AsyncPingRequestWithCorrelationID will send 'PingWithCorrelationID' messages to 'issue130.pingWithCorrelationID' channel.
// NOTE: this won't wait for reply, use the normal version to get the reply or do the catching reply manually.
//
// NOTE: for now, this only support the first message from AsyncAPI list.
// If you need support for other messages, please raise an issue.
func (c *UserController) AsyncPingRequestWithCorrelationID(ctx context.Context, msg PingWithCorrelationID) error {
	// Get channel address
	addr := "issue130.pingWithCorrelationID"

	// Set correlation ID if it does not exist
	if id := msg.CorrelationID(); id == "" {
		msg.SetCorrelationID(uuid.New().String())
	}

	// Set context
	ctx = addUserContextValues(ctx, addr)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "publication")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsCorrelationID, msg.CorrelationID())

	// Convert to BrokerMessage
	brokerMsg, err := msg.toBrokerMessage()
	if err != nil {
		return err
	}

	// Set broker message to context
	ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())

	// Send the message on event-broker through middlewares
	return c.executeMiddlewares(ctx, &brokerMsg, func(ctx context.Context) error {
		return c.broker.Publish(ctx, addr, brokerMsg)
	})
}

// PingRequestWithCorrelationID will send a message and wait for the reply message
// on channel 'issue130.pongWithCorrelationID'.
//
// If a correlation ID is set in the AsyncAPI, then this will wait for the
// reply with the same correlation ID. Otherwise, it will returns the first
// message on the reply channel.
//
// A timeout can be set in context to avoid blocking operation, if needed.
func (c *UserController) PingRequestWithCorrelationID(ctx context.Context, msg PingWithCorrelationID) (PongWithCorrelationID, error) {
	// Get receiving channel address
	addr := "issue130.pongWithCorrelationID"

	// Set context
	ctx = addUserContextValues(ctx, addr)

	// Listening to broker channel
	sub, err := c.broker.Subscribe(ctx, addr)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return PongWithCorrelationID{}, err
	}
	c.logger.Info(ctx, "Listening to channel")

	// Close receiver on leave
	defer func() {
		// Stop the subscription
		sub.Cancel(ctx)

		// Logging unsubscribing
		c.logger.Info(ctx, "Unlistend from channel")
	}()

	// Set correlation ID if it does not exist
	if id := msg.CorrelationID(); id == "" {
		msg.SetCorrelationID(uuid.New().String())
	}

	// Send the message
	if err := c.AsyncPingRequestWithCorrelationID(ctx, msg); err != nil {
		c.logger.Error(ctx, "error happened when sending message", extensions.LogInfo{Key: "error", Value: err.Error()})
		return PongWithCorrelationID{}, fmt.Errorf("error happened when sending message: %w", err)
	}

	// Wait for corresponding response
	for {
		select {
		case brokerMsg, open := <-sub.MessagesChannel():
			// If subscription is closed and there is no more message
			// (i.e. uninitialized message), then the subscription ended before
			// receiving the expected message
			if !open && brokerMsg.IsUninitialized() {
				c.logger.Error(ctx, "Channel closed before getting message")
				return PongWithCorrelationID{}, extensions.ErrSubscriptionCanceled
			}

			// Get new message
			rmsg, err := newPongWithCorrelationIDFromBrokerMessage(brokerMsg)
			if err != nil {
				c.logger.Error(ctx, err.Error())
			}

			// If message doesn't have corresponding correlation ID, then ingore and continue
			if msg.CorrelationID() != rmsg.CorrelationID() {
				continue
			}

			// Set context with received values as it is the expected message
			msgCtx := context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, brokerMsg.String())
			msgCtx = context.WithValue(msgCtx, extensions.ContextKeyIsDirection, "reception")
			msgCtx = context.WithValue(msgCtx, extensions.ContextKeyIsCorrelationID, msg.CorrelationID())

			// Execute middlewares before returning
			if err := c.executeMiddlewares(msgCtx, &brokerMsg, nil); err != nil {
				return PongWithCorrelationID{}, err
			}

			// Return the message to the caller
			//
			// NOTE: it is transformed from the broker again, as it could have
			// been modified by middlewares
			return newPongWithCorrelationIDFromBrokerMessage(brokerMsg)
		case <-ctx.Done(): // Set corrsponding error if context is done
			c.logger.Error(ctx, "Context done before getting message")
			return PongWithCorrelationID{}, extensions.ErrContextCanceled
		}
	}
}

// AsyncAPIVersion is the version of the used AsyncAPI document
const AsyncAPIVersion = "1.0.0"

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

// Message 'Ping' reference another one at '#/components/messages/ping'.
// This should be fixed in a future version to allow message override.
// If you encounter this message, feel free to open an issue on this subject
// to let know that you need this functionnality.

// Message 'Ping' reference another one at '#/components/messages/pingWithCorrelationID'.
// This should be fixed in a future version to allow message override.
// If you encounter this message, feel free to open an issue on this subject
// to let know that you need this functionnality.

// Message 'Pong' reference another one at '#/components/messages/pong'.
// This should be fixed in a future version to allow message override.
// If you encounter this message, feel free to open an issue on this subject
// to let know that you need this functionnality.

// Message 'Pong' reference another one at '#/components/messages/pongWithCorrelationID'.
// This should be fixed in a future version to allow message override.
// If you encounter this message, feel free to open an issue on this subject
// to let know that you need this functionnality.

// Ping is the message expected for 'Ping' channel
type Ping struct {
	// Payload will be inserted in the message payload
	Payload struct {
		Event *string `json:"event"`
	}
}

func NewPing() Ping {
	var msg Ping

	return msg
}

// newPingFromBrokerMessage will fill a new Ping with data from generic broker message
func newPingFromBrokerMessage(bMsg extensions.BrokerMessage) (Ping, error) {
	var msg Ping

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(bMsg.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toBrokerMessage will generate a generic broker message from Ping data
func (msg Ping) toBrokerMessage() (extensions.BrokerMessage, error) {
	// TODO: implement checks on message

	// Marshal payload to JSON
	payload, err := json.Marshal(msg.Payload)
	if err != nil {
		return extensions.BrokerMessage{}, err
	}

	// There is no headers here
	headers := make(map[string][]byte, 0)

	return extensions.BrokerMessage{
		Headers: headers,
		Payload: payload,
	}, nil
}

// PingWithCorrelationID is the message expected for 'PingWithCorrelationID' channel
type PingWithCorrelationID struct {
	// Headers will be used to fill the message headers
	Headers struct {
		// Description: Correlation ID set by user
		CorrelationId *string `json:"correlation_id"`
	}

	// Payload will be inserted in the message payload
	Payload struct {
		Event *string `json:"event"`
	}
}

func NewPingWithCorrelationID() PingWithCorrelationID {
	var msg PingWithCorrelationID

	// Set correlation ID
	u := uuid.New().String()
	msg.Headers.CorrelationId = &u

	return msg
}

// newPingWithCorrelationIDFromBrokerMessage will fill a new PingWithCorrelationID with data from generic broker message
func newPingWithCorrelationIDFromBrokerMessage(bMsg extensions.BrokerMessage) (PingWithCorrelationID, error) {
	var msg PingWithCorrelationID

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(bMsg.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// Get each headers from broker message
	for k, v := range bMsg.Headers {
		switch {
		case k == "correlationId": // Retrieving CorrelationId header
			h := string(v)
			msg.Headers.CorrelationId = &h
		default:
			// TODO: log unknown error
		}
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toBrokerMessage will generate a generic broker message from PingWithCorrelationID data
func (msg PingWithCorrelationID) toBrokerMessage() (extensions.BrokerMessage, error) {
	// TODO: implement checks on message

	// Marshal payload to JSON
	payload, err := json.Marshal(msg.Payload)
	if err != nil {
		return extensions.BrokerMessage{}, err
	}

	// Add each headers to broker message
	headers := make(map[string][]byte, 1)

	// Adding CorrelationId header
	if msg.Headers.CorrelationId != nil {
		headers["correlationId"] = []byte(*msg.Headers.CorrelationId)
	}

	return extensions.BrokerMessage{
		Headers: headers,
		Payload: payload,
	}, nil
}

// CorrelationID will give the correlation ID of the message, based on AsyncAPI spec
func (msg PingWithCorrelationID) CorrelationID() string {
	if msg.Headers.CorrelationId != nil {
		return *msg.Headers.CorrelationId
	}

	return ""
}

// SetCorrelationID will set the correlation ID of the message, based on AsyncAPI spec
func (msg *PingWithCorrelationID) SetCorrelationID(id string) {
	msg.Headers.CorrelationId = &id
}

// SetAsResponseFrom will correlate the message with the one passed in parameter.
// It will assign the 'req' message correlation ID to the message correlation ID,
// both specified in AsyncAPI spec.
func (msg *PingWithCorrelationID) SetAsResponseFrom(req MessageWithCorrelationID) {
	id := req.CorrelationID()
	msg.Headers.CorrelationId = &id
}

// Pong is the message expected for 'Pong' channel
type Pong struct {
	// Payload will be inserted in the message payload
	Payload struct {
		Event *string `json:"event"`
	}
}

func NewPong() Pong {
	var msg Pong

	return msg
}

// newPongFromBrokerMessage will fill a new Pong with data from generic broker message
func newPongFromBrokerMessage(bMsg extensions.BrokerMessage) (Pong, error) {
	var msg Pong

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(bMsg.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toBrokerMessage will generate a generic broker message from Pong data
func (msg Pong) toBrokerMessage() (extensions.BrokerMessage, error) {
	// TODO: implement checks on message

	// Marshal payload to JSON
	payload, err := json.Marshal(msg.Payload)
	if err != nil {
		return extensions.BrokerMessage{}, err
	}

	// There is no headers here
	headers := make(map[string][]byte, 0)

	return extensions.BrokerMessage{
		Headers: headers,
		Payload: payload,
	}, nil
}

// PongWithCorrelationID is the message expected for 'PongWithCorrelationID' channel
type PongWithCorrelationID struct {
	// Headers will be used to fill the message headers
	Headers struct {
		// Description: Correlation ID set by user
		CorrelationId *string `json:"correlation_id"`
	}

	// Payload will be inserted in the message payload
	Payload struct {
		Event *string `json:"event"`
	}
}

func NewPongWithCorrelationID() PongWithCorrelationID {
	var msg PongWithCorrelationID

	// Set correlation ID
	u := uuid.New().String()
	msg.Headers.CorrelationId = &u

	return msg
}

// newPongWithCorrelationIDFromBrokerMessage will fill a new PongWithCorrelationID with data from generic broker message
func newPongWithCorrelationIDFromBrokerMessage(bMsg extensions.BrokerMessage) (PongWithCorrelationID, error) {
	var msg PongWithCorrelationID

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(bMsg.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// Get each headers from broker message
	for k, v := range bMsg.Headers {
		switch {
		case k == "correlationId": // Retrieving CorrelationId header
			h := string(v)
			msg.Headers.CorrelationId = &h
		default:
			// TODO: log unknown error
		}
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toBrokerMessage will generate a generic broker message from PongWithCorrelationID data
func (msg PongWithCorrelationID) toBrokerMessage() (extensions.BrokerMessage, error) {
	// TODO: implement checks on message

	// Marshal payload to JSON
	payload, err := json.Marshal(msg.Payload)
	if err != nil {
		return extensions.BrokerMessage{}, err
	}

	// Add each headers to broker message
	headers := make(map[string][]byte, 1)

	// Adding CorrelationId header
	if msg.Headers.CorrelationId != nil {
		headers["correlationId"] = []byte(*msg.Headers.CorrelationId)
	}

	return extensions.BrokerMessage{
		Headers: headers,
		Payload: payload,
	}, nil
}

// CorrelationID will give the correlation ID of the message, based on AsyncAPI spec
func (msg PongWithCorrelationID) CorrelationID() string {
	if msg.Headers.CorrelationId != nil {
		return *msg.Headers.CorrelationId
	}

	return ""
}

// SetCorrelationID will set the correlation ID of the message, based on AsyncAPI spec
func (msg *PongWithCorrelationID) SetCorrelationID(id string) {
	msg.Headers.CorrelationId = &id
}

// SetAsResponseFrom will correlate the message with the one passed in parameter.
// It will assign the 'req' message correlation ID to the message correlation ID,
// both specified in AsyncAPI spec.
func (msg *PongWithCorrelationID) SetAsResponseFrom(req MessageWithCorrelationID) {
	id := req.CorrelationID()
	msg.Headers.CorrelationId = &id
}
