package asyncapiv2

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
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

	// specificationReferenced is a map of all the outside specifications that
	// are referenced in this specification.
	dependencies map[string]*Specification
}

// NewSpecification creates a new Specification struct.
func NewSpecification() *Specification {
	return &Specification{
		Channels:     make(map[string]*Channel),
		dependencies: make(map[string]*Specification),
	}
}

// AddDependency adds a specification dependency to the Specification.
func (s *Specification) AddDependency(path string, spec asyncapi.Specification) error {
	// Cast to Specification v2
	specV2, ok := spec.(*Specification)
	if !ok {
		return fmt.Errorf(
			"%w: cannot cast %q into 'Specification' (type is %q)",
			extensions.ErrAsyncAPI, path, reflect.TypeOf(spec))
	}

	// Remove local prefix './' if present
	path = strings.TrimPrefix(path, "./")

	// Set in dependencies
	s.dependencies[path] = specV2

	return nil
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
	for _, spec := range s.dependencies {
		if err := spec.generateMetadata(); err != nil {
			return err
		}
	}

	for path, ch := range s.Channels {
		if err := ch.generateMetadata(path); err != nil {
			return err
		}
	}

	return s.Components.generateMetadata()
}

// setDependencies set dependencies between the different elements of the Specification.
func (s *Specification) setDependencies() error {
	for _, spec := range s.dependencies {
		if err := spec.setDependencies(); err != nil {
			return err
		}
	}

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

func (s Specification) getDependencyBasedOnRef(ref string) (*Specification, string, error) {
	// Separate file from path
	fileAndPath := strings.Split(ref, "#")
	if len(fileAndPath) != 2 {
		return nil, "", fmt.Errorf("%w: invalid reference %q", ErrInvalidReference, ref)
	}
	file, ref := fileAndPath[0], fileAndPath[1]

	// If the file if not empty, it should be the file where the reference is
	if file == "" {
		return &s, ref, nil
	}

	// Remove local prefix './' if present
	file = strings.TrimPrefix(file, "./")

	// Get corresponding dependency
	s2, ok := s.dependencies[file]
	if !ok {
		return nil, "", fmt.Errorf(
			"%w: file %q is not referenced in dependencies %+v",
			ErrInvalidReference, file, s.dependencies)
	}

	return s2, ref, nil
}

func (s Specification) reference(ref string) (any, error) {
	// Separate file from path
	usedSpec, ref, err := s.getDependencyBasedOnRef(ref)
	if err != nil {
		return nil, err
	}

	// Separate each part of the reference
	ref = strings.TrimPrefix(ref, "/")
	refPath := strings.Split(ref, "/")

	switch refPath[0] {
	case "components":
		switch refPath[1] {
		case "messages":
			msg := usedSpec.Components.Messages[refPath[2]]
			return msg.referenceFrom(refPath[3:]), nil
		case "schemas":
			schema := usedSpec.Components.Schemas[refPath[2]]
			return schema.referenceFrom(refPath[3:]), nil
		case "parameters":
			return usedSpec.Components.Parameters[refPath[2]], nil
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

// FromUnknownVersion returns an AsyncAPI specification V2 from interface, if compatible.
// Note: Before using this, you should make sure that parsed data is in version 2.
func FromUnknownVersion(s asyncapi.Specification) (*Specification, error) {
	spec, ok := s.(*Specification)
	if !ok {
		return nil, fmt.Errorf("unknown spec format: should have been a v2 format")
	}

	return spec, nil
}
