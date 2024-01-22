package brokers

const (
	// DefaultQueueGroupID is the default queue name used by brokers.
	// Note: empty in order to avoid using a queue group ID when not expected.
	DefaultQueueGroupID = ""

	// BrokerMessagesQueueSize is the size of the broker messages queue that
	// will hold the messages processed from the broker to the universal format.
	BrokerMessagesQueueSize = 64
)
