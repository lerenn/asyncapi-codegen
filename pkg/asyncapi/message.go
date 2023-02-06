package asyncapi

type Message struct {
	Description   string         `json:"description"`
	Headers       *Any           `json:"headers"`
	Payload       *Any           `json:"payload"`
	CorrelationID *CorrelationID `json:"correlationID"`
	Reference     string         `json:"$ref"`

	// --- Non AsyncAPI fields -------------------------------------------------
	Name string `json:"-"`

	// CorrelationIDLocation will indicate where the correlation id is
	// According to: https://www.asyncapi.com/docs/reference/specification/v2.5.0#correlationIDObject
	CorrelationIDLocation string `json:"-"`
}

func (msg *Message) Process(spec Specification) {
	msg.CorrelationIDLocation = msg.getCorrelationIDLocation(spec)
}

func (msg Message) getCorrelationIDLocation(spec Specification) string {
	// Let's check the message before the reference
	if msg.CorrelationID != nil {
		return msg.CorrelationID.Location
	}

	// If there is a reference, check it
	if msg.Reference != "" {
		correlationID := spec.ReferenceMessage(msg.Reference).CorrelationID
		if correlationID != nil {
			return correlationID.Location
		}
	}

	return ""
}
