// Package "issue135" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package issue135

import (
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

// AsyncAPIVersion is the version of the used AsyncAPI document
const AsyncAPIVersion = "1.2.3"

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

const (
	// GroupPath is the constant representing the 'Group' channel path.
	GroupPath = "group"
	// InfoPath is the constant representing the 'Info' channel path.
	InfoPath = "info"
	// ProjectPath is the constant representing the 'Project' channel path.
	ProjectPath = "project"
	// ResourcePath is the constant representing the 'Resource' channel path.
	ResourcePath = "resource"
	// StatusPath is the constant representing the 'Status' channel path.
	StatusPath = "status"
)

// ChannelsPaths is an array of all channels paths
var ChannelsPaths = []string{
	GroupPath,
	InfoPath,
	ProjectPath,
	ResourcePath,
	StatusPath,
}
