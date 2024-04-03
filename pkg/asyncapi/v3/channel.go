package asyncapiv3

import (
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/utils/template"
)

const (
	// ChannelSuffix is the suffix added to the channels name.
	ChannelSuffix = "Channel"
)

var (
	// ErrNoMessageInChannel is the error returned when there is no message in a channel.
	ErrNoMessageInChannel = fmt.Errorf("%w: no message in channel", extensions.ErrAsyncAPI)
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
func (ch *Channel) Process(name string, spec Specification) error {
	// Prevent modification if nil
	if ch == nil {
		return nil
	}

	// Set name
	ch.Name = template.Namify(name)

	// Process reference
	if err := ch.processReference(spec); err != nil {
		return err
	}

	// Process messages
	if err := ch.processMessages(spec); err != nil {
		return err
	}

	// Process servers
	if err := ch.processServers(spec); err != nil {
		return err
	}

	// Process parameters
	if err := ch.processParameters(spec); err != nil {
		return err
	}

	// Process tags
	if err := ch.processTags(spec); err != nil {
		return err
	}

	// Process external documentation
	if err := ch.ExternalDocs.Process(ch.Name+ExternalDocsNameSuffix, spec); err != nil {
		return err
	}

	// Process Bindings
	return ch.Bindings.Process(ch.Name+BindingsSuffix, spec)
}

func (ch *Channel) processParameters(spec Specification) error {
	for name, param := range ch.Parameters {
		if err := param.Process(name+"Parameter", spec); err != nil {
			return err
		}
	}

	return nil
}

func (ch *Channel) processServers(spec Specification) error {
	for i, srv := range ch.Servers {
		if err := srv.Process(fmt.Sprintf("%sServer%d", ch.Name, i), spec); err != nil {
			return err
		}
	}

	return nil
}

func (ch *Channel) processMessages(spec Specification) error {
	for name, msg := range ch.Messages {
		if err := msg.Process(name+"Message", spec); err != nil {
			return err
		}
	}

	return nil
}

func (ch *Channel) processTags(spec Specification) error {
	for i, t := range ch.Tags {
		if err := t.Process(fmt.Sprintf("%sTag%d", ch.Name, i), spec); err != nil {
			return err
		}
	}

	return nil
}

func (ch *Channel) processReference(spec Specification) error {
	if ch.Reference == "" {
		return nil
	}

	// Add pointer to reference if there is one
	refTo, err := spec.ReferenceChannel(ch.Reference)
	if err != nil {
		return err
	}
	ch.ReferenceTo = refTo

	return nil
}

// Follow returns referenced channel if specified or the actual channel.
func (ch *Channel) Follow() *Channel {
	if ch.ReferenceTo != nil {
		return ch.ReferenceTo
	}
	return ch
}

// GetMessage will return the channel message.
func (ch Channel) GetMessage() (*Message, error) {
	for _, m := range ch.Follow().Messages {
		return m.Follow(), nil // TODO: change
	}
	return nil, fmt.Errorf("%w: channel %q", ErrNoMessageInChannel, ch.Name)
}
