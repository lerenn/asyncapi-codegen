package extensions

import (
	"context"
	"fmt"
	"sync"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers"
)

var _ BrokerController = (*VersionWrapper)(nil)

// VersionField is the field that will be added to a message to get the version.
const VersionField = "application-version"

var (
	// ErrNoVersion happens when there is no version in the context or the message.
	ErrNoVersion = fmt.Errorf("%w: no version present", ErrAsyncAPI)
)

type channelsByVersion struct {
	version  string
	messages chan BrokerMessage
	cancel   chan interface{}
	parent   *channelsByBroker
}

func newChannelsByVersion(version string, parent *channelsByBroker) channelsByVersion {
	return channelsByVersion{
		version:  version,
		messages: make(chan BrokerMessage, brokers.BrokerMessagesQueueSize),
		cancel:   make(chan interface{}, 1),
		parent:   parent,
	}
}

func (cbv *channelsByVersion) launchListener() {
	go func() {
		// Wait to receive cancel
		<-cbv.cancel

		// When cancel is received, then remove version listener
		cbv.parent.removeVersionListener(cbv)
	}()
}

func (cbv *channelsByVersion) closeChannels() {
	// Receiving no more messages
	close(cbv.messages)

	// Closing cancel channel to let caller knows that everything is cleaned up
	close(cbv.cancel)
}

type channelsByBroker struct {
	channelName string
	messages    chan BrokerMessage
	cancel      chan interface{}
	parent      *VersionWrapper

	versionsChannels map[string]channelsByVersion
	versionsMutex    sync.Mutex
}

func newChannelsFromBroker(
	channel string,
	messages chan BrokerMessage,
	cancel chan interface{},
	parent *VersionWrapper,
) channelsByBroker {
	return channelsByBroker{
		channelName:      channel,
		messages:         messages,
		cancel:           cancel,
		versionsChannels: make(map[string]channelsByVersion),
		parent:           parent,
	}
}

func (cbb *channelsByBroker) createVersionListener(version string) (channelsByVersion, error) {
	// Lock the versions to avoid conflict
	cbb.versionsMutex.Lock()
	defer cbb.versionsMutex.Unlock()

	// Check if the version doesn't exist already
	_, exists := cbb.versionsChannels[version]
	if exists {
		return channelsByVersion{}, ErrAlreadySubscribedChannel
	}

	// Create the channels necessary
	cbv := newChannelsByVersion(version, cbb)
	cbb.versionsChannels[version] = cbv
	defer cbv.launchListener()

	return cbv, nil
}

func (cbb *channelsByBroker) removeVersionListener(cbv *channelsByVersion) {
	// Lock the versions to avoid conflict
	cbb.versionsMutex.Lock()
	defer cbb.versionsMutex.Unlock()

	// Cleanup the channelsByVersion when leaving
	//
	// NOTE: this is important to make it cleanup at the end of this function as
	// it should be cleanup AFTER the broker have been stopped (in case it was
	// the last version listener), in order to let the caller knows that everything
	// was cleaned up properly.
	defer cbv.closeChannels()

	// Remove the version from the channelsByBroker
	delete(cbb.versionsChannels, cbv.version)

	// Lock the channels to avoid conflict
	cbb.parent.channelsMutex.Lock()
	defer cbb.parent.channelsMutex.Unlock()

	// If there is still version channels, do nothing
	if len(cbb.versionsChannels) > 0 {
		return
	}

	// Otherwise cancel the broker listener and wait for its closure
	cbb.cancel <- true
	<-cbb.cancel

	// Then delete the channelsByBroker from the Version Switch Wrapper
	delete(cbb.parent.channels, cbb.channelName)
}

func (cbb *channelsByBroker) launchListener(ctx context.Context) {
	go func() {
		for {
			// Wait for new messages
			msg := <-cbb.messages

			// Get the version from the message
			version := string(msg.Headers[VersionField])

			// Lock the versions to avoid conflict
			cbb.versionsMutex.Lock()

			// Get the correct channel based on the version
			ch, exists := cbb.versionsChannels[version]
			if !exists {
				// Set context
				ctx = context.WithValue(ctx, ContextKeyIsBrokerMessage, msg)
				ctx = context.WithValue(ctx, ContextKeyIsVersion, version)

				// Log the error
				cbb.parent.logger.Error(ctx, fmt.Sprintf("version %q is not registered", version))
			}

			// Unlock the versions
			cbb.versionsMutex.Unlock()

			// Send the message to the correct channel
			ch.messages <- msg
		}
	}()
}

// VersionWrapper allows to use multiple version of the same App/User Controllers
// on one Broker Controller in order to handle migrations.
type VersionWrapper struct {
	broker BrokerController
	logger Logger

	channels      map[string]*channelsByBroker
	channelsMutex sync.Mutex
}

// VersionWrapperOption adds an option to Version Wrapper.
type VersionWrapperOption func(versionWrapper *VersionWrapper)

// NewVersionWrapper creates a Version Wrapper around a Broker Controller.
func NewVersionWrapper(broker BrokerController, options ...VersionWrapperOption) *VersionWrapper {
	// Create version Wrapper
	vw := VersionWrapper{
		broker:   broker,
		channels: make(map[string]*channelsByBroker),
		logger:   DummyLogger{},
	}

	// Execute options
	for _, option := range options {
		option(&vw)
	}

	return &vw
}

// WithVersionWrapperLogger lets add a logger to the Version Wrapper struct.
func WithVersionWrapperLogger(logger Logger) VersionWrapperOption {
	return func(versionWrapper *VersionWrapper) {
		versionWrapper.logger = logger
	}
}

// Publish a message to the broker.
func (vw *VersionWrapper) Publish(ctx context.Context, channel string, mw BrokerMessage) error {
	// Add version to message
	IfContextSetWith(ctx, ContextKeyIsVersion, func(version string) {
		mw.Headers[VersionField] = []byte(version)
	})

	// Send message
	return vw.broker.Publish(ctx, channel, mw)
}

// Subscribe to messages from the broker.
func (vw *VersionWrapper) Subscribe(ctx context.Context, channel string) (chan BrokerMessage, chan interface{}, error) {
	// Set context
	ctx = context.WithValue(ctx, ContextKeyIsMessageDirection, "reception")
	ctx = context.WithValue(ctx, ContextKeyIsChannel, channel)

	// Get version
	var version string
	IfContextSetWith(ctx, ContextKeyIsVersion, func(v string) { version = v })
	if version == "" {
		return nil, nil, ErrNoVersion
	}

	// Lock the channels to avoid conflict
	vw.channelsMutex.Lock()
	defer vw.channelsMutex.Unlock()

	// Check if the broker channel already exists
	brokerChannel, exists := vw.channels[channel]
	if !exists {
		cbb, err := vw.createBrokerChannels(ctx, channel)
		if err != nil {
			return nil, nil, err
		}
		defer cbb.launchListener(ctx)
		brokerChannel = cbb
	}

	// Check if the version already exists
	cbv, err := brokerChannel.createVersionListener(version)

	return cbv.messages, cbv.cancel, err
}

func (vw *VersionWrapper) createBrokerChannels(ctx context.Context, channel string) (*channelsByBroker, error) {
	// Subscribe to broker
	messages, cancel, err := vw.broker.Subscribe(ctx, channel)
	if err != nil {
		return nil, err
	}

	// Add channels from broker to brokerChannels
	cbb := newChannelsFromBroker(channel, messages, cancel, vw)
	vw.channels[channel] = &cbb // Already locked in parent function

	return &cbb, nil
}
