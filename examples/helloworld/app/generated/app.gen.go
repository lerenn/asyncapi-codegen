// Package "generated" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package generated

import (
	"log"
)

// AppSubscriber represents all application handlers that are expecting messages from clients
type AppSubscriber interface {
	// Hello
	Hello(msg HelloMessage)
}

// AppController is the structure that provides publishing capabilities to the
// developer and and connect the broker with the app
type AppController struct {
	brokerController BrokerController
	stopSubscribers  map[string]chan interface{}
}

// NewAppController links the application to the broker
func NewAppController(bs BrokerController) *AppController {
	// TODO: Check that brokerController is not nil

	return &AppController{
		brokerController: bs,
		stopSubscribers:  make(map[string]chan interface{}),
	}
}

// Close will clean up any existing resources on the controller
func (ac *AppController) Close() {
	ac.UnsubscribeAll()

}

// SubscribeAll will subscribe to channels on which the app is expecting messages
func (ac *AppController) SubscribeAll(as AppSubscriber) error {
	// TODO: Check that as is not nil

	if err := ac.SubscribeHello(as.Hello); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (ac *AppController) UnsubscribeAll() {
	ac.UnsubscribeHello()
}

// SubscribeHello will subscribe to new messages from 'hello' channel
func (ac *AppController) SubscribeHello(fn func(msg HelloMessage)) error {
	// TODO: check if there is already a subscription

	// Subscribe to broker channel
	msgs, stop, err := ac.brokerController.Subscribe("hello")
	if err != nil {
		return err
	}

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for um, open := <-msgs; open; um, open = <-msgs {
			var msg HelloMessage
			if err := msg.fromUniversalMessage(um); err != nil {
				log.Printf("an error happened when receiving an event: %s (msg: %+v)\n", err, msg) // TODO: add proper error handling
				continue
			}

			fn(msg)
		}
	}()

	// Add the stop channel to the inside map
	ac.stopSubscribers["hello"] = stop

	return nil
}

// UnsubscribeHello will unsubscribe messages from 'hello' channel
func (ac *AppController) UnsubscribeHello() {
	stopChan, exists := ac.stopSubscribers["hello"]
	if !exists {
		return
	}

	stopChan <- true
	delete(ac.stopSubscribers, "hello")
}

// Listen will let the controller handle subscriptions and will be interrupted
// only when an struct is sent on the interrupt channel
func (ac *AppController) Listen(irq chan interface{}) {
	<-irq
}
