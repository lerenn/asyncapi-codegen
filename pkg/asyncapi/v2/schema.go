package asyncapiv2

import (
	"github.com/lerenn/asyncapi-codegen/pkg/utils"
	"github.com/lerenn/asyncapi-codegen/pkg/utils/template"
)

// SchemaType is a structure that represents the type of a field.
type SchemaType string

// String returns the string representation of the type.
func (st SchemaType) String() string {
	return string(st)
}

const (
	// SchemaTypeIsArray represents the type of an array.
	SchemaTypeIsArray SchemaType = "array"
	// SchemaTypeIsObject represents the type of an object.
	SchemaTypeIsObject SchemaType = "object"
	// SchemaTypeIsString represents the type of a string.
	SchemaTypeIsString SchemaType = "string"
	// SchemaTypeIsInteger represents the type of an integer.
	SchemaTypeIsInteger SchemaType = "integer"
)

// Schema is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v2.6.0#schemaObject
type Schema struct {
	// --- JSON Schema fields --------------------------------------------------

	AllOf                []*Schema          `json:"allOf"`
	AnyOf                []*Schema          `json:"anyOf"`
	OneOf                []*Schema          `json:"oneOf"`
	Type                 string             `json:"type"`
	Description          string             `json:"description"`
	Format               string             `json:"format"`
	Properties           map[string]*Schema `json:"properties"`
	Items                *Schema            `json:"items"`
	Reference            string             `json:"$ref"`
	Required             []string           `json:"required"`
	AdditionalProperties *Schema            `json:"additionalProperties"`

	// --- Non JSON Schema/AsyncAPI fields -------------------------------------

	Name        string  `json:"-"`
	ReferenceTo *Schema `json:"-"`
	IsRequired  bool    `json:"-"`
	// Embedded extended fields
	Extensions
}

// NewSchema creates a new Schema structure with initialized fields.
func NewSchema() Schema {
	return Schema{
		Properties: make(map[string]*Schema),
		Required:   make([]string, 0),
	}
}

// Process processes the Schema structure to make it ready for code generation.
func (s *Schema) Process(name string, spec Specification, isRequired bool) error {
	s.Name = template.Namify(name)

	// Add pointer to reference if there is one
	if err := s.processReference(spec); err != nil {
		return err
	}

	// Process Properties
	if err := s.processProperties(spec); err != nil {
		return err
	}

	// Process Items
	if err := s.processItems(spec); err != nil {
		return err
	}

	// Process AnyOf
	if err := s.processAnyOf(spec); err != nil {
		return err
	}

	// Process OneOf
	if err := s.processOneOf(spec); err != nil {
		return err
	}

	// Process AllOf
	if err := s.processAllOf(spec); err != nil {
		return err
	}

	// Process AdditionalProperties
	if err := s.processAdditionalProperties(spec); err != nil {
		return err
	}

	// Set IsRequired
	s.IsRequired = isRequired

	return nil
}

func (s *Schema) processReference(spec Specification) error {
	if s.Reference != "" {
		refTo, err := spec.ReferenceSchema(s.Reference)
		if err != nil {
			return err
		}
		s.ReferenceTo = refTo
	}

	return nil
}

func (s *Schema) processProperties(spec Specification) error {
	for n, p := range s.Properties {
		if err := p.Process(
			s.Name+template.Namify(n),
			spec,
			utils.IsInSlice(s.Required, n),
		); err != nil {
			return err
		}
	}

	return nil
}

func (s *Schema) processItems(spec Specification) error {
	if s.Items != nil {
		if err := s.Items.Process(s.Name+"Item", spec, false); err != nil {
			return err
		}
	}

	return nil
}

func (s *Schema) processAnyOf(spec Specification) error {
	for _, v := range s.AnyOf {
		if err := v.Process(s.Name+"AnyOf", spec, false); err != nil {
			return err
		}

		// Merge with other fields as one struct (invalidate references)
		if err := s.MergeWith(spec, *v); err != nil {
			return err
		}
	}

	return nil
}

func (s *Schema) processOneOf(spec Specification) error {
	for _, v := range s.OneOf {
		// Process the OneOf
		if err := v.Process(s.Name+"OneOf", spec, false); err != nil {
			return err
		}

		// Merge the OneOf as one payload
		if err := s.MergeWith(spec, *v); err != nil {
			return err
		}
	}

	return nil
}

func (s *Schema) processAllOf(spec Specification) error {
	for _, v := range s.AllOf {
		if err := v.Process(s.Name+"AllOf", spec, false); err != nil {
			return err
		}

		// Merge with other fields as one struct (invalidate references)
		if err := s.MergeWith(spec, *v); err != nil {
			return err
		}
	}

	return nil
}

func (s *Schema) processAdditionalProperties(spec Specification) error {
	if s.AdditionalProperties != nil {
		if err := s.AdditionalProperties.Process(s.Name+"AdditionalProperties", spec, false); err != nil {
			return err
		}
	}

	return nil
}

// IsFieldRequired checks if a field is required in the asyncapi struct.
func (s Schema) IsFieldRequired(field string) bool {
	return utils.IsInSlice(s.Required, field)
}

func (s *Schema) referenceFrom(ref []string) *Schema {
	if len(ref) == 0 {
		return s
	}

	return s.Properties[ref[0]].referenceFrom(ref[1:])
}

// MergeWith merges the given Schema structure with another one
// (basically for AllOf, AnyOf, OneOf, etc).
func (s *Schema) MergeWith(spec Specification, s2 Schema) error {
	s.Type = SchemaTypeIsObject.String()

	// Getting merged with reference
	if err := s.mergeWithMessageFromReference(s2, spec); err != nil {
		return err
	}

	// Merge AllOf
	if err := s.mergeWithMessageAllOf(s2); err != nil {
		return err
	}

	// Merge AnyOf
	if err := s.mergeWithMessageAnyOf(s2); err != nil {
		return err
	}

	// Merge OneOf
	if err := s.mergeWithMessageOneOf(s2); err != nil {
		return err
	}

	// Merge properties
	if err := s.mergeWithMessageProperties(s2); err != nil {
		return err
	}

	// Merge requirements
	s.Required = append(s.Required, s2.Required...)
	s.Required = utils.RemoveDuplicateFromSlice(s.Required)

	return nil
}

// Follow follows the reference to the end and returns the final Schema.
func (s *Schema) Follow() *Schema {
	if s.ReferenceTo != nil {
		return s.ReferenceTo.Follow()
	}

	return s
}

func (s *Schema) mergeWithMessageFromReference(s2 Schema, spec Specification) error {
	if s2.Reference != "" {
		refAny2, err := spec.ReferenceSchema(s2.Reference)
		if err != nil {
			return err
		}

		if err := s2.MergeWith(spec, *refAny2); err != nil {
			return err
		}
	}

	return nil
}

func (s *Schema) mergeWithMessageAllOf(s2 Schema) error {
	if s2.AllOf != nil {
		if s.AllOf == nil {
			copy(s2.AllOf, s.AllOf)
		} else {
			s.AllOf = append(s.AllOf, s2.AllOf...)
		}
	}

	return nil
}

func (s *Schema) mergeWithMessageAnyOf(s2 Schema) error {
	if s2.AnyOf != nil {
		if s.AnyOf == nil {
			copy(s2.AnyOf, s.AnyOf)
		} else {
			s.AnyOf = append(s.AnyOf, s2.AnyOf...)
		}
	}

	return nil
}

func (s *Schema) mergeWithMessageOneOf(s2 Schema) error {
	if s2.OneOf != nil {
		if s.OneOf == nil {
			copy(s2.OneOf, s.OneOf)
		} else {
			s.OneOf = append(s.OneOf, s2.OneOf...)
		}
	}

	return nil
}

func (s *Schema) mergeWithMessageProperties(s2 Schema) error {
	if s2.Properties != nil {
		if s.Properties == nil {
			s.Properties = make(map[string]*Schema)
		}

		for k, v := range s2.Properties {
			_, exists := s.Properties[k]
			if !exists {
				s.Properties[k] = v
			}
		}
	}

	return nil
}
