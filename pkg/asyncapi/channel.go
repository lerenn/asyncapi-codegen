package asyncapi

import (
	"strings"

	"github.com/lerenn/asyncapi-codegen/pkg/utils"
)

// Channel is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v2.6.0#channelItemObject
type Channel struct {
	Parameters map[string]*Parameter `json:"parameters"`

	Subscribe *Operation `json:"subscribe"`
	Publish   *Operation `json:"publish"`

	// Non AsyncAPI fields
	Name string `json:"-"`
	Path string `json:"-"`
}

// Process processes the Channel to make it ready for code generation.
func (c *Channel) Process(path string, spec Specification) {
	// Set channel name and path
	c.Name = utils.UpperFirstLetter(path)
	c.Path = path

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

	// Process parameters
	for n, p := range c.Parameters {
		p.Process(n, spec)
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
