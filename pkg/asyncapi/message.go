package asyncapi

import (
	"strings"

	"github.com/lerenn/asyncapi-codegen/pkg/utils"
)

type Message struct {
	Description   string         `json:"description"`
	Headers       *Any           `json:"headers"`
	Payload       *Any           `json:"payload"`
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

func (msg *Message) Process(name string, spec Specification) {
	msg.Name = utils.UpperFirstLetter(name)

	// Add pointer to reference if there is one
	if msg.Reference != "" {
		msg.ReferenceTo = spec.ReferenceMessage(msg.Reference)
	}

	// Process Headers and Payload
	if msg.Headers != nil {
		msg.Headers.Process(name+"Headers", spec)
	}
	if msg.Payload != nil {
		msg.Payload.Process(name+"Payload", spec)
	}

	// Process correlation ID
	msg.createCorrelationIDFieldIfMissing()
	msg.CorrelationIDLocation = msg.getCorrelationIDLocation(spec)
	msg.CorrelationIDRequired = msg.isCorrelationIDRequired()
}

func (msg Message) getCorrelationIDLocation(spec Specification) string {
	// Let's check the message before the reference
	if msg.CorrelationID != nil {
		return msg.CorrelationID.Location
	}

	// If there is a reference, check it
	if msg.Reference != "" {
		correlationID := spec.ReferenceMessage(msg.Reference).CorrelationID
		if correlationID != nil {
			return correlationID.Location
		}
	}

	return ""
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

func (msg *Message) createTreeUntilCorrelationID() (correlationIDParent *Any) {
	// Check that correlationID exists
	if msg.CorrelationID == nil || msg.CorrelationID.Location == "" {
		return utils.ToNullable(NewAny())
	}

	// Get root from header or payload
	var child *Any
	path := strings.Split(msg.CorrelationID.Location, "/")
	if strings.HasPrefix(msg.CorrelationID.Location, "$message.header#") {
		if msg.Headers == nil {
			msg.Headers = utils.ToNullable(NewAny())
			msg.Headers.Name = "headers"
			msg.Headers.Type = "object"
		}
		child = msg.Headers
	} else if strings.HasPrefix(msg.CorrelationID.Location, "$message.payload#") && msg.Payload != nil {
		if msg.Payload == nil {
			msg.Payload = utils.ToNullable(NewAny())
			msg.Payload.Name = "headers"
			msg.Payload.Type = "object"
		}
		child = msg.Payload
	}

	// Go down the path to correlation ID
	var exists bool
	for i, v := range path[1:] {
		correlationIDParent = child
		child, exists = correlationIDParent.Properties[v]
		if !exists {
			// Create child
			child = utils.ToNullable(NewAny())
			child.Name = v
			if i == len(path)-2 { // As there is -1 in the loop slice
				child.Type = "string"
			} else {
				child.Type = "object"
			}

			// Add it to parent
			if correlationIDParent.Properties == nil {
				correlationIDParent.Properties = make(map[string]*Any)
			}
			correlationIDParent.Properties[v] = child
		}
	}

	return correlationIDParent
}

func (msg *Message) referenceFrom(ref []string) interface{} {
	if len(ref) == 0 {
		return msg
	}

	var next *Any
	if ref[0] == "payload" {
		next = msg.Payload
	} else if ref[0] == "header" {
		next = msg.Headers
	}

	return next.referenceFrom(ref[1:])
}
