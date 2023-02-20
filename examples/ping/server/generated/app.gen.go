// Package "generated" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package generated

import (
	"fmt"
)

// AppSubscriber represents all handlers that are expecting messages for App
type AppSubscriber interface {
	// Ping
	Ping(msg PingMessage)
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

	if err := c.SubscribePing(as.Ping); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (c *AppController) UnsubscribeAll() {
	c.UnsubscribePing()
}

// SubscribePing will subscribe to new messages from 'ping' channel
func (c *AppController) SubscribePing(fn func(msg PingMessage)) error {
	// Check if there is already a subscription
	_, exists := c.stopSubscribers["ping"]
	if exists {
		return fmt.Errorf("%w: ping channel is already subscribed", ErrAlreadySubscribedChannel)
	}

	// Subscribe to broker channel
	msgs, stop, err := c.brokerController.Subscribe("ping")
	if err != nil {
		return err
	}

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for um, open := <-msgs; open; um, open = <-msgs {
			msg, err := newPingMessageFromUniversalMessage(um)
			if err != nil {
				c.handleError("ping", err)
			} else {
				fn(msg)
			}
		}
	}()

	// Add the stop channel to the inside map
	c.stopSubscribers["ping"] = stop

	return nil
}

// UnsubscribePing will unsubscribe messages from 'ping' channel
func (c *AppController) UnsubscribePing() {
	stopChan, exists := c.stopSubscribers["ping"]
	if !exists {
		return
	}

	stopChan <- true
	delete(c.stopSubscribers, "ping")
}

// PublishPong will publish messages to 'pong' channel
func (c *AppController) PublishPong(msg PongMessage) error {
	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Publish on event broker
	return c.brokerController.Publish("pong", um)
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
