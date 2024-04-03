// Package "issue175" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package issue175

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

// Type1MessagePayload is a schema from the AsyncAPI specification required in messages
type Type1MessagePayload []Type1MessagePayloadItem

// Type1MessagePayloadItem is a schema from the AsyncAPI specification required in messages
type Type1MessagePayloadItem struct {
	Age   *int64  `json:"age"`
	Email *string `json:"email"`
	Name  *string `json:"name"`
}

// Type1Message is the message expected for 'Type1Message' channel.
type Type1Message struct {
	// Payload will be inserted in the message payload
	Payload []Type1MessagePayloadItem
}

func NewType1Message() Type1Message {
	var msg Type1Message

	return msg
}

// newType1MessageFromBrokerMessage will fill a new Type1Message with data from generic broker message
func newType1MessageFromBrokerMessage(bMsg extensions.BrokerMessage) (Type1Message, error) {
	var msg Type1Message

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(bMsg.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toBrokerMessage will generate a generic broker message from Type1Message data
func (msg Type1Message) toBrokerMessage() (extensions.BrokerMessage, error) {
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

// Type2Message is the message expected for 'Type2Message' channel.
type Type2Message struct {
	// Payload will be inserted in the message payload
	Payload ArrayPayloadSchema
}

func NewType2Message() Type2Message {
	var msg Type2Message

	return msg
}

// newType2MessageFromBrokerMessage will fill a new Type2Message with data from generic broker message
func newType2MessageFromBrokerMessage(bMsg extensions.BrokerMessage) (Type2Message, error) {
	var msg Type2Message

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(bMsg.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toBrokerMessage will generate a generic broker message from Type2Message data
func (msg Type2Message) toBrokerMessage() (extensions.BrokerMessage, error) {
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

// Type3MessagePayload is a schema from the AsyncAPI specification required in messages
type Type3MessagePayload []string

// Type3Message is the message expected for 'Type3Message' channel.
type Type3Message struct {
	// Payload will be inserted in the message payload
	Payload []string
}

func NewType3Message() Type3Message {
	var msg Type3Message

	return msg
}

// newType3MessageFromBrokerMessage will fill a new Type3Message with data from generic broker message
func newType3MessageFromBrokerMessage(bMsg extensions.BrokerMessage) (Type3Message, error) {
	var msg Type3Message

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(bMsg.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toBrokerMessage will generate a generic broker message from Type3Message data
func (msg Type3Message) toBrokerMessage() (extensions.BrokerMessage, error) {
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

// ArrayPayloadSchema is a schema from the AsyncAPI specification required in messages
type ArrayPayloadSchema []ArrayPayloadSchemaItem

// ArrayPayloadSchemaItem is a schema from the AsyncAPI specification required in messages
type ArrayPayloadSchemaItem struct {
	Age   *int64  `json:"age"`
	Email *string `json:"email"`
	Name  *string `json:"name"`
}