// Package "generated" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package generated

import (
	"fmt"
	"log"
)

// AppSubscriber represents all application handlers that are expecting messages from clients
type AppSubscriber interface {
	// Ping
	Ping(msg PingMessage)
}

// AppController is the structure that provides publishing capabilities to the
// developer and and connect the broker with the app
type AppController struct {
	brokerController BrokerController
	stopSubscribers  map[string]chan interface{}
}

// NewAppController links the application to the broker
func NewAppController(bs BrokerController) (*AppController, error) {
	if bs == nil {
		return nil, ErrNilBrokerController
	}

	return &AppController{
		brokerController: bs,
		stopSubscribers:  make(map[string]chan interface{}),
	}, nil
}

// Close will clean up any existing resources on the controller
func (ac *AppController) Close() {
	ac.UnsubscribeAll()

}

// SubscribeAll will subscribe to channels on which the app is expecting messages
func (ac *AppController) SubscribeAll(as AppSubscriber) error {
	if as == nil {
		return ErrNilAppController
	}

	if err := ac.SubscribePing(as.Ping); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (ac *AppController) UnsubscribeAll() {
	ac.UnsubscribePing()
}

// SubscribePing will subscribe to new messages from 'ping' channel
func (ac *AppController) SubscribePing(fn func(msg PingMessage)) error {
	// Check if there is already a subscription
	_, exists := ac.stopSubscribers["ping"]
	if exists {
		return fmt.Errorf("%w: ping channel is already subscribed", ErrAlreadySubscribedChannel)
	}

	// Subscribe to broker channel
	msgs, stop, err := ac.brokerController.Subscribe("ping")
	if err != nil {
		return err
	}

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for um, open := <-msgs; open; um, open = <-msgs {
			msg, err := newPingMessageFromUniversalMessage(um)
			if err != nil {
				log.Printf("an error happened when receiving an event: %s (msg: %+v)\n", err, msg) // TODO: add proper error handling
				continue
			}

			fn(msg)
		}
	}()

	// Add the stop channel to the inside map
	ac.stopSubscribers["ping"] = stop

	return nil
}

// UnsubscribePing will unsubscribe messages from 'ping' channel
func (ac *AppController) UnsubscribePing() {
	stopChan, exists := ac.stopSubscribers["ping"]
	if !exists {
		return
	}

	stopChan <- true
	delete(ac.stopSubscribers, "ping")
}

// PublishPong will publish messages to 'pong' channel
func (ac *AppController) PublishPong(msg PongMessage) error {
	// TODO: check that 'ac' is not nil

	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Publish on event broker
	return ac.brokerController.Publish("pong", um)
}

// Listen will let the controller handle subscriptions and will be interrupted
// only when an struct is sent on the interrupt channel
func (ac *AppController) Listen(irq chan interface{}) {
	<-irq
}
