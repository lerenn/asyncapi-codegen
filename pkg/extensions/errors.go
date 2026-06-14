package extensions

import (
	"errors"
	"fmt"
)

var (
	// ErrAsyncAPI is the generic error for AsyncAPI generated code.
	ErrAsyncAPI = errors.New("error when using AsyncAPI")

	// ErrContextCanceled is given when a given context is canceled.
	ErrContextCanceled = fmt.Errorf("%w: context canceled", ErrAsyncAPI)

	// ErrNilBrokerController is raised when a nil broker controller is user.
	ErrNilBrokerController = fmt.Errorf("%w: nil broker controller has been used", ErrAsyncAPI)

	// ErrNilAppSubscriber is raised when a nil app subscriber is used (asyncapiv2 only).
	ErrNilAppSubscriber = fmt.Errorf("%w: nil app subscriber has been used", ErrAsyncAPI)

	// ErrNilUserSubscriber is raised when a nil user subscriber is used (asyncapiv2 only).
	ErrNilUserSubscriber = fmt.Errorf("%w: nil user subscriber has been used", ErrAsyncAPI)

	// ErrAlreadySubscribedChannel is raised when a subscription is done twice
	// or more without unsubscribing.
	ErrAlreadySubscribedChannel = fmt.Errorf("%w: the channel has already been subscribed", ErrAsyncAPI)

	// ErrSubscriptionCanceled is raised when expecting something and the subscription has been canceled before it happens.
	ErrSubscriptionCanceled = fmt.Errorf("%w: the subscription has been canceled", ErrAsyncAPI)

	// ErrNoCorrelationIDSet is raised when a correlation ID is expected, but none is detected.
	ErrNoCorrelationIDSet = fmt.Errorf("%w: no correlation ID but one is expected", ErrAsyncAPI)

	// ErrChannelAddressEmpty is raised when a given channel address is empty,
	// when dynamically set from message.
	ErrChannelAddressEmpty = fmt.Errorf("%w: channel address empty", ErrAsyncAPI)

	// ErrNoMatchingMessage is raised when a message received on a channel that
	// carries several message types does not match any of the expected types,
	// i.e. it could not be unmarshalled and validated against any of them.
	ErrNoMatchingMessage = fmt.Errorf("%w: no expected message type matched the received message", ErrAsyncAPI)
)
