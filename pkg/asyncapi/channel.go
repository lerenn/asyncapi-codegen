package asyncapi

import "strings"

type Channel struct {
	Subscribe *Operation `json:"subscribe"`
	Publish   *Operation `json:"publish"`

	// Non AsyncAPI fields
	Name string `json:"-"`
}

func (c *Channel) Process(spec Specification) {
	c.setMapsValuesName()

	msg := c.GetMessageWithoutReferenceRedirect()
	msg.Process(spec)
}

func (c *Channel) setMapsValuesName() {
	msg := c.GetMessageWithoutReferenceRedirect()

	if msg.Reference != "" {
		msg.Name = strings.Split(msg.Reference, "/")[3]
	} else {
		msg.Name = c.Name
	}
}

func (c Channel) GetMessageWithoutReferenceRedirect() *Message {
	if c.Subscribe != nil {
		return &c.Subscribe.Message
	}

	return &c.Publish.Message
}
