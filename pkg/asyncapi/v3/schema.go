package asyncapiv3

import (
	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
	"github.com/lerenn/asyncapi-codegen/pkg/utils"
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

	Title                string             `json:"title,omitempty"`
	Type                 string             `json:"type,omitempty"`
	Examples             []any              `json:"examples,omitempty"`
	ReadOnly             bool               `json:"readOnly,omitempty"`
	WriteOnly            bool               `json:"writeOnly,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty"`
	PatternProperties    map[string]*Schema `json:"patternProperties,omitempty"`
	AdditionalProperties *Schema            `json:"additionalProperties,omitempty"`
	AdditionalItems      []*Schema          `json:"additionalItems,omitempty"`
	Items                *Schema            `json:"items,omitempty"`
	PropertyNames        []string           `json:"propertyNames,omitempty"`
	Contains             []*Schema          `json:"contains,omitempty"`
	AllOf                []*Schema          `json:"allOf,omitempty"`
	AnyOf                []*Schema          `json:"anyOf,omitempty"`
	OneOf                []*Schema          `json:"oneOf,omitempty"`
	Not                  *Schema            `json:"not,omitempty"`

	// --- AsyncAPI specific ---------------------------------------------------

	Description string `json:"description,omitempty"`
	Format      string `json:"format,omitempty"`
	Default     any    `json:"default,omitempty"`

	Reference string `json:"$ref,omitempty"`

	// --- Non Json Schema/AsyncAPI fields -------------------------------------

	Name        string  `json:"-"`
	ReferenceTo *Schema `json:"-"`

	// Embedded validation fields
	asyncapi.Validations[Schema]

	// --- Embedded extended fields --------------------------------------------

	Extensions
}

// NewSchema creates a new Schema structure with initialized fields.
func NewSchema() Schema {
	return Schema{
		Properties: make(map[string]*Schema),
		Validations: asyncapi.Validations[Schema]{
			Required: make([]string, 0),
		},
	}
}

// generateMetadata generates metadata for the Schema and its children.
//
//nolint:funlen,cyclop // Not necessary to reduce length and cyclop
func (s *Schema) generateMetadata(parentName, name string, number *int, isRequired bool) error {
	// Prevent modification if nil
	if s == nil {
		return nil
	}

	// Set name
	// NOTE: do not specify the type "schema" in the name
	s.Name = generateFullName(parentName, name, "", number)

	// Generate Properties metadata
	for n, p := range s.Properties {
		if err := p.generateMetadata(s.Name, n+"_Property", nil, utils.IsInSlice(s.Required, n)); err != nil {
			return err
		}
	}

	// Generate Pattern Properties metadata
	for n, p := range s.PatternProperties {
		if err := p.generateMetadata(s.Name, n+"_Pattern_Property", nil, utils.IsInSlice(s.Required, n)); err != nil {
			return err
		}
	}

	// Generate AdditionalProperties metadata
	if s.AdditionalProperties != nil {
		if err := s.AdditionalProperties.generateMetadata(s.Name, "Additional_Properties", nil, false); err != nil {
			return err
		}
	}

	// Generate Additional Items metadata
	for i, item := range s.AdditionalItems {
		if err := item.generateMetadata(s.Name, "Additional_Item", &i, false); err != nil {
			return err
		}
	}

	// Generate Items metadata
	// NOTE: give the name of the parent to the items
	if err := s.Items.generateMetadata("", "Item_From_"+s.Name, nil, false); err != nil {
		return err
	}

	// Generate Contains metadata
	for i, item := range s.Contains {
		if err := item.generateMetadata("", s.Name+"_Content", &i, false); err != nil {
			return err
		}
	}

	// Generate AnyOf metadata
	for _, v := range s.AnyOf {
		if err := v.generateMetadata(s.Name, "Any_Of", nil, false); err != nil {
			return err
		}
	}

	// Generate OneOf metadata
	for _, v := range s.OneOf {
		if err := v.generateMetadata(s.Name, "One_Of", nil, false); err != nil {
			return err
		}
	}

	// Generate AllOf metadata
	for _, v := range s.AllOf {
		if err := v.generateMetadata(s.Name, "All_Of", nil, false); err != nil {
			return err
		}
	}

	// Generate Not metadata
	if err := s.Not.generateMetadata(s.Name, "Not_Schema", nil, false); err != nil {
		return err
	}

	// Set IsRequired
	s.IsRequired = isRequired

	s.ShouldOmitEmpty = s.ExtOmitEmpty

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
	for _, p := range s.Properties {
		if err := p.setDependencies(spec); err != nil {
			return err
		}
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
	if err := s.setAnyOfDependenciesAndMerge(spec); err != nil {
		return err
	}

	// Set OneOf dependencies
	if err := s.setOneOfDependenciesAndMerge(spec); err != nil {
		return err
	}

	// Set AllOf dependencies
	if err := s.setAllOfDependenciesAndMerge(spec); err != nil {
		return err
	}

	// Set Not dependencies
	if err := s.Not.setDependencies(spec); err != nil {
		return err
	}

	return nil
}

func (s *Schema) setAllOfDependenciesAndMerge(spec Specification) error {
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

func (s *Schema) setOneOfDependenciesAndMerge(spec Specification) error {
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

func (s *Schema) setAnyOfDependenciesAndMerge(spec Specification) error {
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
	s.Validations.Merge(refTo.Validations)

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
				Validations: asyncapi.Validations[Schema]{
					IsRequired: v.IsRequired,
				},
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
