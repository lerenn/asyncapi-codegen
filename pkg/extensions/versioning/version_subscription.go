package versioning

import (
	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers"
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
			make(chan extensions.BrokerMessage, brokers.BrokerMessagesQueueSize),
			make(chan any, 1),
		),
		parent: parent,
	}
}

func (vs *versionSubcription) launchListener() {
	// Wait for cancellation and remove version listener when it happens
	vs.subscription.WaitForCancellationAsync(func() {
		vs.parent.removeVersionListener(vs)
	})
}
