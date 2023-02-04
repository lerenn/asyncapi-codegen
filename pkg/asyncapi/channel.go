package asyncapi

type Channel struct {
	Subscribe *Operation `json:"subscribe"`
	Publish   *Operation `json:"publish"`

	// Non AsyncAPI fields
	Name string `json:"-"`
}

func (c Channel) GetMessageWithoutReferenceRedirect() Message {
	if c.Subscribe != nil {
		return c.Subscribe.Message
	}

	return c.Publish.Message
}

func (c Channel) CorrelationIDLocation(spec Specification) string {
	msg := c.GetMessageWithoutReferenceRedirect()

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
