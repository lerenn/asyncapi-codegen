package asyncapiv3

import (
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
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

	// NOTE: the JSON null literal is a valid value for address, which cannot be parsed as a Go string
	// Solutions are to change to a *string (a breaking change) or use a wrapper type like sql.NullString
	Address      string                 `json:"address,omitempty"`
	Messages     map[string]*Message    `json:"messages,omitempty"`
	Title        string                 `json:"title,omitempty"`
	Summary      string                 `json:"summary,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Servers      []*Server              `json:"servers,omitempty"` // Reference only
	Parameters   map[string]*Parameter  `json:"parameters,omitempty"`
	Tags         []*Tag                 `json:"tags,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
	Bindings     *ChannelBindings       `json:"bindings,omitempty"`
	Reference    string                 `json:"$ref,omitempty"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string   `json:"-"`
	ReferenceTo *Channel `json:"-"`
}

// generateMetadata generates metadata for the Channel.
// It generates the name of the channel and the metadata of its elements.
func (ch *Channel) generateMetadata(parentName, name string) error {
	// Prevent modification if nil
	if ch == nil {
		return nil
	}

	// Set name
	ch.Name = generateFullName(parentName, name, ChannelSuffix, nil)

	// Generate messages metadata
	for name, msg := range ch.Messages {
		if err := msg.generateMetadata(ch.Name, name, nil); err != nil {
			return err
		}
	}

	// Generate parameters metadata
	for name, param := range ch.Parameters {
		param.generateMetadata(ch.Name, name)
	}

	// Generate tags metadata
	for i, t := range ch.Tags {
		t.generateMetadata(ch.Name, "", &i)
	}

	// Generate external documentation metadata
	fullname := generateFullName(ch.Name, "", ExternalDocsNameSuffix, nil)
	ch.ExternalDocs.generateMetadata(ch.Name, fullname)

	// Generate Bindings metadata
	ch.Bindings.generateMetadata(ch.Name, "")

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
