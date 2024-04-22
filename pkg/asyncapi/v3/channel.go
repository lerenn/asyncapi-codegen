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

// generateMetadata generates metadata for the Channel.
func (ch *Channel) generateMetadata(name string) error {
	// Prevent modification if nil
	if ch == nil {
		return nil
	}

	// Set name
	ch.Name = template.Namify(name)

	// Generate messages metadata
	for name, msg := range ch.Messages {
		if err := msg.generateMetadata(name + "Message"); err != nil {
			return err
		}
	}

	// Generate servers metadata
	for i, srv := range ch.Servers {
		srv.generateMetadata(fmt.Sprintf("%sServer%d", ch.Name, i))
	}

	// Generate parameters metadata
	for name, param := range ch.Parameters {
		param.generateMetadata(name + "Parameter")
	}

	// Generate tags metadata
	for i, t := range ch.Tags {
		t.generateMetadata(fmt.Sprintf("%sTag%d", ch.Name, i))
	}

	// Generate external documentation metadata
	ch.ExternalDocs.generateMetadata(ch.Name + ExternalDocsNameSuffix)

	// Generate Bindings metadata
	ch.Bindings.generateMetadata(ch.Name + BindingsSuffix)
	return nil
}

// setDependencies sets dependencies between the different elements of the Channel.
//
//nolint:cyclop
func (ch *Channel) setDependencies(spec Specification) error {
	// Prevent modification if nil
	if ch == nil {
		return nil
	}

	// Set reference
	if err := ch.setReference(spec); err != nil {
		return err
	}

	// Set messages dependencies
	for _, msg := range ch.Messages {
		if err := msg.setDependencies(spec); err != nil {
			return err
		}
	}

	// Set servers dependencies
	for _, srv := range ch.Servers {
		if err := srv.setDependencies(spec); err != nil {
			return err
		}
	}

	// Set parameters dependencies
	for _, param := range ch.Parameters {
		if err := param.setDependencies(spec); err != nil {
			return err
		}
	}

	// Set tags dependencies
	for _, t := range ch.Tags {
		if err := t.setDependencies(spec); err != nil {
			return err
		}
	}

	// Set external documentation dependencies
	if err := ch.ExternalDocs.setDependencies(spec); err != nil {
		return err
	}

	// Set Bindings dependencies
	return ch.Bindings.setDependencies(spec)
}

func (ch *Channel) setReference(spec Specification) error {
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
