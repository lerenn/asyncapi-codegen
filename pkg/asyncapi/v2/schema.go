package asyncapiv2

import (
	"github.com/TheSadlig/asyncapi-codegen/pkg/asyncapi"
	"github.com/TheSadlig/asyncapi-codegen/pkg/utils"
	"github.com/TheSadlig/asyncapi-codegen/pkg/utils/template"
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
	Type                 string             `json:"type"`
	Description          string             `json:"description"`
	Format               string             `json:"format"`
	Properties           map[string]*Schema `json:"properties"`
	Items                *Schema            `json:"items"`
	Reference            string             `json:"$ref"`
	AdditionalProperties *Schema            `json:"additionalProperties"`

	// --- Non JSON Schema/AsyncAPI fields -------------------------------------

	Name        string  `json:"-"`
	ReferenceTo *Schema `json:"-"`

	// Embedded validation fields
	asyncapi.Validations[Schema]

	// Embedded extended fields
	Extensions
}

// NewSchema creates a new Schema structure with initialized fields.
func NewSchema() Schema {
	return Schema{
		Validations: asyncapi.Validations[Schema]{
			Required: make([]string, 0),
		},
		Properties: make(map[string]*Schema),
	}
}

// generateMetadata generates metadata for the schema and its children.
func (s *Schema) generateMetadata(name string, isRequired bool) error {
	s.Name = template.Namify(name)

	// Generate Properties metadata
	if err := s.generatePropertiesMetadata(); err != nil {
		return err
	}

	// Generate Items metadata
	if err := s.generateItemsMetadata(); err != nil {
		return err
	}

	// Generate AnyOf metadata
	if err := s.generateaAnyOfMetadata(); err != nil {
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

	// Generate AdditionalProperties metadata
	if err := s.generateAdditionalPropertiesMetadata(); err != nil {
		return err
	}

	// Set IsRequired
	s.IsRequired = isRequired

	return nil
}

// setDependencies sets dependencies for the schema from the specification.
func (s *Schema) setDependencies(spec Specification) error {
	// Reference to another schema if specified
	if s.Reference != "" {
		refTo, err := spec.ReferenceSchema(s.Reference)
		if err != nil {
			return err
		}
		s.ReferenceTo = refTo
	}

	// Set properties dependencies
	if err := s.setPropertiesDependencies(spec); err != nil {
		return err
	}

	// Set items dependencies
	if err := s.setItemsDependencies(spec); err != nil {
		return err
	}

	// Set AnyOf links
	if err := s.setAnyOfDependencies(spec); err != nil {
		return err
	}

	// Set OneOf links
	if err := s.setOneOfDependencies(spec); err != nil {
		return err
	}

	// Set AllOf links
	if err := s.setAllOfDependencies(spec); err != nil {
		return err
	}

	// Set AdditionalProperties links
	if err := s.setAdditionalPropertiesDependencies(spec); err != nil {
		return err
	}

	return nil
}

func (s *Schema) generatePropertiesMetadata() error {
	for n, p := range s.Properties {
		if err := p.generateMetadata(
			s.Name+template.Namify(n),
			utils.IsInSlice(s.Required, n),
		); err != nil {
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

func (s *Schema) generateItemsMetadata() error {
	if s.Items != nil {
		if err := s.Items.generateMetadata(s.Name+"Item", false); err != nil {
			return err
		}
	}

	return nil
}

func (s *Schema) setItemsDependencies(spec Specification) error {
	if s.Items != nil {
		if err := s.Items.setDependencies(spec); err != nil {
			return err
		}
	}

	return nil
}

func (s *Schema) generateaAnyOfMetadata() error {
	for _, v := range s.AnyOf {
		if err := v.generateMetadata(s.Name+"AnyOf", false); err != nil {
			return err
		}
	}

	return nil
}

func (s *Schema) setAnyOfDependencies(spec Specification) error {
	for _, v := range s.AnyOf {
		// Set dependencies
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
		// Set dependencies
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
		// Set dependencies
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

func (s *Schema) generateAdditionalPropertiesMetadata() error {
	if s.AdditionalProperties != nil {
		if err := s.AdditionalProperties.generateMetadata(s.Name+"AdditionalProperties", false); err != nil {
			return err
		}
	}

	return nil
}

func (s *Schema) setAdditionalPropertiesDependencies(spec Specification) error {
	if s.AdditionalProperties != nil {
		if err := s.AdditionalProperties.setDependencies(spec); err != nil {
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
