package asyncapiv3

import (
	"github.com/lerenn/asyncapi-codegen/pkg/utils"
)

// Channel is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v2.6.0#channelItemObject
type Channel struct {
	Parameters map[string]*Parameter `json:"parameters"`
	Address    string                `json:"address"`
	Messages   map[string]*Message   `json:"messages"`
	Reference  string                `json:"$ref"`

	// Non AsyncAPI fields
	Name        string   `json:"-"`
	Path        string   `json:"-"`
	ReferenceTo *Channel `json:"-"`
}

// Process processes the Channel to make it ready for code generation.
func (ch *Channel) Process(path string, spec Specification) {
	// Set channel name and path
	ch.Name = utils.UpperFirstLetter(path)
	ch.Path = path

	// Process messages
	for name, msg := range ch.Messages {
		msg.Process(name, spec)
	}

	// Process parameters
	for name, parameter := range ch.Parameters {
		parameter.Process(name, spec)
	}

	// Add pointer to reference if there is one
	if ch.Reference != "" {
		ch.ReferenceTo = spec.ReferenceChannel(ch.Reference)
	}
}

// Follow returns referenced channel if specified or the actual channel.
func (ch *Channel) Follow() *Channel {
	if ch.ReferenceTo != nil {
		return ch.ReferenceTo
	}
	return ch
}

// GetMessage will return the channel message.
func (ch Channel) GetMessage() *Message {
	for _, m := range ch.Follow().Messages {
		return m.Follow() // TODO: change
	}
	return nil
}
