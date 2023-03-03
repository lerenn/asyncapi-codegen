// Package "generated" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package generated

import (
	"context"
	"fmt"
)

// ClientSubscriber represents all handlers that are expecting messages for Client
type ClientSubscriber interface {
	// Pong
	Pong(msg PongMessage, done bool)
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
	// Unsubscribing remaining channels
	c.UnsubscribeAll()
	// Close the channel and put its reference to nil, if not already closed (= being nil)
	if c.errChan != nil {
		close(c.errChan)
		c.errChan = nil
	}
}

// SubscribeAll will subscribe to channels without parameters on which the app is expecting messages.
// For channels with parameters, they should be subscribed independently.
func (c *ClientController) SubscribeAll(as ClientSubscriber) error {
	if as == nil {
		return ErrNilClientSubscriber
	}

	if err := c.SubscribePong(as.Pong); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (c *ClientController) UnsubscribeAll() {
	// Unsubscribe channels with no parameters (if any)
	c.UnsubscribePong()

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
func (c *ClientController) SubscribePong(fn func(msg PongMessage, done bool)) error {
	// Get channel path
	path := "pong"

	// Check if there is already a subscription
	_, exists := c.stopSubscribers[path]
	if exists {
		return fmt.Errorf("%w: pong channel is already subscribed", ErrAlreadySubscribedChannel)
	}

	// Subscribe to broker channel
	msgs, stop, err := c.brokerController.Subscribe(path)
	if err != nil {
		return err
	}

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for {
			// Wait for next message
			um, open := <-msgs

			// Process message
			msg, err := newPongMessageFromUniversalMessage(um)
			if err != nil {
				c.handleError(path, err)
			}

			// Send info if message is correct or susbcription is closed
			if err == nil || !open {
				fn(msg, !open)
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
func (c *ClientController) UnsubscribePong() {
	// Get channel path
	path := "pong"

	// Get stop channel
	stopChan, exists := c.stopSubscribers[path]
	if !exists {
		return
	}

	// Stop the channel and remove the entry
	stopChan <- true
	delete(c.stopSubscribers, path)
}

// PublishPing will publish messages to 'ping' channel
func (c *ClientController) PublishPing(msg PingMessage) error {
	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Publish on event broker
	path := "ping"
	return c.brokerController.Publish(path, um)
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

// WaitForPong will wait for a specific message by its correlation ID
//
// The pub function is the publication function that should be used to send the message
// It will be called after subscribing to the channel to avoid race condition, and potentially loose the message
func (cc *ClientController) WaitForPong(ctx context.Context, msg MessageWithCorrelationID, pub func() error) (PongMessage, error) {
	// Get channel path
	path := "pong"

	// Subscribe to broker channel
	msgs, stop, err := cc.brokerController.Subscribe(path)
	if err != nil {
		return PongMessage{}, err
	}

	// Close subscriber on leave
	defer func() { stop <- true }()

	// Execute publication
	if err := pub(); err != nil {
		return PongMessage{}, err
	}

	// Wait for corresponding response
	for {
		select {
		case um, open := <-msgs:
			// Get new message
			msg, err := newPongMessageFromUniversalMessage(um)
			if err != nil {
				cc.handleError(path, err)
			}

			// If valid message with corresponding correlation ID, return message
			if err == nil &&
				msg.Headers.CorrelationID != nil && msg.CorrelationID() == *msg.Headers.CorrelationID {
				return msg, nil
			} else if !open { // If message is invalid or not corresponding and the subscription is closed, then return error
				return PongMessage{}, ErrSubscriptionCanceled
			}
		case <-ctx.Done(): // Return error if context is done
			return PongMessage{}, ErrContextCanceled
		}
	}
}
