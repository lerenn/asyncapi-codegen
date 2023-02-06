// Package "generated" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package generated

import (
	"log"
)

// AppSubscriber represents all application handlers that are expecting messages from clients
type AppSubscriber interface {
	// BooksListRequest
	BooksListRequest(msg BooksListRequestMessage)
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

	if err := ac.SubscribeBooksListRequest(as.BooksListRequest); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (ac *AppController) UnsubscribeAll() {
	ac.UnsubscribeBooksListRequest()
}

// SubscribeBooksListRequest will subscribe to new messages from 'books.list.request' channel
func (ac *AppController) SubscribeBooksListRequest(fn func(msg BooksListRequestMessage)) error {
	// TODO: check if there is already a subscription

	// Subscribe to broker channel
	msgs, stop, err := ac.brokerController.Subscribe("books.list.request")
	if err != nil {
		return err
	}

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		var um UniversalMessage

		for open := true; open; um, open = <-msgs {
			var msg BooksListRequestMessage
			if err := msg.fromUniversalMessage(um); err != nil {
				log.Println("an error happened when receiving an event:", err) // TODO: add proper error handling
				continue
			}

			fn(msg)
		}
	}()

	// Add the stop channel to the inside map
	ac.stopSubscribers["books.list.request"] = stop

	return nil
}

// UnsubscribeBooksListRequest will unsubscribe messages from 'books.list.request' channel
func (ac *AppController) UnsubscribeBooksListRequest() {
	stopChan, exists := ac.stopSubscribers["books.list.request"]
	if !exists {
		return
	}

	stopChan <- true
	delete(ac.stopSubscribers, "books.list.request")
}

// PublishBooksListResponse will publish messages to 'books.list.response' channel
func (ac *AppController) PublishBooksListResponse(msg BooksListResponseMessage) error {
	// TODO: check that 'ac' is not nil

	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Publish on event broker
	return ac.brokerController.Publish("books.list.response", um)
}

// Listen will let the controller handle subscriptions and will be interrupted
// only when an struct is sent on the interrupt channel
func (ac *AppController) Listen(irq chan interface{}) {
	<-irq
}
