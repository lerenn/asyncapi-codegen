package asyncapiv3

import (
	"strings"
)

// Specification is the asyncapi specification struct that will be used to generate
// code. It should contains every information given in the asyncapi specification.
// Source: https://www.asyncapi.com/docs/reference/specification/v2.6.0#schema
type Specification struct {
	Version    string                `json:"asyncapi"`
	Info       Info                  `json:"info"`
	Channels   map[string]*Channel   `json:"channels"`
	Components Components            `json:"components"`
	Operations map[string]*Operation `json:"operations"`
}

// Process processes the Specification to make it ready for code generation.
func (s *Specification) Process() {
	for path, ch := range s.Channels {
		ch.Process(path, *s)
	}

	for path, op := range s.Operations {
		op.Process(path, *s)
	}

	s.Components.Process(*s)
}

// GetByActionCount gets the count of 'sending' operations and the count
// of 'reception' operations inside the Specification.
func (s Specification) GetByActionCount() (sendCount, receiveCount uint) {
	for _, op := range s.Operations {
		// Check that the publish channel is present
		if op.Action.IsSend() {
			sendCount++
		}

		// Check that the subscribe channel is present
		if op.Action.IsReceive() {
			receiveCount++
		}
	}

	return sendCount, receiveCount
}

// ReferenceChannel returns the Channel struct corresponding to the given reference.
func (s Specification) ReferenceChannel(ref string) *Channel {
	msg, _ := s.reference(ref).(*Channel)
	return msg
}

// ReferenceMessage returns the Message struct corresponding to the given reference.
func (s Specification) ReferenceMessage(ref string) *Message {
	msg, _ := s.reference(ref).(*Message)
	return msg
}

// ReferenceExternalDocumentation returns the ExternalDocumentation struct corresponding to the given reference.
func (s Specification) ReferenceExternalDocumentation(ref string) *ExternalDocumentation {
	msg, _ := s.reference(ref).(*ExternalDocumentation)
	return msg
}

// ReferenceOperationBinding returns the OperationBinding struct corresponding to the given reference.
func (s Specification) ReferenceOperationBinding(ref string) *OperationBinding {
	param, _ := s.reference(ref).(*OperationBinding)
	return param
}

// ReferenceOperationReply returns the OperationReply struct corresponding to the given reference.
func (s Specification) ReferenceOperationReply(ref string) *OperationReply {
	param, _ := s.reference(ref).(*OperationReply)
	return param
}

// ReferenceOperationTrait returns the OperationTrait struct corresponding to the given reference.
func (s Specification) ReferenceOperationTrait(ref string) *OperationTrait {
	param, _ := s.reference(ref).(*OperationTrait)
	return param
}

// ReferenceParameter returns the Parameter struct corresponding to the given reference.
func (s Specification) ReferenceParameter(ref string) *Parameter {
	param, _ := s.reference(ref).(*Parameter)
	return param
}

// ReferenceSecurity returns the Security struct corresponding to the given reference.
func (s Specification) ReferenceSecurity(ref string) *SecurityScheme {
	msg, _ := s.reference(ref).(*SecurityScheme)
	return msg
}

// ReferenceSchema returns the Any struct corresponding to the given reference.
func (s Specification) ReferenceSchema(ref string) *Schema {
	msg, _ := s.reference(ref).(*Schema)
	return msg
}

func (s Specification) reference(ref string) any {
	refPath := strings.Split(ref, "/")[1:]

	switch refPath[0] {
	case "components":
		switch refPath[1] {
		case "messages":
			msg := s.Components.Messages[refPath[2]]
			return msg.referenceFrom(refPath[3:])
		case "schemas":
			schema := s.Components.Schemas[refPath[2]]
			return schema.referenceFrom(refPath[3:])
		case "parameters":
			return s.Components.Parameters[refPath[2]]
		}
	case "channels":
		return s.Channels[refPath[1]]
	}

	return nil
}

// AsyncAPIVersion returns the asyncapi version of this document.
// This function is used mainly by the interface.
func (s Specification) AsyncAPIVersion() string {
	return s.Version
}
