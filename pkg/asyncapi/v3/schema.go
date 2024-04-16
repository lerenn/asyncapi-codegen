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

// generateMetadata generates metadata for the Schema and its children.
//
//nolint:funlen,cyclop // Not necessary to reduce length and cyclop
func (s *Schema) generateMetadata(name string, isRequired bool) error {
	// Prevent modification if nil
	if s == nil {
		return nil
	}

	// Set name
	s.Name = template.Namify(name)

	// Generate Properties metadata
	if err := s.generatePropertiesMetadata(); err != nil {
		return err
	}

	// Generate Pattern Properties metadata
	for n, p := range s.PatternProperties {
		if err := p.generateMetadata(n+"PatternProperty", utils.IsInSlice(s.Required, n)); err != nil {
			return err
		}
	}

	// Generate AdditionalProperties metadata
	if s.AdditionalProperties != nil {
		if err := s.AdditionalProperties.generateMetadata(s.Name+"AdditionalProperties", false); err != nil {
			return err
		}
	}

	// Generate Additional Items metadata
	for i, item := range s.AdditionalItems {
		if err := item.generateMetadata(fmt.Sprintf("%sAdditionalItem%d", s.Name, i), false); err != nil {
			return err
		}
	}

	// Generate Items metadata
	if err := s.Items.generateMetadata(s.Name+"Item", false); err != nil {
		return err
	}

	// Generate Contains metadata
	for i, item := range s.Contains {
		if err := item.generateMetadata(fmt.Sprintf("%sContains%d", name, i), false); err != nil {
			return err
		}
	}

	// Generate AnyOf metadata
	if err := s.generateAnyOfMetadata(); err != nil {
		return err
	}

	// Generate OneOf metadata
	if err := s.generateOneOfMetadata(); err != nil {
		return err
	}

	// Generate AllOf metadata
	if err := s.generateAllOfMetadata(); err != nil {
		return err
	}

	// Generate Not metadata
	if err := s.Not.generateMetadata(s.Name+"Not", false); err != nil {
		return err
	}

	// Set IsRequired
	s.IsRequired = isRequired

	return nil
}

//
//nolint:funlen,cyclop // Not necessary to reduce length and cyclop
func (s *Schema) setDependencies(spec Specification) error {
	// Prevent modification if nil
	if s == nil {
		return nil
	}

	// Set reference
	if err := s.setReference(spec); err != nil {
		return err
	}

	// Set Properties dependencies
	if err := s.setPropertiesDependencies(spec); err != nil {
		return err
	}

	// Set Pattern Properties dependencies
	for _, p := range s.PatternProperties {
		if err := p.setDependencies(spec); err != nil {
			return err
		}
	}

	// Set AdditionalProperties dependencies
	if s.AdditionalProperties != nil {
		if err := s.AdditionalProperties.setDependencies(spec); err != nil {
			return err
		}
	}

	// Set Additional Items dependencies
	for _, item := range s.AdditionalItems {
		if err := item.setDependencies(spec); err != nil {
			return err
		}
	}

	// Set Items dependencies
	if err := s.Items.setDependencies(spec); err != nil {
		return err
	}

	// Set Contains dependencies
	for _, item := range s.Contains {
		if err := item.setDependencies(spec); err != nil {
			return err
		}
	}

	// Set AnyOf dependencies
	if err := s.setAnyOfDependencies(spec); err != nil {
		return err
	}

	// Set OneOf dependencies
	if err := s.setOneOfDependencies(spec); err != nil {
		return err
	}

	// Set AllOf dependencies
	if err := s.setAllOfDependencies(spec); err != nil {
		return err
	}

	// Set Not dependencies
	if err := s.Not.setDependencies(spec); err != nil {
		return err
	}

	return nil
}

func (s *Schema) generateAllOfMetadata() error {
	for _, v := range s.AllOf {
		if err := v.generateMetadata(s.Name+"AllOf", false); err != nil {
			return err
		}
	}
	return nil
}

func (s *Schema) setAllOfDependencies(spec Specification) error {
	for _, v := range s.AllOf {
		if err := v.setDependencies(spec); err != nil {
			return err
		}

		// Merge with other fields as one struct (invalidate references)
		if err := s.MergeWith(spec, *v); err != nil {
			return err
		}
	}

	return nil
}

func (s *Schema) generateOneOfMetadata() error {
	for _, v := range s.OneOf {
		if err := v.generateMetadata(s.Name+"OneOf", false); err != nil {
			return err
		}
	}

	return nil
}

func (s *Schema) setOneOfDependencies(spec Specification) error {
	for _, v := range s.OneOf {
		if err := v.setDependencies(spec); err != nil {
			return err
		}

		// Merge with other fields as one struct (invalidate references)
		if err := s.MergeWith(spec, *v); err != nil {
			return err
		}
	}

	return nil
}

func (s *Schema) generateAnyOfMetadata() error {
	for _, v := range s.AnyOf {
		if err := v.generateMetadata(s.Name+"AnyOf", false); err != nil {
			return err
		}
	}

	return nil
}

func (s *Schema) setAnyOfDependencies(spec Specification) error {
	for _, v := range s.AnyOf {
		if err := v.setDependencies(spec); err != nil {
			return err
		}

		// Merge with other fields as one struct (invalidate references)
		if err := s.MergeWith(spec, *v); err != nil {
			return err
		}
	}

	return nil
}

func (s *Schema) generatePropertiesMetadata() error {
	for n, p := range s.Properties {
		if err := p.generateMetadata(n+"Property", utils.IsInSlice(s.Required, n)); err != nil {
			return err
		}
	}
	return nil
}

func (s *Schema) setPropertiesDependencies(spec Specification) error {
	for _, p := range s.Properties {
		if err := p.setDependencies(spec); err != nil {
			return err
		}
	}
	return nil
}

func (s *Schema) setReference(spec Specification) error {
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
func (s *Schema) MergeWith(spec Specification, s2 Schema) error {
	if s == nil {
		return nil
	}

	s.Type = SchemaTypeIsObject.String()

	// Merge with other fields
	s.mergeWithSchemaAllOf(s2)
	s.mergeWithSchemaAnyOf(s2)
	s.mergeWithSchemaOneOf(s2)
	s.mergeWithSchemaProperties(s2)

	// Merge requirements
	s.Required = append(s.Required, s2.Required...)
	s.Required = utils.RemoveDuplicateFromSlice(s.Required)

	return nil
}

func (s *Schema) mergeWithSchemaAllOf(s2 Schema) {
	// Return if there are no AllOf to merge
	if s2.AllOf == nil && (s2.ReferenceTo == nil || s2.ReferenceTo.AllOf == nil) {
		return
	}

	// Initialize AllOf if they are nil
	if s.AllOf == nil {
		s.AllOf = make([]*Schema, 0)
	}

	// Add AllOf from s2 to s
	if s2.AllOf != nil {
		s.AllOf = append(s.AllOf, s2.AllOf...)
	}

	// Add AllOf from s2 reference to s
	if s2.ReferenceTo != nil && s2.ReferenceTo.AllOf != nil {
		for _, v := range s2.ReferenceTo.AllOf {
			s.AllOf = append(s.AllOf, &Schema{ReferenceTo: v})
		}
	}
}

func (s *Schema) mergeWithSchemaAnyOf(s2 Schema) {
	// Return if there are no AnyOf to merge
	if s2.AnyOf == nil && (s2.ReferenceTo == nil || s2.ReferenceTo.AnyOf == nil) {
		return
	}

	// Initialize AnyOf if they are nil
	if s.AnyOf == nil {
		s.AnyOf = make([]*Schema, 0)
	}

	// Add AnyOf from s2 to s
	if s2.AnyOf != nil {
		s.AnyOf = append(s.AnyOf, s2.AnyOf...)
	}

	// Add AnyOf from s2 reference to s
	if s2.ReferenceTo != nil && s2.ReferenceTo.AnyOf != nil {
		for _, v := range s2.ReferenceTo.AnyOf {
			s.AnyOf = append(s.AnyOf, &Schema{ReferenceTo: v})
		}
	}
}

func (s *Schema) mergeWithSchemaOneOf(s2 Schema) {
	// Return if there are no OneOf to merge
	if s2.OneOf == nil && (s2.ReferenceTo == nil || s2.ReferenceTo.OneOf == nil) {
		return
	}

	// Initialize OneOf if they are nil
	if s.OneOf == nil {
		s.OneOf = make([]*Schema, 0)
	}

	// Add OneOf from s2 to s
	if s2.OneOf != nil {
		s.OneOf = append(s.OneOf, s2.OneOf...)
	}

	// Add OneOf from s2 reference to s
	if s2.ReferenceTo != nil && s2.ReferenceTo.OneOf != nil {
		for _, v := range s2.ReferenceTo.OneOf {
			s.OneOf = append(s.OneOf, &Schema{ReferenceTo: v})
		}
	}
}

func (s *Schema) mergeWithSchemaProperties(s2 Schema) {
	// Return if there are no properties to merge
	if s2.Properties == nil && (s2.ReferenceTo == nil || s2.ReferenceTo.Properties == nil) {
		return
	}

	// Initialize properties if they are nil
	if s.Properties == nil {
		s.Properties = make(map[string]*Schema)
	}

	// Add properties from s2 to s
	for k, v := range s2.Properties {
		_, exists := s.Properties[k]
		if exists {
			continue
		}

		s.Properties[k] = v
	}

	// Add properties from s2 reference to s
	s.mergeWithSchemaReferenceProperties(s2)
}

func (s *Schema) mergeWithSchemaReferenceProperties(s2 Schema) {
	// Return if there are no properties to merge
	if s2.ReferenceTo == nil || s2.ReferenceTo.Properties == nil {
		return
	}

	// Add properties from s2 reference to s
	for k, v := range s2.ReferenceTo.Properties {
		// Skip if the property already exists
		_, exists := s.Properties[k]
		if exists {
			continue
		}

		// Add the property
		if v.Type == "object" {
			s.Properties[k] = &Schema{
				IsRequired:  v.IsRequired,
				ReferenceTo: v,
			}
		} else {
			s.Properties[k] = v
		}

		// Add to required if it is required
		if v.IsRequired {
			s.Required = append(s.Required, k)
		}
	}
}

// Follow returns referenced schema if specified or the actual schema.
func (s *Schema) Follow() *Schema {
	if s.ReferenceTo != nil {
		return s.ReferenceTo
	}
	return s
}
