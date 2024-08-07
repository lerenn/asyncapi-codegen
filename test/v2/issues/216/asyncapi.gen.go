// Package "issue216" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package issue216

import (
	"encoding/json"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

// AsyncAPIVersion is the version of the used AsyncAPI document
const AsyncAPIVersion = "0.1.0"

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

// EventSuccessMessagePayload is a schema from the AsyncAPI specification required in messages
type EventSuccessMessagePayload struct {
	EventObjects []EventSuccessMessagePayloadEventObjectsItem `json:"event_objects,omitempty"`
}

// EventSuccessMessagePayloadEventObjectsItem is a schema from the AsyncAPI specification required in messages
type EventSuccessMessagePayloadEventObjectsItem struct {
	// Description: The identifier of the event
	EventId *string `json:"event_id,omitempty"`

	// Description: The type of the event
	EventType *string `json:"event_type,omitempty"`
}

// EventSuccessMessage is the message expected for 'EventSuccessMessage' channel.
type EventSuccessMessage struct {
	// Payload will be inserted in the message payload
	Payload EventSuccessMessagePayload
}

func NewEventSuccessMessage() EventSuccessMessage {
	var msg EventSuccessMessage

	return msg
}

// brokerMessageToEventSuccessMessage will fill a new EventSuccessMessage with data from generic broker message
func brokerMessageToEventSuccessMessage(bMsg extensions.BrokerMessage) (EventSuccessMessage, error) {
	var msg EventSuccessMessage

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(bMsg.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toBrokerMessage will generate a generic broker message from EventSuccessMessage data
func (msg EventSuccessMessage) toBrokerMessage() (extensions.BrokerMessage, error) {
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
