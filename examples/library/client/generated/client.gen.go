// Package "generated" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package generated

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

// ClientSubscriber represents all application handlers that are expecting messages from application
type ClientSubscriber interface {
	// BooksListResponse
	BooksListResponse(msg BooksListResponseMessage)
}

// ClientController is the structure that provides publishing capabilities to the
// developer and and connect the broker with the client
type ClientController struct {
	brokerController BrokerController
	stopSubscribers  map[string]chan interface{}
}

// NewClientController links the client to the broker
func NewClientController(bs BrokerController) *ClientController {
	// TODO: Check that brokerController is not nil

	return &ClientController{
		brokerController: bs,
		stopSubscribers:  make(map[string]chan interface{}),
	}
}

// Close will clean up any existing resources on the controller
func (cc *ClientController) Close() {
	cc.UnsubscribeAll()

}

// SubscribeAll will subscribe to channels on which the client is expecting messages
func (cc *ClientController) SubscribeAll(cs ClientSubscriber) error {
	// TODO: Check that cs is not nil

	if err := cc.SubscribeBooksListResponse(cs.BooksListResponse); err != nil {
		return err
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (cc *ClientController) UnsubscribeAll() {
	cc.UnsubscribeBooksListResponse()
}

// SubscribeBooksListResponse will subscribe to new messages from 'books.list.response' channel
func (cc *ClientController) SubscribeBooksListResponse(fn func(msg BooksListResponseMessage)) error {
	// Subscribe to broker channel
	msgs, stop, err := cc.brokerController.Subscribe("books.list.response")
	if err != nil {
		return err
	}

	// Asynchronously listen to new messages and pass them to client subscriber
	go func() {
		var msg BooksListResponseMessage
		var um UniversalMessage

		for open := true; open; {
			um, open = <-msgs

			err := json.Unmarshal(um.Payload, &msg.Payload)
			if err != nil {
				log.Println("an error happened when receiving an event:", err) // TODO: add proper error handling
				continue
			}

			// TODO: run checks on data type

			fn(msg)
		}
	}()

	// Add the stop channel to the inside map
	cc.stopSubscribers["books.list.response"] = stop

	return nil
}

// UnsubscribeBooksListResponse will unsubscribe messages from 'books.list.response' channel
func (cc *ClientController) UnsubscribeBooksListResponse() {
	stopChan, exists := cc.stopSubscribers["books.list.response"]
	if !exists {
		return
	}

	stopChan <- true
	delete(cc.stopSubscribers, "books.list.response")
}

// PublishBooksListRequest will publish messages to 'books.list.request' channel
func (cc *ClientController) PublishBooksListRequest(msg BooksListRequestMessage) error {
	// TODO: check that 'cc' is not nil
	// TODO: implement checks on message

	// Convert to JSON payload
	payload, err := json.Marshal(msg.Payload)
	if err != nil {
		return err
	}

	// Create a new correlationID if none is specified
	correlationID := uuid.New().String()
	// TODO: get if from another place according to spec
	if msg.Headers.CorrelationID != "" {
		correlationID = msg.Headers.CorrelationID
	}

	// Create universal message
	um := UniversalMessage{
		Payload:       payload,
		CorrelationID: correlationID,
	}

	// Publish on event broker
	return cc.brokerController.Publish("books.list.request", um)
}
