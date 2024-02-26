package asyncapiv3

import (
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/utils"
)

// Channel is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#channelObject
type Channel struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Address      string                 `json:"address"`
	Messages     map[string]*Message    `json:"messages"`
	Title        string                 `json:"title"`
	Summary      string                 `json:"summary"`
	Description  string                 `json:"description"`
	Servers      []*Server              `json:"servers"`
	Parameters   map[string]*Parameter  `json:"parameters"`
	Tags         []*Tag                 `json:"tags"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs"`
	Bindings     *ChannelBindings       `json:"bindings"`
	Reference    string                 `json:"$ref"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string   `json:"-"`
	ReferenceTo *Channel `json:"-"`
}

// Process processes the Channel to make it ready for code generation.
func (ch *Channel) Process(path string, spec Specification) {
	// Prevent modification if nil
	if ch == nil {
		return
	}

	// Set name
	ch.Name = utils.UpperFirstLetter(path)

	// Add pointer to reference if there is one
	if ch.Reference != "" {
		ch.ReferenceTo = spec.ReferenceChannel(ch.Reference)
	}

	// Process messages
	for name, msg := range ch.Messages {
		msg.Process(name, spec)
	}

	// Process servers
	for i, srv := range ch.Servers {
		srv.Process(fmt.Sprintf("%sServer%d", ch.Name, i), spec)
	}

	// Process parameters
	for name, parameter := range ch.Parameters {
		parameter.Process(name, spec)
	}

	// Process tags
	for i, t := range ch.Tags {
		t.Process(fmt.Sprintf("%sTag%d", ch.Name, i), spec)
	}

	// Process external documentation
	ch.ExternalDocs.Process(ch.Name+ExternalDocsNameSuffix, spec)

	// Process Bindings
	ch.Bindings.Process(ch.Name+BindingsSuffix, spec)
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
