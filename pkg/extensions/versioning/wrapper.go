package versioning

import (
	"context"
	"fmt"
	"sync"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

// Check that it still fills the interface.
var _ extensions.BrokerController = (*Wrapper)(nil)

// DefaultVersionHeaderKey is the field that will be added to a message to get the version.
const DefaultVersionHeaderKey = "application-version"

var (
	// ErrNoVersion happens when there is no version in the context or the message.
	ErrNoVersion = fmt.Errorf("%w: no version present", extensions.ErrAsyncAPI)
)

// Wrapper allows to use multiple version of the same App/User Controllers
// on one Broker Controller in order to handle migrations.
type Wrapper struct {
	broker           extensions.BrokerController
	logger           extensions.Logger
	defaultVersion   *string
	versionHeaderKey string

	channels      map[string]*brokerSubscription
	channelsMutex sync.Mutex
}

// WrapperOption adds an option to Version Wrapper.
type WrapperOption func(versionWrapper *Wrapper)

// NewWrapper creates a Version Wrapper around a Broker Controller.
func NewWrapper(broker extensions.BrokerController, options ...WrapperOption) *Wrapper {
	// Create version Wrapper
	w := Wrapper{
		broker:           broker,
		channels:         make(map[string]*brokerSubscription),
		logger:           extensions.DummyLogger{},
		versionHeaderKey: DefaultVersionHeaderKey,
	}

	// Execute options
	for _, option := range options {
		option(&w)
	}

	return &w
}

// WithLogger lets add a logger to the Wrapper struct.
func WithLogger(logger extensions.Logger) WrapperOption {
	return func(wrapper *Wrapper) {
		wrapper.logger = logger
	}
}

// WithDefaultVersion lets add a default version to tag messages that don't have
// versions tagged.
func WithDefaultVersion(version string) WrapperOption {
	return func(wrapper *Wrapper) {
		wrapper.defaultVersion = &version
	}
}

// WithVersionHeaderKey lets use a different version header key to add application
// version to published messages and retrieve it from received messages.
func WithVersionHeaderKey(headerKey string) WrapperOption {
	return func(wrapper *Wrapper) {
		wrapper.versionHeaderKey = headerKey
	}
}

// Publish a message to the broker.
func (w *Wrapper) Publish(ctx context.Context, channel string, mw extensions.BrokerMessage) error {
	// Add version to message
	extensions.IfContextSetWith(ctx, extensions.ContextKeyIsVersion, func(version string) {
		mw.Headers[w.versionHeaderKey] = []byte(version)
	})

	// Send message
	return w.broker.Publish(ctx, channel, mw)
}

// Subscribe to messages from the broker.
func (w *Wrapper) Subscribe(ctx context.Context, channel string) (extensions.BrokerChannelSubscription, error) {
	// Set context
	ctx = context.WithValue(ctx, extensions.ContextKeyIsMessageDirection, "reception")
	ctx = context.WithValue(ctx, extensions.ContextKeyIsChannel, channel)

	// Get version
	var version string
	extensions.IfContextSetWith(ctx, extensions.ContextKeyIsVersion, func(v string) { version = v })
	if version == "" {
		return extensions.BrokerChannelSubscription{}, ErrNoVersion
	}

	// Lock the channels to avoid conflict
	w.channelsMutex.Lock()
	defer w.channelsMutex.Unlock()

	// Check if the broker channel already exists
	brokerChannel, exists := w.channels[channel]
	if !exists {
		cbb, err := w.createBrokerChannels(ctx, channel)
		if err != nil {
			return extensions.BrokerChannelSubscription{}, err
		}
		defer cbb.launchListener(ctx)
		brokerChannel = cbb
	}

	// Check if the version already exists
	cbv, err := brokerChannel.createVersionListener(version)

	return cbv.subscription, err
}

func (w *Wrapper) createBrokerChannels(ctx context.Context, channel string) (*brokerSubscription, error) {
	// Subscribe to broker
	subscription, err := w.broker.Subscribe(ctx, channel)
	if err != nil {
		return nil, err
	}

	// Add channels from broker to brokerChannels
	cbb := newBrokerSubscription(channel, subscription, w)
	w.channels[channel] = &cbb // Already locked in parent function

	return &cbb, nil
}
