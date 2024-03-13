package asyncapiv3

import (
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/utils"
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
	AdditionalProperties map[string]*Schema `json:"additionalProperties"`
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
func (s *Schema) Process(name string, spec Specification, isRequired bool) {
	// Prevent modification if nil
	if s == nil {
		return
	}

	// Set name
	s.Name = utils.UpperFirstLetter(name)

	// Add pointer to reference if there is one
	if s.Reference != "" {
		s.ReferenceTo = spec.ReferenceSchema(s.Reference)
	}

	// Process Properties
	for n, p := range s.Properties {
		p.Process(n+"Property", spec, utils.IsInSlice(s.Required, n))
	}

	// Process Pattern Properties
	for n, p := range s.PatternProperties {
		p.Process(n+"PatternProperty", spec, utils.IsInSlice(s.Required, n))
	}

	// Process Additional Properties
	for n, p := range s.AdditionalProperties {
		p.Process(n+"AdditionalProperties", spec, utils.IsInSlice(s.Required, n))
	}

	// Process Additional Items
	for i, item := range s.AdditionalItems {
		item.Process(fmt.Sprintf("%sAdditionalItem%d", name, i), spec, false)
	}

	// Process Items
	s.Items.Process(name+"Items", spec, false)

	// Process Contains
	for i, item := range s.Contains {
		item.Process(fmt.Sprintf("%sContains%d", name, i), spec, false)
	}

	// Process AnyOf
	for _, v := range s.AnyOf {
		v.Process(name+"AnyOf", spec, false)

		// Merge with other fields as one struct (invalidate references)
		s.MergeWith(spec, *v)
	}

	// Process OneOf
	for _, v := range s.OneOf {
		v.Process(name+"OneOf", spec, false)

		// Merge with other fields as one struct (invalidate references)
		s.MergeWith(spec, *v)
	}

	// Process AllOf
	for _, v := range s.AllOf {
		v.Process(name+"AllOf", spec, false)

		// Merge with other fields as one struct (invalidate references)
		s.MergeWith(spec, *v)
	}

	// Process Not
	s.Not.Process(name+"Not", spec, false)

	// Set IsRequired
	s.IsRequired = isRequired
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
func (s *Schema) MergeWith(spec Specification, s2 Schema) {
	if s == nil {
		return
	}

	s.Type = MessageTypeIsObject.String()

	// Getting merged with reference
	if s2.Reference != "" {
		refAny2 := spec.ReferenceSchema(s2.Reference)
		s2.MergeWith(spec, *refAny2)
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
}
