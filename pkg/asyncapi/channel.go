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

func (c Channel) CorrelationIdLocation(spec Specification) string {
	msg := c.GetMessageWithoutReferenceRedirect()

	// Let's check the message before the reference
	if msg.CorrelationId != nil {
		return msg.CorrelationId.Location
	}

	// If there is a reference, check it
	if msg.Reference != "" {
		correlationId := spec.ReferenceMessage(msg.Reference).CorrelationId
		if correlationId != nil {
			return correlationId.Location
		}
	}

	return ""
}
