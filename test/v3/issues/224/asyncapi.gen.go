// Package "issue224" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package issue224

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

// AppSubscriber contains all handlers that are listening messages for App
type AppSubscriber interface {
	// ReceiveTestOperationReceived receive all TestMessageFromTestChannel messages from Test channel.
	ReceiveTestOperationReceived(ctx context.Context, msg TestMessageFromTestChannel) error
}

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
		errorHandler:  extensions.DefaultErrorHandler(),
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
	c.UnsubscribeFromAllChannels(ctx)

	c.logger.Info(ctx, "Closed app controller")
}

// SubscribeToAllChannels will receive messages from channels where channel has
// no parameter on which the app is expecting messages. For channels with parameters,
// they should be subscribed independently.
func (c *AppController) SubscribeToAllChannels(ctx context.Context, as AppSubscriber) error {
	if as == nil {
		return extensions.ErrNilAppSubscriber
	}

	if err := c.SubscribeToReceiveTestOperation(ctx, as.ReceiveTestOperationReceived); err != nil {
		return err
	}

	return nil
}

// UnsubscribeFromAllChannels will stop the subscription of all remaining subscribed channels
func (c *AppController) UnsubscribeFromAllChannels(ctx context.Context) {
	c.UnsubscribeFromReceiveTestOperation(ctx)
}

// SubscribeToReceiveTestOperation will receive TestMessageFromTestChannel messages from Test channel.
//
// Callback function 'fn' will be called each time a new message is received.
//
// NOTE: for now, this only support the first message from AsyncAPI list.
//
// NOTE: for now, this only support the first message from AsyncAPI list.
// If you need support for other messages, please raise an issue.
func (c *AppController) SubscribeToReceiveTestOperation(
	ctx context.Context,
	fn func(ctx context.Context, msg TestMessageFromTestChannel) error,
) error {
	// Get channel address
	addr := "v3.issue224.test"

	// Set context
	ctx = addAppContextValues(ctx, addr)
	ctx = context.WithValue(ctx, extensions.ContextKeyIsDirection, "reception")

	// Check if the controller is already subscribed
	_, exists := c.subscriptions[addr]
	if exists {
		err := fmt.Errorf("%w: controller is already subscribed on channel %q", extensions.ErrAlreadySubscribedChannel, addr)
		c.logger.Error(ctx, err.Error())
		return err
	}

	// Subscribe to broker channel
	sub, err := c.broker.Subscribe(ctx, addr)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return err
	}
	c.logger.Info(ctx, "Subscribed to channel")

	// Asynchronously listen to new messages and pass them to app receiver
	go func() {
		for {
			// Listen to next message
			stop, err := c.listenToReceiveTestOperationNextMessage(addr, sub, fn)
			if err != nil {
				c.logger.Error(ctx, err.Error())
			}

			// Stop if required
			if stop {
				return
			}
		}
	}()

	// Add the cancel channel to the inside map
	c.subscriptions[addr] = sub

	return nil
}

func (c *AppController) listenToReceiveTestOperationNextMessage(
	addr string,
	sub extensions.BrokerChannelSubscription,
	fn func(ctx context.Context, msg TestMessageFromTestChannel) error,
) (stop bool, err error) {
	// Create a context for the received response
	msgCtx, cancel := context.WithCancel(context.Background())
	msgCtx = addAppContextValues(msgCtx, addr)
	msgCtx = context.WithValue(msgCtx, extensions.ContextKeyIsDirection, "reception")
	defer cancel()

	// Wait for next message
	acknowledgeableBrokerMessage, open := <-sub.MessagesChannel()

	// If subscription is closed and there is no more message
	// (i.e. uninitialized message), then exit the function
	if !open && acknowledgeableBrokerMessage.IsUninitialized() {
		return true, nil
	}

	// Set broker message to context
	msgCtx = context.WithValue(msgCtx, extensions.ContextKeyIsBrokerMessage, acknowledgeableBrokerMessage.String())

	// Execute middlewares before handling the message
	if err := c.executeMiddlewares(msgCtx, &acknowledgeableBrokerMessage.BrokerMessage, func(middlewareCtx context.Context) error {
		// Process message
		msg, err := brokerMessageToTestMessageFromTestChannel(acknowledgeableBrokerMessage.BrokerMessage)
		if err != nil {
			return err
		}

		// Execute the subscription function
		if err := fn(middlewareCtx, msg); err != nil {
			return err
		}

		acknowledgeableBrokerMessage.Ack()

		return nil
	}); err != nil {
		c.errorHandler(msgCtx, addr, &acknowledgeableBrokerMessage, err)
		// On error execute the acknowledgeableBrokerMessage nack() function and
		// let the BrokerAcknowledgment decide what is the right nack behavior for the broker
		acknowledgeableBrokerMessage.Nak()
	}

	return false, nil
}

// UnsubscribeFromReceiveTestOperation will stop the reception of TestMessageFromTestChannel messages from Test channel.
// A timeout can be set in context to avoid blocking operation, if needed.
func (c *AppController) UnsubscribeFromReceiveTestOperation(
	ctx context.Context,
) {
	// Get channel address
	addr := "v3.issue224.test"

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

	c.logger.Info(ctx, "Unsubscribed from channel")
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
		errorHandler:  extensions.DefaultErrorHandler(),
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
}

// SendToReceiveTestOperation will send a TestMessageFromTestChannel message on Test channel.
//
// NOTE: for now, this only support the first message from AsyncAPI list.
// If you need support for other messages, please raise an issue.
func (c *UserController) SendToReceiveTestOperation(
	ctx context.Context,
	msg TestMessageFromTestChannel,
) error {
	// Set channel address
	addr := "v3.issue224.test"

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
	// handler to handle errors from consumers and middlewares
	errorHandler extensions.ErrorHandler
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

// WithErrorHandler attaches a errorhandler to handle errors from subscriber functions
func WithErrorHandler(handler extensions.ErrorHandler) ControllerOption {
	return func(controller *controller) {
		controller.errorHandler = handler
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

// TestMessageFromTestChannel is the message expected for 'TestMessageFromTestChannel' channel.
type TestMessageFromTestChannel struct {
	// Payload will be inserted in the message payload
	Payload ColliderDictionarySchema
}

func NewTestMessageFromTestChannel() TestMessageFromTestChannel {
	var msg TestMessageFromTestChannel

	return msg
}

// brokerMessageToTestMessageFromTestChannel will fill a new TestMessageFromTestChannel with data from generic broker message
func brokerMessageToTestMessageFromTestChannel(bMsg extensions.BrokerMessage) (TestMessageFromTestChannel, error) {
	var msg TestMessageFromTestChannel

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(bMsg.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toBrokerMessage will generate a generic broker message from TestMessageFromTestChannel data
func (msg TestMessageFromTestChannel) toBrokerMessage() (extensions.BrokerMessage, error) {
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

// ColliderSchema is a schema from the AsyncAPI specification required in messages
type ColliderSchema struct {
	Margin *float32                        `json:"margin,omitempty"`
	Pose   *PoseSchema                     `json:"pose,omitempty"`
	Shape  ShapePropertyFromColliderSchema `json:"shape"`
}

// ShapePropertyFromColliderSchema is a schema from the AsyncAPI specification required in messages
type ShapePropertyFromColliderSchema struct {
	Radius    float64 `json:"radius"`
	ShapeType string  `json:"shape_type" validate:"oneof='sphere'"`
}

// ColliderDictionarySchema is a schema from the AsyncAPI specification required in messages
type ColliderDictionarySchema struct {
	// AdditionalProperties represents the object additional properties.
	AdditionalProperties map[string]ColliderSchema `json:"-"`
}

// MarshalJSON marshals the schema into JSON with support for additional properties.
func (t ColliderDictionarySchema) MarshalJSON() ([]byte, error) {
	type alias ColliderDictionarySchema

	// Copy original into alias and marshal the alias to avoid JSON marshal recursion
	b, err := json.Marshal(alias(t))
	if err != nil {
		return nil, err
	}

	// Remove the end of the json (i.e. '}')
	b = b[:len(b)-1]

	// When there are no properties, we cant start with a separator
	needSeparator := len(b) > 1

	// Add additional properties
	for k, v := range t.AdditionalProperties {
		if needSeparator {
			b = append(b, ',')
		}
		needSeparator = true

		vBytes, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		b = append(b, fmt.Sprintf("%q:%s", k, vBytes)...)
	}

	// Close JSON and return
	return append(b, []byte("}")...), nil
}

// UnmarshalJSON unmarshals schema from JSON with support for additional properties.
func (t *ColliderDictionarySchema) UnmarshalJSON(data []byte) error {
	type alias ColliderDictionarySchema

	// Unmarshal to map to get all fields
	var m map[string]ColliderSchema
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	// Unmarshal into the alias then copy the alias content into the original
	// object. This is done to avoid JSON unmarshal recursion.
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	*t = ColliderDictionarySchema(a)

	// Get all fields that are additional and add them to the AdditionalProperties field.
	t.AdditionalProperties = make(map[string]ColliderSchema, len(m))
	for k, v := range m {
		switch k {
		default:
			t.AdditionalProperties[k] = v
		}
	}

	return nil
}

// PoseSchema is a schema from the AsyncAPI specification required in messages
type PoseSchema struct {
	Orientation *Vector3dSchema `json:"orientation,omitempty"`
	Position    *Vector3dSchema `json:"position,omitempty"`
}

// SphereSchema is a schema from the AsyncAPI specification required in messages
type SphereSchema struct {
	Radius    float64 `json:"radius"`
	ShapeType string  `json:"shape_type" validate:"oneof='sphere'"`
}

// Vector3dSchema is a schema from the AsyncAPI specification required in messages
type Vector3dSchema []float64

const (
	// TestChannelPath is the constant representing the 'TestChannel' channel path.
	TestChannelPath = "v3.issue224.test"
)

// ChannelsPaths is an array of all channels paths
var ChannelsPaths = []string{
	TestChannelPath,
}
