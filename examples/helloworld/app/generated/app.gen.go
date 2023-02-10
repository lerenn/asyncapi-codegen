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
	Hello(msg HelloMessage)
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

	if err := c.SubscribeHello(as.Hello); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (c *AppController) UnsubscribeAll() {
	c.UnsubscribeHello()
}

// SubscribeHello will subscribe to new messages from 'hello' channel
func (c *AppController) SubscribeHello(fn func(msg HelloMessage)) error {
	// Check if there is already a subscription
	_, exists := c.stopSubscribers["hello"]
	if exists {
		return fmt.Errorf("%w: hello channel is already subscribed", ErrAlreadySubscribedChannel)
	}

	// Subscribe to broker channel
	msgs, stop, err := c.brokerController.Subscribe("hello")
	if err != nil {
		return err
	}

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for um, open := <-msgs; open; um, open = <-msgs {
			msg, err := newHelloMessageFromUniversalMessage(um)
			if err != nil {
				c.errChan <- Error{
					Channel: "hello",
					Err:     err,
				}
			} else {
				fn(msg)
			}
		}
	}()

	// Add the stop channel to the inside map
	c.stopSubscribers["hello"] = stop

	return nil
}

// UnsubscribeHello will unsubscribe messages from 'hello' channel
func (c *AppController) UnsubscribeHello() {
	stopChan, exists := c.stopSubscribers["hello"]
	if !exists {
		return
	}

	stopChan <- true
	delete(c.stopSubscribers, "hello")
}

// Listen will let the controller handle subscriptions and will be interrupted
// only when an struct is sent on the interrupt channel
func (c *AppController) Listen(irq <-chan interface{}) {
	<-irq
}
