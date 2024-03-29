package asyncapiv2

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/utils"
	"github.com/lerenn/asyncapi-codegen/pkg/utils/template"
	"github.com/mohae/deepcopy"
)

const (
	// MessageSuffix is the suffix added to the name of the message.
	MessageSuffix = "Message"
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

// Message is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v2.6.0#messageObject
type Message struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Description   string         `json:"description"`
	Headers       *Schema        `json:"headers"`
	OneOf         []*Message     `json:"oneOf"`
	Payload       *Schema        `json:"payload"`
	CorrelationID *CorrelationID `json:"correlationID"`
	Reference     string         `json:"$ref"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string   `json:"-"`
	ReferenceTo *Message `json:"-"`

	// CorrelationIDLocation will indicate where the correlation id is
	// According to: https://www.asyncapi.com/docs/reference/specification/v2.6.0#correlationIDObject
	CorrelationIDLocation string `json:"-"`
	CorrelationIDRequired bool   `json:"-"`
}

// Process processes the Message to make it ready for code generation.
func (msg *Message) Process(name string, spec Specification) error {
	msg.Name = template.Namify(name)

	// Add pointer to reference if there is one
	if msg.Reference != "" {
		refTo, err := spec.ReferenceMessage(msg.Reference)
		if err != nil {
			return err
		}
		msg.ReferenceTo = refTo
	}

	// Process Payload
	if msg.Payload != nil {
		if err := msg.Payload.Process(msg.Name+"Payload", spec, false); err != nil {
			return err
		}
	}

	// Process Headers
	if err := msg.processHeaders(spec); err != nil {
		return err
	}

	// Process OneOf
	if err := msg.processOneOf(spec); err != nil {
		return err
	}

	// Process correlation ID
	return msg.processCorrelationID(spec)
}

func (msg *Message) processHeaders(spec Specification) error {
	// check headers exists
	if msg.Headers == nil {
		return nil
	}

	// Process headers
	if err := msg.Headers.Process(msg.Name+"Headers", spec, false); err != nil {
		return err
	}

	// Check Headers is an object, after processing
	if msg.Headers.Follow().Type != SchemaTypeIsObject.String() {
		return fmt.Errorf(
			"%w: %q headers must be an object, is %q",
			extensions.ErrAsyncAPI, msg.Name, msg.Headers.Follow().Type)
	}

	return nil
}

func (msg *Message) processCorrelationID(spec Specification) error {
	msg.createCorrelationIDFieldIfMissing()
	msg.CorrelationIDRequired = msg.isCorrelationIDRequired()
	loc, err := msg.getCorrelationIDLocation(spec)
	if err != nil {
		return err
	}
	msg.CorrelationIDLocation = loc

	return nil
}

func (msg *Message) processOneOf(spec Specification) error {
	for k, v := range msg.OneOf {
		// Process the OneOf
		if err := v.Process(msg.Name+MessageSuffix+strconv.Itoa(k), spec); err != nil {
			return err
		}

		// Merge the OneOf as one payload
		if err := msg.MergeWith(spec, *v); err != nil {
			return err
		}
	}

	return nil
}

func (msg Message) getCorrelationIDLocation(spec Specification) (string, error) {
	// Let's check the message before the reference
	if msg.CorrelationID != nil {
		return msg.CorrelationID.Location, nil
	}

	// If there is a reference, check it
	if msg.Reference != "" {
		msg, err := spec.ReferenceMessage(msg.Reference)
		if err != nil {
			return "", err
		}

		correlationID := msg.CorrelationID
		if correlationID != nil {
			return correlationID.Location, nil
		}
	}

	return "", nil
}

func (msg Message) isCorrelationIDRequired() bool {
	if msg.CorrelationID == nil || msg.CorrelationID.Location == "" {
		return false
	}

	correlationIDParent := msg.createTreeUntilCorrelationID()
	path := strings.Split(msg.CorrelationID.Location, "/")
	return correlationIDParent.IsFieldRequired(path[len(path)-1])
}

func (msg *Message) createCorrelationIDFieldIfMissing() {
	_ = msg.createTreeUntilCorrelationID()
}

func (msg *Message) createTreeUntilCorrelationID() (correlationIDParent *Schema) {
	// Check that correlationID exists
	if msg.CorrelationID == nil || msg.CorrelationID.Location == "" {
		return utils.ToPointer(NewSchema())
	}

	// Check that the correlation ID is in header
	if strings.HasPrefix(msg.CorrelationID.Location, "$message.header#") {
		return msg.createTreeUntilCorrelationIDFromMessageType(MessageFieldIsHeader)
	}

	// Check that the correlation ID is in payload
	if strings.HasPrefix(msg.CorrelationID.Location, "$message.payload#") {
		return msg.createTreeUntilCorrelationIDFromMessageType(MessageFieldIsPayload)
	}

	// Default to nothing
	return utils.ToPointer(NewSchema())
}

func (msg *Message) createTreeUntilCorrelationIDFromMessageType(t MessageField) (correlationIDParent *Schema) {
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
	return msg.downToCorrelationID(child)
}

func (msg Message) downToCorrelationID(child *Schema) (correlationIDParent *Schema) {
	var exists bool

	path := strings.Split(msg.CorrelationID.Location, "/")
	for i, v := range path[1:] {
		// Keep the parent
		correlationIDParent = child

		// Get the corresponding child
		child, exists = correlationIDParent.Properties[v]
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
			if correlationIDParent.Properties == nil {
				correlationIDParent.Properties = make(map[string]*Schema)
			}
			correlationIDParent.Properties[v] = child
		}
	}

	return correlationIDParent
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

// Follow will follow the reference to the end.
func (msg *Message) Follow() *Message {
	if msg.ReferenceTo != nil {
		return msg.ReferenceTo.Follow()
	}
	return msg
}
