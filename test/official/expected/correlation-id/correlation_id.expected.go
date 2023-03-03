// Package "correlationID" provides primitives to interact with the AsyncAPI specification.
//
// Code generated by github.com/lerenn/asyncapi-codegen version (devel) DO NOT EDIT.
package correlationID

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AppSubscriber represents all handlers that are expecting messages for App
type AppSubscriber interface {
	// SmartylightingStreetlights10EventStreetlightIDLightingMeasured
	SmartylightingStreetlights10EventStreetlightIDLightingMeasured(msg LightMeasuredMessage, done bool)
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
	// Unsubscribing remaining channels
	c.UnsubscribeAll()
	// Close the channel and put its reference to nil, if not already closed (= being nil)
	if c.errChan != nil {
		close(c.errChan)
		c.errChan = nil
	}
}

// SubscribeAll will subscribe to channels without parameters on which the app is expecting messages.
// For channels with parameters, they should be subscribed independently.
func (c *AppController) SubscribeAll(as AppSubscriber) error {
	if as == nil {
		return ErrNilAppSubscriber
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (c *AppController) UnsubscribeAll() {
	// Unsubscribe channels with no parameters (if any)

	// Unsubscribe remaining channels
	for n, stopChan := range c.stopSubscribers {
		stopChan <- true
		delete(c.stopSubscribers, n)
	}
}

// SubscribeSmartylightingStreetlights10EventStreetlightIDLightingMeasured will subscribe to new messages from 'smartylighting/streetlights/1/0/event/{streetlightId}/lighting/measured' channel.
//
// Callback function 'fn' will be called each time a new message is received.
// The 'done' argument indicates when the subscription is canceled and can be
// used to clean up resources.
func (c *AppController) SubscribeSmartylightingStreetlights10EventStreetlightIDLightingMeasured(params SmartylightingStreetlights10EventStreetlightIDLightingMeasuredParameters, fn func(msg LightMeasuredMessage, done bool)) error {
	// Get channel path
	path := fmt.Sprintf("smartylighting/streetlights/1/0/event/%s/lighting/measured", params.StreetlightID)

	// Check if there is already a subscription
	_, exists := c.stopSubscribers[path]
	if exists {
		return fmt.Errorf("%w: smartylighting/streetlights/1/0/event/{streetlightId}/lighting/measured channel is already subscribed", ErrAlreadySubscribedChannel)
	}

	// Subscribe to broker channel
	msgs, stop, err := c.brokerController.Subscribe(path)
	if err != nil {
		return err
	}

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for {
			// Wait for next message
			um, open := <-msgs

			// Process message
			msg, err := newLightMeasuredMessageFromUniversalMessage(um)
			if err != nil {
				c.handleError(path, err)
			}

			// Send info if message is correct or susbcription is closed
			if err == nil || !open {
				fn(msg, !open)
			}

			// If subscription is closed, then exit the function
			if !open {
				return
			}
		}
	}()

	// Add the stop channel to the inside map
	c.stopSubscribers[path] = stop

	return nil
}

// UnsubscribeSmartylightingStreetlights10EventStreetlightIDLightingMeasured will unsubscribe messages from 'smartylighting/streetlights/1/0/event/{streetlightId}/lighting/measured' channel
func (c *AppController) UnsubscribeSmartylightingStreetlights10EventStreetlightIDLightingMeasured(params SmartylightingStreetlights10EventStreetlightIDLightingMeasuredParameters) {
	// Get channel path
	path := fmt.Sprintf("smartylighting/streetlights/1/0/event/%s/lighting/measured", params.StreetlightID)

	// Get stop channel
	stopChan, exists := c.stopSubscribers[path]
	if !exists {
		return
	}

	// Stop the channel and remove the entry
	stopChan <- true
	delete(c.stopSubscribers, path)
}

// PublishSmartylightingStreetlights10ActionStreetlightIDDim will publish messages to 'smartylighting/streetlights/1/0/action/{streetlightId}/dim' channel
func (c *AppController) PublishSmartylightingStreetlights10ActionStreetlightIDDim(params SmartylightingStreetlights10ActionStreetlightIDDimParameters, msg DimLightMessage) error {
	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Publish on event broker
	path := fmt.Sprintf("smartylighting/streetlights/1/0/action/%s/dim", params.StreetlightID)
	return c.brokerController.Publish(path, um)
}

func (c *AppController) handleError(channelName string, err error) {
	// Wrap error with the channel name
	errWrapped := Error{
		Channel: channelName,
		Err:     err,
	}

	// Send it to the error channel
	select {
	case c.errChan <- errWrapped:
	default:
		// Drop error if it's full or closed
	}
}

// ClientSubscriber represents all handlers that are expecting messages for Client
type ClientSubscriber interface {
	// SmartylightingStreetlights10ActionStreetlightIDDim
	SmartylightingStreetlights10ActionStreetlightIDDim(msg DimLightMessage, done bool)
}

// ClientController is the structure that provides publishing capabilities to the
// developer and and connect the broker with the Client
type ClientController struct {
	brokerController BrokerController
	stopSubscribers  map[string]chan interface{}
	errChan          chan Error
}

// NewClientController links the Client to the broker
func NewClientController(bs BrokerController) (*ClientController, error) {
	if bs == nil {
		return nil, ErrNilBrokerController
	}

	return &ClientController{
		brokerController: bs,
		stopSubscribers:  make(map[string]chan interface{}),
		errChan:          make(chan Error, 256),
	}, nil
}

// Errors will give back the channel that contains errors and that you can listen to handle errors
// Please take a look at Error struct form information on error
func (c ClientController) Errors() <-chan Error {
	return c.errChan
}

// Close will clean up any existing resources on the controller
func (c *ClientController) Close() {
	// Unsubscribing remaining channels
	c.UnsubscribeAll()
	// Close the channel and put its reference to nil, if not already closed (= being nil)
	if c.errChan != nil {
		close(c.errChan)
		c.errChan = nil
	}
}

// SubscribeAll will subscribe to channels without parameters on which the app is expecting messages.
// For channels with parameters, they should be subscribed independently.
func (c *ClientController) SubscribeAll(as ClientSubscriber) error {
	if as == nil {
		return ErrNilClientSubscriber
	}

	return nil
}

// UnsubscribeAll will unsubscribe all remaining subscribed channels
func (c *ClientController) UnsubscribeAll() {
	// Unsubscribe channels with no parameters (if any)

	// Unsubscribe remaining channels
	for n, stopChan := range c.stopSubscribers {
		stopChan <- true
		delete(c.stopSubscribers, n)
	}
}

// SubscribeSmartylightingStreetlights10ActionStreetlightIDDim will subscribe to new messages from 'smartylighting/streetlights/1/0/action/{streetlightId}/dim' channel.
//
// Callback function 'fn' will be called each time a new message is received.
// The 'done' argument indicates when the subscription is canceled and can be
// used to clean up resources.
func (c *ClientController) SubscribeSmartylightingStreetlights10ActionStreetlightIDDim(params SmartylightingStreetlights10ActionStreetlightIDDimParameters, fn func(msg DimLightMessage, done bool)) error {
	// Get channel path
	path := fmt.Sprintf("smartylighting/streetlights/1/0/action/%s/dim", params.StreetlightID)

	// Check if there is already a subscription
	_, exists := c.stopSubscribers[path]
	if exists {
		return fmt.Errorf("%w: smartylighting/streetlights/1/0/action/{streetlightId}/dim channel is already subscribed", ErrAlreadySubscribedChannel)
	}

	// Subscribe to broker channel
	msgs, stop, err := c.brokerController.Subscribe(path)
	if err != nil {
		return err
	}

	// Asynchronously listen to new messages and pass them to app subscriber
	go func() {
		for {
			// Wait for next message
			um, open := <-msgs

			// Process message
			msg, err := newDimLightMessageFromUniversalMessage(um)
			if err != nil {
				c.handleError(path, err)
			}

			// Send info if message is correct or susbcription is closed
			if err == nil || !open {
				fn(msg, !open)
			}

			// If subscription is closed, then exit the function
			if !open {
				return
			}
		}
	}()

	// Add the stop channel to the inside map
	c.stopSubscribers[path] = stop

	return nil
}

// UnsubscribeSmartylightingStreetlights10ActionStreetlightIDDim will unsubscribe messages from 'smartylighting/streetlights/1/0/action/{streetlightId}/dim' channel
func (c *ClientController) UnsubscribeSmartylightingStreetlights10ActionStreetlightIDDim(params SmartylightingStreetlights10ActionStreetlightIDDimParameters) {
	// Get channel path
	path := fmt.Sprintf("smartylighting/streetlights/1/0/action/%s/dim", params.StreetlightID)

	// Get stop channel
	stopChan, exists := c.stopSubscribers[path]
	if !exists {
		return
	}

	// Stop the channel and remove the entry
	stopChan <- true
	delete(c.stopSubscribers, path)
}

// PublishSmartylightingStreetlights10EventStreetlightIDLightingMeasured will publish messages to 'smartylighting/streetlights/1/0/event/{streetlightId}/lighting/measured' channel
func (c *ClientController) PublishSmartylightingStreetlights10EventStreetlightIDLightingMeasured(params SmartylightingStreetlights10EventStreetlightIDLightingMeasuredParameters, msg LightMeasuredMessage) error {
	// Convert to UniversalMessage
	um, err := msg.toUniversalMessage()
	if err != nil {
		return err
	}

	// Publish on event broker
	path := fmt.Sprintf("smartylighting/streetlights/1/0/event/%s/lighting/measured", params.StreetlightID)
	return c.brokerController.Publish(path, um)
}

func (c *ClientController) handleError(channelName string, err error) {
	// Wrap error with the channel name
	errWrapped := Error{
		Channel: channelName,
		Err:     err,
	}

	// Send it to the error channel
	select {
	case c.errChan <- errWrapped:
	default:
		// Drop error if it's full or closed
	}
}

const (
	// CorrelationIDField is the name of the field that will contain the correlation ID
	CorrelationIDField = "correlation_id"
)

// UniversalMessage is a wrapper that will contain all information regarding a message
type UniversalMessage struct {
	CorrelationID *string
	Payload       []byte
}

// BrokerController represents the functions that should be implemented to connect
// the broker to the application or the client
type BrokerController interface {
	// Publish a message to the broker
	Publish(channel string, mw UniversalMessage) error

	// Subscribe to messages from the broker
	Subscribe(channel string) (msgs chan UniversalMessage, stop chan interface{}, err error)
}

var (
	// Generic error for AsyncAPI generated code
	ErrAsyncAPI = errors.New("error when using AsyncAPI")

	// ErrContextCanceled is given when a given context is canceled
	ErrContextCanceled = fmt.Errorf("%w: context canceled", ErrAsyncAPI)

	// ErrNilBrokerController is raised when a nil broker controller is user
	ErrNilBrokerController = fmt.Errorf("%w: nil broker controller has been used", ErrAsyncAPI)

	// ErrNilAppSubscriber is raised when a nil app subscriber is user
	ErrNilAppSubscriber = fmt.Errorf("%w: nil app subscriber has been used", ErrAsyncAPI)

	// ErrNilClientSubscriber is raised when a nil client subscriber is user
	ErrNilClientSubscriber = fmt.Errorf("%w: nil client subscriber has been used", ErrAsyncAPI)

	// ErrAlreadySubscribedChannel is raised when a subscription is done twice
	// or more without unsubscribing
	ErrAlreadySubscribedChannel = fmt.Errorf("%w: the channel has already been subscribed", ErrAsyncAPI)

	// ErrSubscriptionCanceled is raised when expecting something and the subscription has been canceled before it happens
	ErrSubscriptionCanceled = fmt.Errorf("%w: the subscription has been canceled", ErrAsyncAPI)
)

type MessageWithCorrelationID interface {
	CorrelationID() string
}

type Error struct {
	Channel string
	Err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("channel %q: err %v", e.Channel, e.Err)
}

// SmartylightingStreetlights10ActionStreetlightIDDimParameters represents SmartylightingStreetlights10ActionStreetlightIDDim channel parameters
type SmartylightingStreetlights10ActionStreetlightIDDimParameters struct {
	StreetlightID string
}

// SmartylightingStreetlights10EventStreetlightIDLightingMeasuredParameters represents SmartylightingStreetlights10EventStreetlightIDLightingMeasured channel parameters
type SmartylightingStreetlights10EventStreetlightIDLightingMeasuredParameters struct {
	StreetlightID string
}

// DimLightMessage is the message expected for 'DimLight' channel
type DimLightMessage struct {
	// Payload will be inserted in the message payload
	Payload DimLightPayloadSchema
}

func NewDimLightMessage() DimLightMessage {
	var msg DimLightMessage

	return msg
}

// newDimLightMessageFromUniversalMessage will fill a new DimLightMessage with data from UniversalMessage
func newDimLightMessageFromUniversalMessage(um UniversalMessage) (DimLightMessage, error) {
	var msg DimLightMessage

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(um.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// TODO: run checks on msg type

	return msg, nil
}

// toUniversalMessage will generate an UniversalMessage from DimLightMessage data
func (msg DimLightMessage) toUniversalMessage() (UniversalMessage, error) {
	// TODO: implement checks on message

	// Marshal payload to JSON
	payload, err := json.Marshal(msg.Payload)
	if err != nil {
		return UniversalMessage{}, err
	}

	return UniversalMessage{
		Payload: payload,
	}, nil
}

// LightMeasuredMessage is the message expected for 'LightMeasured' channel
type LightMeasuredMessage struct {
	// Headers will be used to fill the message headers
	Headers struct {
		Mqmd *struct {
			CorrelID *string `json:"correl_id"`
		} `json:"mqmd"`
	}

	// Payload will be inserted in the message payload
	Payload LightMeasuredPayloadSchema
}

func NewLightMeasuredMessage() LightMeasuredMessage {
	var msg LightMeasuredMessage

	// Set correlation ID
	u := uuid.New().String()
	msg.Headers.Mqmd.CorrelID = &u

	return msg
}

// newLightMeasuredMessageFromUniversalMessage will fill a new LightMeasuredMessage with data from UniversalMessage
func newLightMeasuredMessageFromUniversalMessage(um UniversalMessage) (LightMeasuredMessage, error) {
	var msg LightMeasuredMessage

	// Unmarshal payload to expected message payload format
	err := json.Unmarshal(um.Payload, &msg.Payload)
	if err != nil {
		return msg, err
	}

	// Get correlation ID
	msg.Headers.Mqmd.CorrelID = um.CorrelationID

	// TODO: run checks on msg type

	return msg, nil
}

// toUniversalMessage will generate an UniversalMessage from LightMeasuredMessage data
func (msg LightMeasuredMessage) toUniversalMessage() (UniversalMessage, error) {
	// TODO: implement checks on message

	// Marshal payload to JSON
	payload, err := json.Marshal(msg.Payload)
	if err != nil {
		return UniversalMessage{}, err
	}

	// Set correlation ID if it does not exist
	var correlationID *string
	if msg.Headers.Mqmd.CorrelID != nil {
		correlationID = msg.Headers.Mqmd.CorrelID
	} else {
		u := uuid.New().String()
		correlationID = &u
	}

	return UniversalMessage{
		Payload:       payload,
		CorrelationID: correlationID,
	}, nil
}

// CorrelationID will give the correlation ID of the message, based on AsyncAPI spec
func (msg LightMeasuredMessage) CorrelationID() string {
	if msg.Headers.Mqmd.CorrelID != nil {
		return *msg.Headers.Mqmd.CorrelID
	}

	return ""
}

// SetAsResponseFrom will correlate the message with the one passed in parameter.
// It will assign the 'req' message correlation ID to the message correlation ID,
// both specified in AsyncAPI spec.
func (msg *LightMeasuredMessage) SetAsResponseFrom(req MessageWithCorrelationID) {
	id := req.CorrelationID()
	msg.Headers.Mqmd.CorrelID = &id
}

// DimLightPayloadSchema is a schema from the AsyncAPI specification required in messages
type DimLightPayloadSchema struct {
	// Description: Percentage to which the light should be dimmed to.
	Percentage *int64 `json:"percentage"`

	// Description: Date and time when the message was sent.
	SentAt *SentAtSchema `json:"sent_at"`
}

// LightMeasuredPayloadSchema is a schema from the AsyncAPI specification required in messages
type LightMeasuredPayloadSchema struct {
	// Description: Light intensity measured in lumens.
	Lumens *int64 `json:"lumens"`

	// Description: Date and time when the message was sent.
	SentAt *SentAtSchema `json:"sent_at"`
}

// SentAtSchema is a schema from the AsyncAPI specification required in messages
// Description: Date and time when the message was sent.
type SentAtSchema time.Time

// MarshalJSON will override the marshal as this is not a normal 'time.Time' type
func (t SentAtSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t))
}

// UnmarshalJSON will override the unmarshal as this is not a normal 'time.Time' type
func (t *SentAtSchema) UnmarshalJSON(data []byte) error {
	var timeFormat time.Time
	if err := json.Unmarshal(data, &timeFormat); err != nil {
		return err
	}

	*t = SentAtSchema(timeFormat)
	return nil
}
