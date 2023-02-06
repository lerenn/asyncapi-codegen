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

	msg := c.GetChannelMessage()
	msg.Process(spec)
}

func (c *Channel) setMapsValuesName() {
	msg := c.GetChannelMessage()

	if msg.Reference != "" {
		msg.Name = strings.Split(msg.Reference, "/")[3]
	} else {
		msg.Name = c.Name
	}
}

// GetChannelMessage will return the channel message
// WARNING: if there is a reference, then it won't be followed.
func (c Channel) GetChannelMessage() *Message {
	if c.Subscribe != nil {
		return &c.Subscribe.Message
	}

	return &c.Publish.Message
}
