package broker

// Message is a wrapper that will contain all information regarding a message
type Message struct {
	CorrelationID *string
	Payload       []byte
}
