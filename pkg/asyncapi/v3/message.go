package asyncapiv3

import (
	"fmt"
	"strings"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/utils"
	"github.com/mohae/deepcopy"
)

// MessageField is a structure that represents the type of a field.
type MessageField string

// String returns the string representation of the type.
func (t MessageField) String() string {
	return string(t)
}

const (
	// MessageFieldIsHeader represents the message field of a header.
	MessageFieldIsHeader MessageField = "header"
	// MessageFieldIsPayload represents the message field of a payload.
	MessageFieldIsPayload MessageField = "payload"
)

const (
	// MessageHeadersSuffix is the suffix for the headers schema name.
	MessageHeadersSuffix = "Headers"
	// MessagePayloadSuffix is the suffix for the payload schema name.
	MessagePayloadSuffix = "Payload"
)

// Message is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#messageObject
type Message struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Headers       *Schema                `json:"headers,omitempty"`
	Payload       *Schema                `json:"payload,omitempty"`
	OneOf         []*Message             `json:"oneOf,omitempty"`
	CorrelationID *CorrelationID         `json:"correlationID,omitempty"`
	ContentType   string                 `json:"contentType,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Title         string                 `json:"title,omitempty"`
	Summary       string                 `json:"summary,omitempty"`
	Description   string                 `json:"description,omitempty"`
	Tags          []*Tag                 `json:"tags,omitempty"`
	ExternalDocs  *ExternalDocumentation `json:"externalDocs,omitempty"`
	Bindings      *MessageBindings       `json:"bindings,omitempty"`
	Examples      []*MessageExample      `json:"examples,omitempty"`
	Traits        []*MessageTrait        `json:"traits,omitempty"`
	Reference     string                 `json:"$ref,omitempty"`

	// --- Non AsyncAPI fields -------------------------------------------------

	ReferenceTo *Message `json:"-"`

	// CorrelationIDLocation will indicate where the correlation id is
	// According to: https://www.asyncapi.com/docs/reference/specification/v3.0.0#correlationIdObject
	CorrelationIDRequired bool `json:"-"`
}

// generateMetadata generates metadata for the Message.
func (msg *Message) generateMetadata(parentName, name string, number *int) error {
	// Prevent modification if nil
	if msg == nil {
		return nil
	}

	// Set name
	msg.Name = generateFullName(parentName, name, "Message", number)

	// Generate Payload metadata
	// NOTE: Suffix Message name with payload for name and set no parent
	if err := msg.Payload.generateMetadata("", msg.Name+"_Payload", nil, false); err != nil {
		return err
	}

	// Generate Headers metadata
	if err := msg.generateHeadersMetadata(); err != nil {
		return err
	}

	// Generate OneOf metadata
	if err := msg.generateOneOfMetadata(); err != nil {
		return err
	}

	// Generate tags metadata
	msg.generateTagsMetadata()

	// Generate external documentation metadata
	msg.ExternalDocs.generateMetadata(msg.Name, ExternalDocsNameSuffix)

	// Generate Bindings metadata
	msg.Bindings.generateMetadata(msg.Name, "")

	// Generate Message Examples metadata
	msg.generateExamplesMetadata()

	// Generate traits metadata
	if err := msg.generateTraitsMetadata(); err != nil {
		return err
	}

	// Process correlation ID
	msg.createCorrelationIDFieldIfMissing()
	msg.CorrelationIDRequired = msg.isCorrelationIDRequired()

	return nil
}

// setDependencies sets dependencies between the different elements of the Message.
func (msg *Message) setDependencies(spec Specification) error {
	// Prevent modification if nil
	if msg == nil {
		return nil
	}

	// Add pointer to reference if there is one
	if err := msg.setReference(spec); err != nil {
		return err
	}

	// Set Payload dependencies
	if err := msg.Payload.setDependencies(spec); err != nil {
		return err
	}

	// Set Headers dependencies
	if err := msg.setHeadersDependencies(spec); err != nil {
		return err
	}

	// Set OneOf dependencies
	if err := msg.setOneOfDependencies(spec); err != nil {
		return err
	}

	// Set tags dependencies
	if err := msg.setTagsDependencies(spec); err != nil {
		return err
	}

	// Set external documentation dependencies
	if err := msg.ExternalDocs.setDependencies(spec); err != nil {
		return err
	}

	// Set Bindings dependencies
	if err := msg.Bindings.setDependencies(spec); err != nil {
		return err
	}

	// Set Message Examples dependencies
	if err := msg.setExamplesDependencies(spec); err != nil {
		return err
	}

	// Set traits dependencies
	return msg.setTraitsDependencies(spec)
}

func (msg *Message) setReference(spec Specification) error {
	if msg.Reference == "" {
		return nil
	}

	refMsg, err := spec.ReferenceMessage(msg.Reference)
	if err != nil {
		return err
	}
	msg.ReferenceTo = refMsg

	return nil
}

func (msg *Message) generateTagsMetadata() {
	for i, t := range msg.Tags {
		t.generateMetadata(msg.Name, "", &i)
	}
}

func (msg *Message) setTagsDependencies(spec Specification) error {
	for _, t := range msg.Tags {
		if err := t.setDependencies(spec); err != nil {
			return err
		}
	}

	return nil
}

func (msg *Message) generateExamplesMetadata() {
	for i, ex := range msg.Examples {
		ex.generateMetadata(msg.Name, "", &i)
	}
}

func (msg *Message) setExamplesDependencies(spec Specification) error {
	for _, ex := range msg.Examples {
		if err := ex.setDependencies(spec); err != nil {
			return err
		}
	}

	return nil
}

func (msg *Message) generateHeadersMetadata() error {
	if msg.Headers == nil {
		return nil
	}

	if err := msg.Headers.generateMetadata(msg.Name, "Headers", nil, false); err != nil {
		return err
	}

	if msg.Headers.Follow().Type != SchemaTypeIsObject.String() {
		return fmt.Errorf(
			"%w: %q headers must be an object, is %q",
			extensions.ErrAsyncAPI, msg.Name, msg.Headers.Follow().Type)
	}

	return nil
}

func (msg *Message) setHeadersDependencies(spec Specification) error {
	if msg.Headers == nil {
		return nil
	}

	return msg.Headers.setDependencies(spec)
}

func (msg *Message) generateOneOfMetadata() error {
	for i, v := range msg.OneOf {
		if err := v.generateMetadata(msg.Name, "", &i); err != nil {
			return err
		}
	}

	return nil
}

func (msg *Message) setOneOfDependencies(spec Specification) error {
	for _, v := range msg.OneOf {
		// Set OneOf dependencies
		if err := v.setDependencies(spec); err != nil {
			return err
		}

		// Merge the OneOf as one payload
		if err := msg.MergeWith(spec, *v); err != nil {
			return err
		}
	}

	return nil
}

func (msg *Message) generateTraitsMetadata() error {
	for i, t := range msg.Traits {
		if err := t.generateMetadata(msg.Name, "", &i); err != nil {
			return err
		}
	}

	return nil
}

func (msg *Message) setTraitsDependencies(spec Specification) error {
	for _, t := range msg.Traits {
		if err := t.setDependencies(spec); err != nil {
			return err
		}
		if err := msg.ApplyTrait(t.Follow(), spec); err != nil {
			return err
		}
	}

	return nil
}

func (msg Message) isCorrelationIDRequired() bool {
	if msg.CorrelationID == nil || msg.CorrelationID.Location == "" {
		return false
	}

	correlationIDParent := msg.createTreeUntilLocation(msg.CorrelationID.Location)
	path := strings.Split(msg.CorrelationID.Location, "/")
	return correlationIDParent.IsFieldRequired(path[len(path)-1])
}

func (msg *Message) createCorrelationIDFieldIfMissing() {
	if msg.CorrelationID == nil {
		return
	}

	_ = msg.createTreeUntilLocation(msg.CorrelationID.Location)
}

func (msg *Message) createTreeUntilLocation(location string) (locationParent *Schema) {
	// Check location
	if location == "" {
		return utils.ToPointer(NewSchema())
	}

	// Check that the correlation ID is in header
	if strings.HasPrefix(location, "$message.header#") {
		return msg.createTreeUntilLocationFromMessageType(MessageFieldIsHeader, location)
	}

	// Check that the correlation ID is in payload
	if strings.HasPrefix(location, "$message.payload#") {
		return msg.createTreeUntilLocationFromMessageType(MessageFieldIsPayload, location)
	}

	// Default to nothing
	return utils.ToPointer(NewSchema())
}

func (msg *Message) createTreeUntilLocationFromMessageType(t MessageField, location string) (locationParent *Schema) {
	// Get correct top level placeholder
	var placeholder **Schema
	if t == MessageFieldIsHeader {
		placeholder = &msg.Headers
	} else {
		placeholder = &msg.Payload
	}

	var child *Schema
	switch {
	case (*placeholder) != nil && (*placeholder).ReferenceTo != nil: // If there is a reference
		// Use it as child
		child = (*placeholder).ReferenceTo
	case (*placeholder) == nil: // If there is no header and no reference
		// Create a default one for the message
		(*placeholder) = utils.ToPointer(NewSchema())
		(*placeholder).Name = MessageFieldIsHeader.String()
		(*placeholder).Type = SchemaTypeIsObject.String()
		fallthrough
	default:
		// Set the child as the message headers
		child = (*placeholder)
	}

	// Go down the path to correlation ID
	return msg.downToLocation(child, location)
}

func (msg Message) downToLocation(child *Schema, location string) (locationParent *Schema) {
	var exists bool

	path := strings.Split(location, "/")
	for i, v := range path[1:] {
		// Keep the parent
		locationParent = child

		// Get the corresponding child
		child, exists = locationParent.Properties[v]
		if !exists { // If it doesn't exist
			// Create child
			child = utils.ToPointer(NewSchema())
			child.Name = v
			if i == len(path)-2 { // As there is -1 in the loop slice
				child.Type = SchemaTypeIsString.String()
			} else {
				child.Type = MessageFieldIsHeader.String()
			}

			// Add it to parent
			if locationParent.Properties == nil {
				locationParent.Properties = make(map[string]*Schema)
			}
			locationParent.Properties[v] = child
		}
	}

	return locationParent
}

func (msg *Message) referenceFrom(ref []string) any {
	if len(ref) == 0 {
		return msg
	}

	var next *Schema
	if ref[0] == "payload" {
		next = msg.Payload
	} else if ref[0] == MessageFieldIsHeader.String() {
		next = msg.Headers
	}

	return next.referenceFrom(ref[1:])
}

// MergeWith merges the Message with another one.
func (msg *Message) MergeWith(spec Specification, msg2 Message) error {
	// Remove reference if merging
	if err := msg.mergeWithSelfReference(spec); err != nil {
		return err
	}

	// Get reference from msg2
	if err := msg.mergeWithMessageFromReference(msg2, spec); err != nil {
		return err
	}

	// Merge Payload
	if err := msg.mergeWithMessagePayload(msg2, spec); err != nil {
		return err
	}

	// Merge Headers
	if err := msg.mergeWithMessageHeaders(msg2, spec); err != nil {
		return err
	}

	return nil
}

func (msg *Message) mergeWithSelfReference(spec Specification) error {
	if msg.Reference != "" {
		// Get referenced message
		refMsg, err := spec.ReferenceMessage(msg.Reference)
		if err != nil {
			return err
		}

		// Remove reference
		msg.Reference = ""

		// Merge with referenced message
		if err := msg.MergeWith(spec, *refMsg); err != nil {
			return err
		}
	}

	return nil
}

func (msg *Message) mergeWithMessageFromReference(msg2 Message, spec Specification) error {
	if msg2.Reference != "" {
		// Get referenced message
		refMsg2, err := spec.ReferenceMessage(msg2.Reference)
		if err != nil {
			return err
		}

		// Merge with referenced message
		if err := msg2.MergeWith(spec, *refMsg2); err != nil {
			return err
		}
	}

	return nil
}

func (msg *Message) mergeWithMessagePayload(msg2 Message, spec Specification) error {
	if msg2.Payload != nil {
		if msg.Payload == nil {
			msg.Payload = deepcopy.Copy(msg2.Payload).(*Schema)
		} else {
			if err := msg.Payload.MergeWith(spec, *msg2.Payload); err != nil {
				return err
			}
		}
	}

	return nil
}

func (msg *Message) mergeWithMessageHeaders(msg2 Message, spec Specification) error {
	if msg2.Headers != nil {
		if msg.Headers == nil {
			msg.Headers = deepcopy.Copy(msg2.Headers).(*Schema)
		} else {
			if err := msg.Headers.MergeWith(spec, *msg2.Headers); err != nil {
				return err
			}
		}
	}
	return nil
}

// Follow returns referenced message if specified or the actual message.
func (msg *Message) Follow() *Message {
	if msg.ReferenceTo != nil {
		return msg.ReferenceTo
	}
	return msg
}

// ApplyTrait applies a trait to the message.
//
//nolint:cyclop
func (msg *Message) ApplyTrait(mt *MessageTrait, spec Specification) error {
	// Check message is not nil
	if msg == nil {
		return nil
	}

	// Merge headers
	if err := msg.mergeHeaders(spec, mt.Headers); err != nil {
		return err
	}

	// Merge payload
	if err := msg.mergePayload(spec, mt.Payload); err != nil {
		return err
	}

	// Add correlation ID if present and not overriding
	if msg.CorrelationID == nil && mt.CorrelationID != nil {
		corelID := *mt.CorrelationID
		msg.CorrelationID = &corelID
	}

	// Override content type if not set
	if msg.ContentType == "" {
		msg.ContentType = mt.ContentType
	}

	// Override title if not set
	if msg.Title == "" {
		msg.Title = mt.Title
	}

	// Override summary if not set
	if msg.Summary == "" {
		msg.Summary = mt.Summary
	}

	// Override description if not set
	if msg.Description == "" {
		msg.Description = mt.Description
	}

	// Merge tags
	msg.Tags = append(msg.Tags, mt.Tags...)
	msg.Tags = RemoveDuplicateTags(msg.Tags)

	// Override external docs if not set
	if msg.ExternalDocs == nil && mt.ExternalDocs != nil {
		extDoc := *mt.ExternalDocs
		msg.ExternalDocs = &extDoc
	}

	// Merge examples
	msg.Examples = append(msg.Examples, mt.Examples...)

	return nil
}

func (msg *Message) mergeHeaders(spec Specification, headers *Schema) error {
	// Check if headers are nil
	if headers == nil {
		return nil
	}

	// Check if message headers are nil, then create them
	if msg.Headers == nil {
		newHeaders := utils.ToValue(headers.Follow())
		msg.Headers = &newHeaders
		if err := newHeaders.generateMetadata(msg.Name, "Headers", nil, false); err != nil {
			return err
		}
		return newHeaders.setDependencies(spec)
	}

	// Merge headers
	return msg.Headers.MergeWith(spec, *headers)
}

func (msg *Message) mergePayload(spec Specification, payload *Schema) error {
	// Check if payload is nil
	if payload == nil {
		return nil
	}

	// Check if message payload is nil, then create them
	if msg.Payload == nil {
		newPayload := utils.ToValue(payload.Follow())
		msg.Payload = &newPayload
		if err := newPayload.generateMetadata(msg.Name, "Payload", nil, false); err != nil {
			return err
		}
		return newPayload.setDependencies(spec)
	}

	// Merge payload
	return msg.Headers.MergeWith(spec, *payload)
}

// HaveCorrelationID check that the message have a correlation ID.
func (msg Message) HaveCorrelationID() bool {
	return msg.Follow().CorrelationID.Exists()
}
