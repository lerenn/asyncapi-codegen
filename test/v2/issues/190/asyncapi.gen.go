// Package "issue190" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package issue190

import (
	"encoding/json"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

// AsyncAPIVersion is the version of the used AsyncAPI document
const AsyncAPIVersion = "1.0.0"

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

// Issue169Msg1MessagePayload is a schema from the AsyncAPI specification required in messages
type Issue169Msg1MessagePayload struct {
	Data *Issue169Msg1MessagePayloadData `json:"data"`
}

// Issue169Msg1MessagePayloadData is a schema from the AsyncAPI specification required in messages
type Issue169Msg1MessagePayloadData struct {
	Hello *string `json:"hello"`
	Id    *string `json:"id"`
}

// Issue169Msg1Message is the message expected for 'Issue169Msg1Message' channel.
type Issue169Msg1Message struct {
	// Payload will be inserted in the message payload
	Payload Issue169Msg1MessagePayload
}

func NewIssue169Msg1Message() Issue169Msg1Message {
	var msg Issue169Msg1Message

	return msg
}

// brokerMessageToIssue169Msg1Message will fill a new Issue169Msg1Message with data from generic broker message
func brokerMessageToIssue169Msg1Message(bMsg extensions.BrokerMessage) (Issue169Msg1Message, error) {
	var msg Issue169Msg1Message

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(bMsg.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toBrokerMessage will generate a generic broker message from Issue169Msg1Message data
func (msg Issue169Msg1Message) toBrokerMessage() (extensions.BrokerMessage, error) {
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

// Issue169Msg2MessagePayload is a schema from the AsyncAPI specification required in messages
type Issue169Msg2MessagePayload struct {
	Data *Issue169Msg2MessagePayloadData `json:"data"`
}

// Issue169Msg2MessagePayloadData is a schema from the AsyncAPI specification required in messages
type Issue169Msg2MessagePayloadData struct {
	Bar *string `json:"bar"`
	Id  *string `json:"id"`
}

// Issue169Msg2Message is the message expected for 'Issue169Msg2Message' channel.
type Issue169Msg2Message struct {
	// Payload will be inserted in the message payload
	Payload Issue169Msg2MessagePayload
}

func NewIssue169Msg2Message() Issue169Msg2Message {
	var msg Issue169Msg2Message

	return msg
}

// brokerMessageToIssue169Msg2Message will fill a new Issue169Msg2Message with data from generic broker message
func brokerMessageToIssue169Msg2Message(bMsg extensions.BrokerMessage) (Issue169Msg2Message, error) {
	var msg Issue169Msg2Message

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(bMsg.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toBrokerMessage will generate a generic broker message from Issue169Msg2Message data
func (msg Issue169Msg2Message) toBrokerMessage() (extensions.BrokerMessage, error) {
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

const (
	// Issue169Msg1Path is the constant representing the 'Issue169Msg1' channel path.
	Issue169Msg1Path = "issue169.msg1"
	// Issue169Msg2Path is the constant representing the 'Issue169Msg2' channel path.
	Issue169Msg2Path = "issue169.msg2"
)

// ChannelsPaths is an array of all channels paths
var ChannelsPaths = []string{
	Issue169Msg1Path,
	Issue169Msg2Path,
}
