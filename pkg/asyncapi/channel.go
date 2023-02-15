package asyncapi

import "strings"

type Channel struct {
	Subscribe *Operation `json:"subscribe"`
	Publish   *Operation `json:"publish"`

	// Non AsyncAPI fields
	Name string `json:"-"`
}

func (c *Channel) Process(name string, spec Specification) {
	// Set channel name
	c.Name = name

	// Get message
	msg := c.GetChannelMessage()

	// Get message name
	var msgName string
	if msg.Reference != "" {
		msgName = strings.Split(msg.Reference, "/")[3]
	} else {
		msgName = c.Name
	}

	// Process message
	msg.Process(msgName, spec)
}

// GetChannelMessage will return the channel message
// WARNING: if there is a reference, then it won't be followed.
func (c Channel) GetChannelMessage() *Message {
	if c.Subscribe != nil {
		return &c.Subscribe.Message
	}

	return &c.Publish.Message
}
