package asyncapiv3

import (
	"fmt"

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
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#schemaObject
type Schema struct {
	// --- Original JSON Schema Definition -------------------------------------

	Title                string             `json:"title"`
	Type                 string             `json:"type"`
	Required             []string           `json:"required"`
	MultipleOf           []string           `json:"multipleOf"`
	Maximum              string             `json:"maximum"`
	ExclusiveMaximum     float64            `json:"exclusiveMaximum"`
	Minimum              string             `json:"minimum"`
	ExclusiveMinimum     float64            `json:"exclusiveMinimum"`
	MaxLength            uint               `json:"maxLength"`
	MinLength            uint               `json:"minLength"`
	Pattern              string             `json:"pattern"`
	MaxItems             uint               `json:"maxItems"`
	MinItems             uint               `json:"minItems"`
	UniqueItems          bool               `json:"uniqueItems"`
	MaxProperties        uint               `json:"maxProperties"`
	MinProperties        uint               `json:"minProperties"`
	Enum                 []any              `json:"enum"`
	Const                any                `json:"const"`
	Examples             []any              `json:"examples"`
	ReadOnly             bool               `json:"readOnly"`
	WriteOnly            bool               `json:"writeOnly"`
	Properties           map[string]*Schema `json:"properties"`
	PatternProperties    map[string]*Schema `json:"patternProperties"`
	AdditionalProperties *Schema            `json:"additionalProperties"`
	AdditionalItems      []*Schema          `json:"additionalItems"`
	Items                *Schema            `json:"items"`
	PropertyNames        []string           `json:"propertyNames"`
	Contains             []*Schema          `json:"contains"`
	AllOf                []*Schema          `json:"allOf"`
	AnyOf                []*Schema          `json:"anyOf"`
	OneOf                []*Schema          `json:"oneOf"`
	Not                  *Schema            `json:"not"`

	// --- AsyncAPI specific ---------------------------------------------------

	Description string `json:"description"`
	Format      string `json:"format"`
	Default     any    `json:"default"`

	Reference string `json:"$ref"`

	// --- Non Json Schema/AsyncAPI fields -------------------------------------

	Name        string  `json:"-"`
	ReferenceTo *Schema `json:"-"`
	IsRequired  bool    `json:"-"`

	// --- Embedded extended fields --------------------------------------------

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
//
//nolint:funlen,cyclop // Not necessary to reduce length and cyclop
func (s *Schema) Process(name string, spec Specification, isRequired bool) error {
	// Prevent modification if nil
	if s == nil {
		return nil
	}

	// Set name
	s.Name = template.Namify(name)

	// Process reference
	if err := s.processReference(spec); err != nil {
		return err
	}

	// Process Properties
	if err := s.processProperties(spec); err != nil {
		return err
	}

	// Process Pattern Properties
	for n, p := range s.PatternProperties {
		if err := p.Process(n+"PatternProperty", spec, utils.IsInSlice(s.Required, n)); err != nil {
			return err
		}
	}

	// Process AdditionalProperties
	if s.AdditionalProperties != nil {
		if err := s.AdditionalProperties.Process(s.Name+"AdditionalProperties", spec, false); err != nil {
			return err
		}
	}

	// Process Additional Items
	for i, item := range s.AdditionalItems {
		if err := item.Process(fmt.Sprintf("%sAdditionalItem%d", s.Name, i), spec, false); err != nil {
			return err
		}
	}

	// Process Items
	if err := s.Items.Process(s.Name+"Items", spec, false); err != nil {
		return err
	}

	// Process Contains
	for i, item := range s.Contains {
		if err := item.Process(fmt.Sprintf("%sContains%d", name, i), spec, false); err != nil {
			return err
		}
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

	// Process Not
	if err := s.Not.Process(s.Name+"Not", spec, false); err != nil {
		return err
	}

	// Set IsRequired
	s.IsRequired = isRequired

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

func (s *Schema) processOneOf(spec Specification) error {
	for _, v := range s.OneOf {
		if err := v.Process(s.Name+"OneOf", spec, false); err != nil {
			return err
		}

		// Merge with other fields as one struct (invalidate references)
		if err := s.MergeWith(spec, *v); err != nil {
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

func (s *Schema) processProperties(spec Specification) error {
	for n, p := range s.Properties {
		if err := p.Process(n+"Property", spec, utils.IsInSlice(s.Required, n)); err != nil {
			return err
		}
	}
	return nil
}

func (s *Schema) processReference(spec Specification) error {
	if s.Reference == "" {
		return nil
	}

	// Add pointer to reference
	refTo, err := spec.ReferenceSchema(s.Reference)
	if err != nil {
		return err
	}
	s.ReferenceTo = refTo

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
//
//nolint:cyclop
func (s *Schema) MergeWith(spec Specification, s2 Schema) error {
	if s == nil {
		return nil
	}

	s.Type = SchemaTypeIsObject.String()

	// Getting merged with reference
	if s2.Reference != "" {
		refAny2, err := spec.ReferenceSchema(s2.Reference)
		if err != nil {
			return err
		}

		if err := s2.MergeWith(spec, *refAny2); err != nil {
			return err
		}
	}

	// Merge AnyOf
	if s2.AnyOf != nil {
		if s.AnyOf == nil {
			copy(s2.AnyOf, s.AnyOf)
		} else {
			s.AnyOf = append(s.AnyOf, s2.AnyOf...)
		}
	}

	// Merge OneOf
	if s2.OneOf != nil {
		if s.OneOf == nil {
			copy(s2.OneOf, s.OneOf)
		} else {
			s.OneOf = append(s.OneOf, s2.OneOf...)
		}
	}

	// Merge properties
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

	// Merge requirements
	s.Required = append(s.Required, s2.Required...)
	s.Required = utils.RemoveDuplicateFromSlice(s.Required)

	return nil
}

// Follow returns referenced schema if specified or the actual schema.
func (s *Schema) Follow() *Schema {
	if s.ReferenceTo != nil {
		return s.ReferenceTo
	}
	return s
}
