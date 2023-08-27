// Package "generated" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package generated

import (
	"context"
	"fmt"

	apiContext "github.com/lerenn/asyncapi-codegen/pkg/context"
	"github.com/lerenn/asyncapi-codegen/pkg/log"
	"github.com/lerenn/asyncapi-codegen/pkg/middleware"
)

// ClientSubscriber represents all handlers that are expecting messages for Client
type ClientSubscriber interface {
	// Pong
	Pong(ctx context.Context, msg PongMessage, done bool)
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

func addClientContextValues(ctx context.Context, path string) context.Context {
	ctx = context.WithValue(ctx, apiContext.KeyIsProvider, "client")
	return context.WithValue(ctx, apiContext.KeyIsChannel, path)
}

// Close will clean up any existing resources on the controller
func (c *ClientController) Close(ctx context.Context) {
	// Unsubscribing remaining channels
	c.UnsubscribeAll(ctx)
	c.logger.Info(ctx, "Closed client controller")
}

// SubscribeAll will subscribe to channels without parameters on which the app is expecting messages.
// For channels with parameters, they should be subscribed independently.
func (c *ClientController) SubscribeAll(ctx context.Context, as ClientSubscriber) error {
	if as == nil {
		return ErrNilClientSubscriber
	}

	if err := c.SubscribePong(ctx, as.Pong); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (c *ClientController) UnsubscribeAll(ctx context.Context) {
	// Unsubscribe channels with no parameters (if any)
	c.UnsubscribePong(ctx)

	// Unsubscribe remaining channels
	for n, stopChan := range c.stopSubscribers {
		stopChan <- true
		delete(c.stopSubscribers, n)
	}
}

// SubscribePong will subscribe to new messages from 'pong' channel.
//
// Callback function 'fn' will be called each time a new message is received.
// The 'done' argument indicates when the subscription is canceled and can be
// used to clean up resources.
func (c *ClientController) SubscribePong(ctx context.Context, fn func(ctx context.Context, msg PongMessage, done bool)) error {
	// Get channel path
	path := "pong"

	// Set context
	ctx = addClientContextValues(ctx, path)

	// Check if there is already a subscription
	_, exists := c.stopSubscribers[path]
	if exists {
		err := fmt.Errorf("%w: %q channel is already subscribed", ErrAlreadySubscribedChannel, path)
		c.logger.Error(ctx, err.Error())
		return err
	}

	// Subscribe to broker channel
	msgs, stop, err := c.brokerController.Subscribe(ctx, path)
	if err != nil {
		c.logger.Error(ctx, err.Error())
		return err
	}
	c.logger.Info(ctx, "Subscribed to channel")

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for {
			// Wait for next message
			um, open := <-msgs

			// Add correlation ID to context if it exists
			if um.CorrelationID != nil {
				ctx = context.WithValue(ctx, apiContext.KeyIsCorrelationID, *um.CorrelationID)
			}

			// Process message
			msg, err := newPongMessageFromUniversalMessage(um)
			if err != nil {
				ctx = context.WithValue(ctx, apiContext.KeyIsMessage, um)
				c.logger.Error(ctx, err.Error())
			}

			// Add context
			msgCtx := context.WithValue(ctx, apiContext.KeyIsMessage, msg)
			msgCtx = context.WithValue(msgCtx, apiContext.KeyIsMessageDirection, "reception")

			// Process message if no error and still open
			if err == nil && open {
				// Execute middlewares with the callback
				c.executeMiddlewares(msgCtx, func(ctx context.Context) {
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

// UnsubscribePong will unsubscribe messages from 'pong' channel
func (c *ClientController) UnsubscribePong(ctx context.Context) {
	// Get channel path
	path := "pong"

	// Set context
	ctx = addClientContextValues(ctx, path)

	// Get stop channel
	stopChan, exists := c.stopSubscribers[path]
	if !exists {
		return
	}

	// Stop the channel and remove the entry
	stopChan <- true
	delete(c.stopSubscribers, path)

	c.logger.Info(ctx, "Unsubscribed from channel")
}

// PublishPing will publish messages to 'ping' channel
func (c *ClientController) PublishPing(ctx context.Context, msg PingMessage) error {
	// Get channel path
	path := "ping"

	// Set context
	ctx = addClientContextValues(ctx, path)
	ctx = context.WithValue(ctx, apiContext.KeyIsMessage, msg)
	ctx = context.WithValue(ctx, apiContext.KeyIsMessageDirection, "publication")

	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Add correlation ID to context if it exists
	if um.CorrelationID != nil {
		ctx = context.WithValue(ctx, apiContext.KeyIsCorrelationID, *um.CorrelationID)
	}

	// Publish the message on event-broker through middlewares
	c.executeMiddlewares(ctx, func(ctx context.Context) {
		err = c.brokerController.Publish(ctx, path, um)
	})

	// Return error from publication on broker
	return err
}

// WaitForPong will wait for a specific message by its correlation ID
//
// The pub function is the publication function that should be used to send the message
// It will be called after subscribing to the channel to avoid race condition, and potentially loose the message
func (cc *ClientController) WaitForPong(ctx context.Context, publishMsg MessageWithCorrelationID, pub func(ctx context.Context) error) (PongMessage, error) {
	// Get channel path
	path := "pong"

	// Set context
	ctx = addClientContextValues(ctx, path)

	// Subscribe to broker channel
	msgs, stop, err := cc.brokerController.Subscribe(ctx, path)
	if err != nil {
		cc.logger.Error(ctx, err.Error())
		return PongMessage{}, err
	}
	cc.logger.Info(ctx, "Subscribed to channel")

	// Close subscriber on leave
	defer func() {
		// Unsubscribe
		stop <- true

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
		case um, open := <-msgs:
			// Get new message
			msg, err := newPongMessageFromUniversalMessage(um)
			if err != nil {
				cc.logger.Error(ctx, err.Error())
			}

			// If valid message with corresponding correlation ID, return message
			if err == nil && publishMsg.CorrelationID() == msg.CorrelationID() {
				// Set context with received values
				msgCtx := context.WithValue(ctx, apiContext.KeyIsMessage, msg)
				msgCtx = context.WithValue(msgCtx, apiContext.KeyIsMessageDirection, "reception")
				msgCtx = context.WithValue(msgCtx, apiContext.KeyIsCorrelationID, publishMsg.CorrelationID())

				// Execute middlewares before returning
				cc.executeMiddlewares(msgCtx, func(_ context.Context) {
					/* Nothing to do more */
				})

				return msg, nil
			} else if !open { // If message is invalid or not corresponding and the subscription is closed, then set corresponding error
				cc.logger.Error(ctx, "Channel closed before getting message")
				return PongMessage{}, ErrSubscriptionCanceled
			}
		case <-ctx.Done(): // Set corrsponding error if context is done
			cc.logger.Error(ctx, "Context done before getting message")
			return PongMessage{}, ErrContextCanceled
		}
	}
}
