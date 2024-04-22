package asyncapiv3

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

const (
	// MajorVersion is the major version of this AsyncAPI implementation.
	MajorVersion = 3
)

var (
	// ErrInvalidReference is sent when a reference is invalid.
	ErrInvalidReference = fmt.Errorf("%w: invalid reference", extensions.ErrAsyncAPI)
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
func (s *Specification) Process() error {
	if err := s.generateMetadata(); err != nil {
		return err
	}

	return s.setDependencies()
}

// generateMetadata generate metadata for the Specification and its children.
func (s *Specification) generateMetadata() error {
	// Prevent modification if nil
	if s == nil {
		return nil
	}

	// Generate metadata for components
	if err := s.Components.generateMetadata(); err != nil {
		return err
	}

	// Generate metadata for info
	if err := s.Info.generateMetadata(); err != nil {
		return err
	}

	// Generate servers metadata
	for i, srv := range s.Servers {
		srv.generateMetadata(fmt.Sprintf("Server%d", i))
	}

	// Generate metadata for channels
	for name, ch := range s.Channels {
		if err := ch.generateMetadata(name + ChannelSuffix); err != nil {
			return err
		}
	}

	// Generate metadata for operations
	for name, op := range s.Operations {
		if err := op.generateMetadata(name + "Operation"); err != nil {
			return err
		}
	}

	return nil
}

// setDependencies set dependencies between the different elements of the Specification.
func (s *Specification) setDependencies() error {
	// Prevent modification if nil
	if s == nil {
		return nil
	}

	// Set dependencies for components
	if err := s.Components.setDependencies(*s); err != nil {
		return err
	}

	// Set dependencies for info
	if err := s.Info.setDependencies(*s); err != nil {
		return err
	}

	// Set dependencies for servers
	for _, srv := range s.Servers {
		if err := srv.setDependencies(*s); err != nil {
			return err
		}
	}

	// Set dependencies for channels
	for _, ch := range s.Channels {
		if err := ch.setDependencies(*s); err != nil {
			return err
		}
	}

	// Set dependencies for operations
	for _, op := range s.Operations {
		if err := op.setDependencies(*s); err != nil {
			return err
		}
	}

	return nil
}

// GetOperationCountByAction gets the count of 'sending' operations and the count
// of 'reception' operations inside the Specification.
func (s Specification) GetOperationCountByAction() (sendCount, receiveCount uint) {
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
func (s Specification) ReferenceChannel(ref string) (*Channel, error) {
	// Get object pointed by reference
	obj := s.reference(ref)

	// Cast to channel
	channel, ok := obj.(*Channel)
	if !ok {
		return nil, fmt.Errorf(
			"%w: cannot cast %q into 'Channel' (type is %q)",
			ErrInvalidReference, ref, reflect.TypeOf(obj))
	}

	// Check that channel is not nil
	if channel == nil {
		return nil, fmt.Errorf("%w: empty target for channel reference %q", ErrInvalidReference, ref)
	}

	return channel, nil
}

// ReferenceChannelBindings returns the ChannelBindings struct corresponding to the given reference.
func (s Specification) ReferenceChannelBindings(ref string) (*ChannelBindings, error) {
	// Get object pointed by reference
	obj := s.reference(ref)

	// Cast to channel bindings
	bindings, ok := obj.(*ChannelBindings)
	if !ok {
		return nil, fmt.Errorf(
			"%w: cannot cast %q into 'ChannelBindings' (type is %q)",
			ErrInvalidReference, ref, reflect.TypeOf(obj))
	}

	// Check that channel bindings is not nil
	if bindings == nil {
		return nil, fmt.Errorf("%w: empty target for channel bindings reference %q", ErrInvalidReference, ref)
	}

	return bindings, nil
}

// ReferenceExternalDocumentation returns the ExternalDocumentation struct corresponding to the given reference.
func (s Specification) ReferenceExternalDocumentation(ref string) (*ExternalDocumentation, error) {
	// Get object pointed by reference
	obj := s.reference(ref)

	// Cast to external documentation
	doc, ok := obj.(*ExternalDocumentation)
	if !ok {
		return nil, fmt.Errorf(
			"%w: cannot cast %q into 'ExternalDocumentation' (type is %q)",
			ErrInvalidReference, ref, reflect.TypeOf(obj))
	}

	// Check that external documentation is not nil
	if doc == nil {
		return nil, fmt.Errorf("%w: empty target for external documentation reference %q", ErrInvalidReference, ref)
	}

	return doc, nil
}

// ReferenceMessage returns the Message struct corresponding to the given reference.
func (s Specification) ReferenceMessage(ref string) (*Message, error) {
	// Get object pointed by reference
	obj := s.reference(ref)

	// Cast to message
	msg, ok := obj.(*Message)
	if !ok {
		return nil, fmt.Errorf(
			"%w: cannot cast %q into 'Message' (type is %q)",
			ErrInvalidReference, ref, reflect.TypeOf(obj))
	}

	// Check that message is not nil
	if msg == nil {
		return nil, fmt.Errorf("%w: empty target for message reference %q", ErrInvalidReference, ref)
	}

	return msg, nil
}

// ReferenceMessageBindings returns the MessageBindings struct corresponding to the given reference.
func (s Specification) ReferenceMessageBindings(ref string) (*MessageBindings, error) {
	// Get object pointed by reference
	obj := s.reference(ref)

	// Cast to message bindings
	bindings, ok := obj.(*MessageBindings)
	if !ok {
		return nil, fmt.Errorf(
			"%w: cannot cast %q into 'MessageBindings' (type is %q)",
			ErrInvalidReference, ref, reflect.TypeOf(obj))
	}

	// Check that message bindings is not nil
	if bindings == nil {
		return nil, fmt.Errorf("%w: empty target for message bindings reference %q", ErrInvalidReference, ref)
	}

	return bindings, nil
}

// ReferenceMessageExample returns the MessageExample struct corresponding to the given reference.
func (s Specification) ReferenceMessageExample(ref string) (*MessageExample, error) {
	// Get object pointed by reference
	obj := s.reference(ref)

	// Cast to message example
	example, ok := obj.(*MessageExample)
	if !ok {
		return nil, fmt.Errorf(
			"%w: cannot cast %q into 'MessageExample' (type is %q)",
			ErrInvalidReference, ref, reflect.TypeOf(obj))
	}

	// Check that message example is not nil
	if example == nil {
		return nil, fmt.Errorf("%w: empty target for message example reference %q", ErrInvalidReference, ref)
	}

	return example, nil
}

// ReferenceMessageTrait returns the MessageTrait struct corresponding to the given reference.
func (s Specification) ReferenceMessageTrait(ref string) (*MessageTrait, error) {
	// Get object pointed by reference
	obj := s.reference(ref)

	// Cast to message trait
	trait, ok := obj.(*MessageTrait)
	if !ok {
		return nil, fmt.Errorf(
			"%w: cannot cast %q into 'MessageTrait' (type is %q)",
			ErrInvalidReference, ref, reflect.TypeOf(obj))
	}

	// Check that message trait is not nil
	if trait == nil {
		return nil, fmt.Errorf("%w: empty target for message trait reference %q", ErrInvalidReference, ref)
	}

	return trait, nil
}

// ReferenceOperation returns the Operation struct corresponding to the given reference.
func (s Specification) ReferenceOperation(ref string) (*Operation, error) {
	// Get object pointed by reference
	obj := s.reference(ref)

	// Cast to operation
	operation, ok := obj.(*Operation)
	if !ok {
		return nil, fmt.Errorf(
			"%w: cannot cast %q into 'Operation' (type is %q)",
			ErrInvalidReference, ref, reflect.TypeOf(obj))
	}

	// Check that operation is not nil
	if operation == nil {
		return nil, fmt.Errorf("%w: empty target for operation reference %q", ErrInvalidReference, ref)
	}

	return operation, nil
}

// ReferenceOperationBindings returns the OperationBindings struct corresponding to the given reference.
func (s Specification) ReferenceOperationBindings(ref string) (*OperationBindings, error) {
	// Get object pointed by reference
	obj := s.reference(ref)

	// Cast to operation bindings
	bindings, ok := obj.(*OperationBindings)
	if !ok {
		return nil, fmt.Errorf(
			"%w: cannot cast %q into 'OperationBindings' (type is %q)",
			ErrInvalidReference, ref, reflect.TypeOf(obj))
	}

	// Check that operation bindings is not nil
	if bindings == nil {
		return nil, fmt.Errorf("%w: empty target for operation bindings reference %q", ErrInvalidReference, ref)
	}

	return bindings, nil
}

// ReferenceOperationReply returns the OperationReply struct corresponding to the given reference.
func (s Specification) ReferenceOperationReply(ref string) (*OperationReply, error) {
	// Get object pointed by reference
	obj := s.reference(ref)

	// Cast to operation reply
	reply, ok := obj.(*OperationReply)
	if !ok {
		return nil, fmt.Errorf(
			"%w: cannot cast %q into 'OperationReply' (type is %q)",
			ErrInvalidReference, ref, reflect.TypeOf(obj))
	}

	// Check that operation reply is not nil
	if reply == nil {
		return nil, fmt.Errorf("%w: empty target for operation reply reference %q", ErrInvalidReference, ref)
	}

	return reply, nil
}

// ReferenceOperationReplyAddress returns the OperationReplyAddress struct corresponding to the given reference.
func (s Specification) ReferenceOperationReplyAddress(ref string) (*OperationReplyAddress, error) {
	// Get object pointed by reference
	obj := s.reference(ref)

	// Cast to operation reply address
	address, ok := obj.(*OperationReplyAddress)
	if !ok {
		return nil, fmt.Errorf(
			"%w: cannot cast %q into 'OperationReplyAddress' (type is %q)",
			ErrInvalidReference, ref, reflect.TypeOf(obj))
	}

	// Check that operation reply address is not nil
	if address == nil {
		return nil, fmt.Errorf("%w: empty target for operation reply address reference %q", ErrInvalidReference, ref)
	}

	return address, nil
}

// ReferenceOperationTrait returns the OperationTrait struct corresponding to the given reference.
func (s Specification) ReferenceOperationTrait(ref string) (*OperationTrait, error) {
	// Get object pointed by reference
	obj := s.reference(ref)

	// Cast to operation trait
	trait, ok := obj.(*OperationTrait)
	if !ok {
		return nil, fmt.Errorf(
			"%w: cannot cast %q into 'OperationTrait' (type is %q)",
			ErrInvalidReference, ref, reflect.TypeOf(obj))
	}

	// Check that operation trait is not nil
	if trait == nil {
		return nil, fmt.Errorf("%w: empty target for operation trait reference %q", ErrInvalidReference, ref)
	}

	return trait, nil
}

// ReferenceParameter returns the Parameter struct corresponding to the given reference.
func (s Specification) ReferenceParameter(ref string) (*Parameter, error) {
	// Get object pointed by reference
	obj := s.reference(ref)

	// Cast to parameter
	param, ok := obj.(*Parameter)
	if !ok {
		return nil, fmt.Errorf(
			"%w: cannot cast %q into 'Parameter' (type is %q)",
			ErrInvalidReference, ref, reflect.TypeOf(obj))
	}

	// Check that parameter is not nil
	if param == nil {
		return nil, fmt.Errorf("%w: empty target for parameter reference %q", ErrInvalidReference, ref)
	}

	return param, nil
}

// ReferenceSecurity returns the SecurityScheme struct corresponding to the given reference.
func (s Specification) ReferenceSecurity(ref string) (*SecurityScheme, error) {
	// Get object pointed by reference
	obj := s.reference(ref)

	// Cast to security scheme
	security, ok := obj.(*SecurityScheme)
	if !ok {
		return nil, fmt.Errorf(
			"%w: cannot cast %q into 'SecurityScheme' (type is %q)",
			ErrInvalidReference, ref, reflect.TypeOf(obj))
	}

	// Check that security scheme is not nil
	if security == nil {
		return nil, fmt.Errorf("%w: empty target for security scheme reference %q", ErrInvalidReference, ref)
	}

	return security, nil
}

// ReferenceSchema returns the Schema struct corresponding to the given reference.
func (s Specification) ReferenceSchema(ref string) (*Schema, error) {
	// Get object pointed by reference
	obj := s.reference(ref)

	// Cast to schema
	schema, ok := obj.(*Schema)
	if !ok {
		return nil, fmt.Errorf(
			"%w: cannot cast %q into 'Schema' (type is %q)",
			ErrInvalidReference, ref, reflect.TypeOf(obj))
	}

	// Check that schema is not nil
	if schema == nil {
		return nil, fmt.Errorf("%w: empty target for schema reference %q", ErrInvalidReference, ref)
	}

	return schema, nil
}

// ReferenceServer returns the Server struct corresponding to the given reference.
func (s Specification) ReferenceServer(ref string) (*Server, error) {
	// Get object pointed by reference
	obj := s.reference(ref)

	// Cast to server
	server, ok := obj.(*Server)
	if !ok {
		return nil, fmt.Errorf(
			"%w: cannot cast %q into 'Server' (type is %q)",
			ErrInvalidReference, ref, reflect.TypeOf(obj))
	}

	// Check that server is not nil
	if server == nil {
		return nil, fmt.Errorf("%w: empty target for server reference %q", ErrInvalidReference, ref)
	}

	return server, nil
}

// ReferenceServerBindings returns the ServerBindings struct corresponding to the given reference.
func (s Specification) ReferenceServerBindings(ref string) (*ServerBindings, error) {
	// Get object pointed by reference
	obj := s.reference(ref)

	// Cast to server bindings
	bindings, ok := obj.(*ServerBindings)
	if !ok {
		return nil, fmt.Errorf(
			"%w: cannot cast %q into 'ServerBindings' (type is %q)",
			ErrInvalidReference, ref, reflect.TypeOf(obj))
	}

	// Check that server bindings is not nil
	if bindings == nil {
		return nil, fmt.Errorf("%w: empty target for server bindings reference %q", ErrInvalidReference, ref)
	}

	return bindings, nil
}

// ReferenceServerVariable returns the ServerVariable struct corresponding to the given reference.
func (s Specification) ReferenceServerVariable(ref string) (*ServerVariable, error) {
	// Get object pointed by reference
	obj := s.reference(ref)

	// Cast to server variable
	variable, ok := obj.(*ServerVariable)
	if !ok {
		return nil, fmt.Errorf(
			"%w: cannot cast %q into 'ServerVariable' (type is %q)",
			ErrInvalidReference, ref, reflect.TypeOf(obj))
	}

	// Check that server variable is not nil
	if variable == nil {
		return nil, fmt.Errorf("%w: empty target for server variable reference %q", ErrInvalidReference, ref)
	}

	return variable, nil
}

// ReferenceTag returns the Tag struct corresponding to the given reference.
func (s Specification) ReferenceTag(ref string) (*Tag, error) {
	// Get object pointed by reference
	obj := s.reference(ref)

	// Cast to tag
	tag, ok := obj.(*Tag)
	if !ok {
		return nil, fmt.Errorf(
			"%w: cannot cast %q into 'Tag' (type is %q)",
			ErrInvalidReference, ref, reflect.TypeOf(obj))
	}

	// Check that tag is not nil
	if tag == nil {
		return nil, fmt.Errorf("%w: empty target for tag reference %q", ErrInvalidReference, ref)
	}

	return tag, nil
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
		default:
			return fmt.Errorf("%w: %q from reference %q is not supported", ErrInvalidReference, refPath[1], ref)
		}
	case "channels":
		return s.Channels[refPath[1]]
	default:
		return fmt.Errorf("%w: %q from reference %q is not supported", ErrInvalidReference, refPath[0], ref)
	}
}

// MajorVersion returns the asyncapi major version of this document.
// This function is used mainly by the interface.
func (s Specification) MajorVersion() int {
	return MajorVersion
}
