package asyncapiv3

import (
	"fmt"
	"strings"
)

const (
	// MajorVersion is the major version of this AsyncAPI implementation.
	MajorVersion = 3
)

// Specification is the asyncapi specification struct that will be used to generate
// code. It should contains every information given in the asyncapi specification.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#schema
type Specification struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Version            string                `json:"asyncapi"`
	ID                 string                `json:"id"`
	Info               Info                  `json:"info"`
	Servers            []*Server             `json:"servers"`
	DefaultContentType string                `json:"defaultContentType"`
	Channels           map[string]*Channel   `json:"channels"`
	Operations         map[string]*Operation `json:"operations"`
	Components         Components            `json:"components"`

	// --- Non AsyncAPI fields -------------------------------------------------
}

// Process processes the Specification to make it ready for code generation.
func (s *Specification) Process() {
	// Prevent modification if nil
	if s == nil {
		return
	}

	// Process info
	s.Info.Process(*s)

	// Process servers
	for i, srv := range s.Servers {
		srv.Process(fmt.Sprintf("Server%d", i), *s)
	}

	// Process channels
	for path, ch := range s.Channels {
		ch.Process(path, *s)
	}

	// Process operations
	for path, op := range s.Operations {
		op.Process(path, *s)
	}

	// Process components
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
	ch, _ := s.reference(ref).(*Channel)
	return ch
}

// ReferenceChannelBindings returns the ChannelBindings struct corresponding to the given reference.
func (s Specification) ReferenceChannelBindings(ref string) *ChannelBindings {
	bindings, _ := s.reference(ref).(*ChannelBindings)
	return bindings
}

// ReferenceExternalDocumentation returns the ExternalDocumentation struct corresponding to the given reference.
func (s Specification) ReferenceExternalDocumentation(ref string) *ExternalDocumentation {
	extDoc, _ := s.reference(ref).(*ExternalDocumentation)
	return extDoc
}

// ReferenceMessage returns the Message struct corresponding to the given reference.
func (s Specification) ReferenceMessage(ref string) *Message {
	msg, _ := s.reference(ref).(*Message)
	return msg
}

// ReferenceMessageBindings returns the MessageBindings struct corresponding to the given reference.
func (s Specification) ReferenceMessageBindings(ref string) *MessageBindings {
	bindings, _ := s.reference(ref).(*MessageBindings)
	return bindings
}

// ReferenceMessageExample returns the MessageExample struct corresponding to the given reference.
func (s Specification) ReferenceMessageExample(ref string) *MessageExample {
	bindings, _ := s.reference(ref).(*MessageExample)
	return bindings
}

// ReferenceMessageTrait returns the MessageTrait struct corresponding to the given reference.
func (s Specification) ReferenceMessageTrait(ref string) *MessageTrait {
	bindings, _ := s.reference(ref).(*MessageTrait)
	return bindings
}

// ReferenceOperation returns the Operation struct corresponding to the given reference.
func (s Specification) ReferenceOperation(ref string) *Operation {
	op, _ := s.reference(ref).(*Operation)
	return op
}

// ReferenceOperationBindings returns the OperationBindings struct corresponding to the given reference.
func (s Specification) ReferenceOperationBindings(ref string) *OperationBindings {
	bindings, _ := s.reference(ref).(*OperationBindings)
	return bindings
}

// ReferenceOperationReply returns the OperationReply struct corresponding to the given reference.
func (s Specification) ReferenceOperationReply(ref string) *OperationReply {
	opReply, _ := s.reference(ref).(*OperationReply)
	return opReply
}

// ReferenceOperationReplyAddress returns the OperationReplyAddress struct corresponding to the given reference.
func (s Specification) ReferenceOperationReplyAddress(ref string) *OperationReplyAddress {
	opReply, _ := s.reference(ref).(*OperationReplyAddress)
	return opReply
}

// ReferenceOperationTrait returns the OperationTrait struct corresponding to the given reference.
func (s Specification) ReferenceOperationTrait(ref string) *OperationTrait {
	opTrait, _ := s.reference(ref).(*OperationTrait)
	return opTrait
}

// ReferenceParameter returns the Parameter struct corresponding to the given reference.
func (s Specification) ReferenceParameter(ref string) *Parameter {
	param, _ := s.reference(ref).(*Parameter)
	return param
}

// ReferenceSecurity returns the SecurityScheme struct corresponding to the given reference.
func (s Specification) ReferenceSecurity(ref string) *SecurityScheme {
	security, _ := s.reference(ref).(*SecurityScheme)
	return security
}

// ReferenceSchema returns the Schema struct corresponding to the given reference.
func (s Specification) ReferenceSchema(ref string) *Schema {
	schema, _ := s.reference(ref).(*Schema)
	return schema
}

// ReferenceServer returns the Server struct corresponding to the given reference.
func (s Specification) ReferenceServer(ref string) *Server {
	bindings, _ := s.reference(ref).(*Server)
	return bindings
}

// ReferenceServerBindings returns the ServerBindings struct corresponding to the given reference.
func (s Specification) ReferenceServerBindings(ref string) *ServerBindings {
	bindings, _ := s.reference(ref).(*ServerBindings)
	return bindings
}

// ReferenceServerVariable returns the ServerVariable struct corresponding to the given reference.
func (s Specification) ReferenceServerVariable(ref string) *ServerVariable {
	bindings, _ := s.reference(ref).(*ServerVariable)
	return bindings
}

// ReferenceTag returns the Tag struct corresponding to the given reference.
func (s Specification) ReferenceTag(ref string) *Tag {
	bindings, _ := s.reference(ref).(*Tag)
	return bindings
}

//nolint:funlen,cyclop // Not necessary to reduce statements and cyclop
func (s Specification) reference(ref string) any {
	refPath := strings.Split(ref, "/")[1:]

	switch refPath[0] {
	case "components":
		switch refPath[1] {
		case "schemas":
			schema := s.Components.Schemas[refPath[2]]
			return schema.referenceFrom(refPath[3:])
		case "servers":
			return s.Components.Servers[refPath[2]]
		case "channels":
			return s.Components.Channels[refPath[2]]
		case "operations":
			return s.Components.Operations[refPath[2]]
		case "messages":
			msg := s.Components.Messages[refPath[2]]
			return msg.referenceFrom(refPath[3:])
		case "securitySchemes":
			return s.Components.SecuritySchemes[refPath[2]]
		case "serverVariables":
			return s.Components.ServerVariables[refPath[2]]
		case "parameters":
			return s.Components.Parameters[refPath[2]]
		case "correlationIds":
			return s.Components.CorrelationIDs[refPath[2]]
		case "replies":
			return s.Components.Replies[refPath[2]]
		case "replyAddresses":
			return s.Components.ReplyAddresses[refPath[2]]
		case "externalDocs":
			return s.Components.ExternalDocs[refPath[2]]
		case "tags":
			return s.Components.Tags[refPath[2]]
		case "operationTraits":
			return s.Components.OperationTraits[refPath[2]]
		case "messageTraits":
			return s.Components.MessageTraits[refPath[2]]
		case "serverBindings":
			return s.Components.ServerBindings[refPath[2]]
		case "channelBindings":
			return s.Components.ChannelBindings[refPath[2]]
		case "operationBindings":
			return s.Components.OperationBindings[refPath[2]]
		case "messageBindings":
			return s.Components.MessageBindings[refPath[2]]
		}
	case "channels":
		return s.Channels[refPath[1]]
	}

	return nil
}

// MajorVersion returns the asyncapi major version of this document.
// This function is used mainly by the interface.
func (s Specification) MajorVersion() int {
	return MajorVersion
}
