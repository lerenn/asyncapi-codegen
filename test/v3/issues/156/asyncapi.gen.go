// Package "issue156" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package issue156

import (
	"encoding/json"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

// AsyncAPIVersion is the version of the used AsyncAPI document
const AsyncAPIVersion = ""

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

// TestingMessage is the message expected for 'TestingMessage' channel.
type TestingMessage struct {
	// Payload will be inserted in the message payload
	Payload TestSchema
}

func NewTestingMessage() TestingMessage {
	var msg TestingMessage

	return msg
}

// brokerMessageToTestingMessage will fill a new TestingMessage with data from generic broker message
func brokerMessageToTestingMessage(bMsg extensions.BrokerMessage) (TestingMessage, error) {
	var msg TestingMessage

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(bMsg.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toBrokerMessage will generate a generic broker message from TestingMessage data
func (msg TestingMessage) toBrokerMessage() (extensions.BrokerMessage, error) {
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

// SubTestSchema is a schema from the AsyncAPI specification required in messages
type SubTestSchema string

// TestSchema is a schema from the AsyncAPI specification required in messages
type TestSchema struct {
	Subtest *SubTestSchema `json:"subtest,omitempty"`
}
