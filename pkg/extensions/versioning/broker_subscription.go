package versioning

import (
	"context"
	"fmt"
	"sync"

	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions"
)

type brokerSubscription struct {
	channel      string
	subscription extensions.BrokerChannelSubscription
	parent       *Wrapper

	versionsChannels map[string]versionSubcription
	versionsMutex    sync.Mutex
}

func newBrokerSubscription(
	channel string,
	sub extensions.BrokerChannelSubscription,
	parent *Wrapper,
) brokerSubscription {
	return brokerSubscription{
		channel:          channel,
		subscription:     sub,
		parent:           parent,
		versionsChannels: make(map[string]versionSubcription),
	}
}

func (bs *brokerSubscription) createVersionListener(ctx context.Context, version string) (versionSubcription, error) {
	// Lock the versions to avoid conflict
	bs.versionsMutex.Lock()
	defer bs.versionsMutex.Unlock()

	// Check if the version doesn't exist already
	_, exists := bs.versionsChannels[version]
	if exists {
		return versionSubcription{}, extensions.ErrAlreadySubscribedChannel
	}

	// Create the channels necessary
	cbv := newVersionSubscription(version, bs)
	bs.versionsChannels[version] = cbv
	defer cbv.launchListener(ctx)

	return cbv, nil
}

func (bs *brokerSubscription) removeVersionListener(ctx context.Context, vs *versionSubcription) {
	// Lock the versions to avoid conflict
	bs.versionsMutex.Lock()
	defer bs.versionsMutex.Unlock()

	// Remove the version from the channelsByBroker
	delete(bs.versionsChannels, vs.version)

	// Lock the channels to avoid conflict
	bs.parent.channelsMutex.Lock()
	defer bs.parent.channelsMutex.Unlock()

	// If there is still version channels, do nothing
	if len(bs.versionsChannels) > 0 {
		return
	}

	// Otherwise cancel the broker listener
	bs.subscription.Cancel(ctx)

	// Then delete the channelsByBroker from the Version Switch Wrapper
	delete(bs.parent.channels, bs.channel)
}

func (bs *brokerSubscription) launchListener(ctx context.Context) {
	go func() {
		for {
			// Wait for new messages
			msg, open := <-bs.subscription.MessagesChannel()
			if !open {
				break
			}

			// Get the version from the message
			bVersion, exists := msg.Headers[bs.parent.versionHeaderKey]
			version := string(bVersion)

			// Add default version if none is specified
			if !exists || version == "" {
				// If there is a default version activated, then go on with it
				if bs.parent.defaultVersion != nil {
					version = *bs.parent.defaultVersion
				} else {
					ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, msg)
					bs.parent.logger.Error(ctx, "no version in the message and no default version")
					continue
				}
			}

			// Lock the versions to avoid conflict
			bs.versionsMutex.Lock()

			// Get the correct channel based on the version
			vc, exists := bs.versionsChannels[version]
			if !exists {
				// Set context
				ctx = context.WithValue(ctx, extensions.ContextKeyIsBrokerMessage, msg)
				ctx = context.WithValue(ctx, extensions.ContextKeyIsVersion, version)

				// Log the error
				bs.parent.logger.Error(ctx, fmt.Sprintf("version %q is not registered", version))
				continue
			}

			// Unlock the versions
			bs.versionsMutex.Unlock()

			// Send the message to the correct channel
			vc.subscription.TransmitReceivedMessage(msg)
		}
	}()
}
