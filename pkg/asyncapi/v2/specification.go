package asyncapiv2

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
)

const (
	// MajorVersion is the major version of this AsyncAPI implementation.
	MajorVersion = 2
)

var (
	// ErrInvalidReference is sent when a reference is invalid.
	ErrInvalidReference = fmt.Errorf("%w: invalid reference", extensions.ErrAsyncAPI)
)

// Specification is the asyncapi specification struct that will be used to generate
// code. It should contains every information given in the asyncapi specification.
// Source: https://www.asyncapi.com/docs/reference/specification/v2.6.0#schema
type Specification struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Version    string              `json:"asyncapi"`
	Info       Info                `json:"info"`
	Channels   map[string]*Channel `json:"channels"`
	Components Components          `json:"components"`

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
	for path, ch := range s.Channels {
		if err := ch.generateMetadata(path); err != nil {
			return err
		}
	}

	return s.Components.generateMetadata()
}

// setDependencies set dependencies between the different elements of the Specification.
func (s *Specification) setDependencies() error {
	for _, ch := range s.Channels {
		if err := ch.setDependencies(*s); err != nil {
			return err
		}
	}

	return s.Components.setDependencies(*s)
}

// GetPublishSubscribeCount gets the count of 'publish' channels and the count
// of 'subscribe' channels inside the Specification.
func (s Specification) GetPublishSubscribeCount() (publishCount, subscribeCount uint) {
	for _, c := range s.Channels {
		// Check that the publish channel is present
		if c.Publish != nil {
			publishCount++
		}

		// Check that the subscribe channel is present
		if c.Subscribe != nil {
			subscribeCount++
		}
	}

	return publishCount, subscribeCount
}

// ReferenceParameter returns the Parameter struct corresponding to the given reference.
func (s Specification) ReferenceParameter(ref string) (*Parameter, error) {
	// Get object pointed by reference
	obj, err := s.reference(ref)
	if err != nil {
		return nil, err
	}

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

// ReferenceMessage returns the Message struct corresponding to the given reference.
func (s Specification) ReferenceMessage(ref string) (*Message, error) {
	// Get object pointed by reference
	obj, err := s.reference(ref)
	if err != nil {
		return nil, err
	}

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

// ReferenceSchema returns the Any struct corresponding to the given reference.
func (s Specification) ReferenceSchema(ref string) (*Schema, error) {
	// Get object pointed by reference
	obj, err := s.reference(ref)
	if err != nil {
		return nil, err
	}

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

func (s Specification) reference(ref string) (any, error) {
	refPath := strings.Split(ref, "/")[1:]

	switch refPath[0] {
	case "components":
		switch refPath[1] {
		case "messages":
			msg := s.Components.Messages[refPath[2]]
			return msg.referenceFrom(refPath[3:]), nil
		case "schemas":
			schema := s.Components.Schemas[refPath[2]]
			return schema.referenceFrom(refPath[3:]), nil
		case "parameters":
			return s.Components.Parameters[refPath[2]], nil
		default:
			return nil, fmt.Errorf("%w: %q from reference %q is not supported", ErrInvalidReference, refPath[1], ref)
		}
	default:
		return nil, fmt.Errorf("%w: %q from reference %q is not supported", ErrInvalidReference, refPath[0], ref)
	}
}

// MajorVersion returns the asyncapi major version of this document.
// This function is used mainly by the interface.
func (s Specification) MajorVersion() int {
	return MajorVersion
}
