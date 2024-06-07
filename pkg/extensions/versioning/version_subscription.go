package versioning

import (
	"context"
	"time"

	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions"
	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions/brokers"
)

type versionSubcription struct {
	version      string
	subscription extensions.BrokerChannelSubscription
	parent       *brokerSubscription
}

func newVersionSubscription(version string, parent *brokerSubscription) versionSubcription {
	return versionSubcription{
		version: version,
		subscription: extensions.NewBrokerChannelSubscription(
			make(chan extensions.AcknowledgeableBrokerMessage, brokers.BrokerMessagesQueueSize),
			make(chan any, 1),
		),
		parent: parent,
	}
}

func (vs *versionSubcription) launchListener(ctx context.Context) {
	// Wait for cancellation and remove version listener when it happens
	vs.subscription.WaitForCancellationAsync(func() {
		// Create cancel function in case there is a problem with broker removal
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		// Remove the version listener
		vs.parent.removeVersionListener(ctx, vs)
	})
}
