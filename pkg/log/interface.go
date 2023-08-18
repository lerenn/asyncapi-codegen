package log

// Context will contain information about the context where the log is called
// This is interesting if you're implementing your own logger
type Context struct {
	// Module is the name of the module this data is coming from.
	// When coming from a generated client, it is `asyncapi`
	Module string

	// Provider is the name of the provider this data is coming from.
	// When coming from generated code, it is `app`, `client` or `broker`
	Provider string

	// Action is the name of the action this data is coming from.
	// When coming from generated code, it is the name of the channel
	Action string

	// Operation is the name of the operation this data is coming from.
	// When coming from generated code, it is `subscribe`, `publish`, `wait-for`, etc
	Operation string

	// Message is the message that has been sent or received
	Message any

	// CorrelationID is the correlation ID of the message
	CorrelationID string
}

// AdditionalInfo is a key-value pair that will be added to the log
type AdditionalInfo struct {
	Key   string
	Value interface{}
}

// Logger is the interface that must be implemented by a logger
type Logger interface {
	// Info logs information based on a message and key-value elements
	Info(ctx Context, msg string, info ...AdditionalInfo)

	// Error logs error based on a message and key-value elements
	Error(ctx Context, msg string, info ...AdditionalInfo)
}
