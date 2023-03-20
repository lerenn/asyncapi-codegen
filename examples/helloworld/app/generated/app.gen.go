// Package "generated" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package generated

import (
	"fmt"
)

// AppSubscriber represents all handlers that are expecting messages for App
type AppSubscriber interface {
	// Hello
	Hello(msg HelloMessage, done bool)
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
func (c *AppController) SubscribeAll(as AppSubscriber) error {
	if as == nil {
		return ErrNilAppSubscriber
	}

	if err := c.SubscribeHello(as.Hello); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (c *AppController) UnsubscribeAll() {
	// Unsubscribe channels with no parameters (if any)
	c.UnsubscribeHello()

	// Unsubscribe remaining channels
	for n, stopChan := range c.stopSubscribers {
		stopChan <- true
		delete(c.stopSubscribers, n)
	}
}

// SubscribeHello will subscribe to new messages from 'hello' channel.
//
// Callback function 'fn' will be called each time a new message is received.
// The 'done' argument indicates when the subscription is canceled and can be
// used to clean up resources.
func (c *AppController) SubscribeHello(fn func(msg HelloMessage, done bool)) error {
	// Get channel path
	path := "hello"

	// Check if there is already a subscription
	_, exists := c.stopSubscribers[path]
	if exists {
		return fmt.Errorf("%w: %q channel is already subscribed", ErrAlreadySubscribedChannel, path)
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
			msg, err := newHelloMessageFromUniversalMessage(um)
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

// UnsubscribeHello will unsubscribe messages from 'hello' channel
func (c *AppController) UnsubscribeHello() {
	// Get channel path
	path := "hello"

	// Get stop channel
	stopChan, exists := c.stopSubscribers[path]
	if !exists {
		return
	}

	// Stop the channel and remove the entry
	stopChan <- true
	delete(c.stopSubscribers, path)
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
