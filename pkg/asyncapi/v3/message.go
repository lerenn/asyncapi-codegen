package asyncapiv3

import (
	"fmt"
	"strings"

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

// MessageType is a structure that represents the type of a field.
type MessageType string

// String returns the string representation of the type.
func (t MessageType) String() string {
	return string(t)
}

const (
	// MessageTypeIsArray represents the type of an array.
	MessageTypeIsArray MessageType = "array"
	// MessageTypeIsHeader represents the type of a header.
	MessageTypeIsHeader MessageType = "header"
	// MessageTypeIsObject represents the type of an object.
	MessageTypeIsObject MessageType = "object"
	// MessageTypeIsString represents the type of a string.
	MessageTypeIsString MessageType = "string"
	// MessageTypeIsInteger represents the type of an integer.
	MessageTypeIsInteger MessageType = "integer"
	// MessageTypeIsPayload represents the type of a payload.
	MessageTypeIsPayload MessageType = "payload"
)

// Message is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#messageObject
type Message struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Headers       *Schema                `json:"headers"`
	Payload       *Schema                `json:"payload"`
	OneOf         []*Message             `json:"oneOf"`
	CorrelationID *CorrelationID         `json:"correlationID"`
	ContentType   string                 `json:"contentType"`
	Name          string                 `json:"name"`
	Title         string                 `json:"title"`
	Summary       string                 `json:"summary"`
	Description   string                 `json:"description"`
	Tags          []*Tag                 `json:"tags"`
	ExternalDocs  *ExternalDocumentation `json:"externalDocs"`
	Bindings      *MessageBindings       `json:"bindings"`
	Examples      []*MessageExample      `json:"examples"`
	Traits        []*MessageTrait        `json:"traits"`
	Reference     string                 `json:"$ref"`

	// --- Non AsyncAPI fields -------------------------------------------------

	ReferenceTo *Message `json:"-"`

	// CorrelationIDLocation will indicate where the correlation id is
	// According to: https://www.asyncapi.com/docs/reference/specification/v3.0.0#correlationIdObject
	CorrelationIDRequired bool `json:"-"`
}

// Process processes the Message to make it ready for code generation.
func (msg *Message) Process(name string, spec Specification) {
	// Prevent modification if nil
	if msg == nil {
		return
	}

	// Set name
	if msg.Name == "" {
		msg.Name = utils.UpperFirstLetter(name)
	} else {
		msg.Name = utils.UpperFirstLetter(msg.Name)
	}

	// Add pointer to reference if there is one
	if msg.Reference != "" {
		msg.ReferenceTo = spec.ReferenceMessage(msg.Reference)
	}

	// Process Headers and Payload
	msg.Headers.Process(name+"Headers", spec, false)
	msg.Payload.Process(name+"Payload", spec, false)

	// Process OneOf
	for i, v := range msg.OneOf {
		// Process the OneOf
		v.Process(fmt.Sprintf("%sOneOf%d", name, i), spec)

		// Merge the OneOf as one payload
		msg.MergeWith(spec, *v)
	}

	// Process tags
	for i, t := range msg.Tags {
		t.Process(fmt.Sprintf("%sTag%d", msg.Name, i), spec)
	}

	// Process external documentation
	msg.ExternalDocs.Process(msg.Name+ExternalDocsNameSuffix, spec)

	// Process Bindings
	msg.Bindings.Process(msg.Name+BindingsSuffix, spec)

	// Process Message Examples
	for i, ex := range msg.Examples {
		ex.Process(fmt.Sprintf("%sExample%d", msg.Name, i), spec)
	}

	// Process traits and apply them
	for i, t := range msg.Traits {
		t.Process(fmt.Sprintf("%sTrait%d", msg.Name, i), spec)
		msg.ApplyTrait(t.Follow(), spec)
	}

	// Process correlation ID
	msg.createCorrelationIDFieldIfMissing()
	msg.CorrelationIDRequired = msg.isCorrelationIDRequired()
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
		return msg.createTreeUntilLocationFromMessageType(MessageTypeIsHeader, location)
	}

	// Check that the correlation ID is in payload
	if strings.HasPrefix(location, "$message.payload#") {
		return msg.createTreeUntilLocationFromMessageType(MessageTypeIsPayload, location)
	}

	// Default to nothing
	return utils.ToPointer(NewSchema())
}

func (msg *Message) createTreeUntilLocationFromMessageType(t MessageType, location string) (locationParent *Schema) {
	// Get correct top level placeholder
	var placeholder **Schema
	if t == MessageTypeIsHeader {
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
		(*placeholder).Name = MessageTypeIsHeader.String()
		(*placeholder).Type = MessageTypeIsObject.String()
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
				child.Type = MessageTypeIsString.String()
			} else {
				child.Type = MessageTypeIsHeader.String()
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
	} else if ref[0] == MessageTypeIsHeader.String() {
		next = msg.Headers
	}

	return next.referenceFrom(ref[1:])
}

// MergeWith merges the Message with another one.
func (msg *Message) MergeWith(spec Specification, msg2 Message) {
	// Remove reference if merging
	if msg.Reference != "" {
		refMsg := spec.ReferenceMessage(msg.Reference)
		msg.Reference = ""
		msg.MergeWith(spec, *refMsg)
	}

	// Get reference from msg2
	if msg2.Reference != "" {
		refMsg2 := spec.ReferenceMessage(msg2.Reference)
		msg2.MergeWith(spec, *refMsg2)
	}

	// Merge Payload
	if msg2.Payload != nil {
		if msg.Payload == nil {
			msg.Payload = deepcopy.Copy(msg2.Payload).(*Schema)
		} else {
			msg.Payload.MergeWith(spec, *msg2.Payload)
		}
	}

	// Merge Headers
	if msg2.Headers != nil {
		if msg.Headers == nil {
			msg.Headers = deepcopy.Copy(msg2.Headers).(*Schema)
		} else {
			msg.Headers.MergeWith(spec, *msg2.Headers)
		}
	}
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
func (msg *Message) ApplyTrait(mt *MessageTrait, spec Specification) {
	// Check message is not nil
	if msg == nil {
		return
	}

	// Merge headers if present
	if mt.Headers != nil {
		msg.Headers.MergeWith(spec, *mt.Headers)
	}

	// Merge payload if present
	if mt.Payload != nil {
		msg.Payload.MergeWith(spec, *mt.Payload)
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
}

func (msg Message) HaveCorrelationID() bool {
	return msg.Follow().CorrelationID.Exists()
}
