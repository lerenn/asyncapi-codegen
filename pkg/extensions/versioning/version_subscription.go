package versioning

import (
	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers"
)

type versionSubcription struct {
	version  string
	messages chan extensions.BrokerMessage
	cancel   chan any
	parent   *brokerSubscription
}

func newVersionSubscription(version string, parent *brokerSubscription) versionSubcription {
	return versionSubcription{
		version:  version,
		messages: make(chan extensions.BrokerMessage, brokers.BrokerMessagesQueueSize),
		cancel:   make(chan any, 1),
		parent:   parent,
	}
}

func (vs *versionSubcription) launchListener() {
	go func() {
		// Wait to receive cancel
		<-vs.cancel

		// When cancel is received, then remove version listener
		vs.parent.removeVersionListener(vs)
	}()
}

func (vs *versionSubcription) closeChannels() {
	// Receiving no more messages
	close(vs.messages)

	// Closing cancel channel to let caller knows that everything is cleaned up
	close(vs.cancel)
}
