package brokers

const (
	// DefaultQueueGroupID is the default queue name used by brokers.
	DefaultQueueGroupID = "asyncapi"

	// BrokerMessagesQueueSize is the size of the broker messages queue that
	// will hold the messages processed from the broker to the universal format.
	BrokerMessagesQueueSize = 64
)
