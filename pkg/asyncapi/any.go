package asyncapi

import (
	"github.com/lerenn/asyncapi-codegen/pkg/utils"
)

type Any struct {
	AllOf       []*Any          `json:"allOf"`
	AnyOf       []*Any          `json:"anyOf"`
	OneOf       []*Any          `json:"oneOf"`
	Type        string          `json:"type"`
	Description string          `json:"description"`
	Format      string          `json:"format"`
	Properties  map[string]*Any `json:"properties"`
	Items       *Any            `json:"items"`
	Reference   string          `json:"$ref"`
	Required    []string        `json:"required"`

	// Non AsyncAPI fields
	Name        string `json:"-"`
	ReferenceTo *Any   `json:"-"`
	IsRequired  bool   `json:"-"`

	// Embedded extended fields
	Extensions
}

func NewAny() Any {
	return Any{
		Properties: make(map[string]*Any),
		Required:   make([]string, 0),
	}
}

func (a *Any) Process(name string, spec Specification, isRequired bool) {
	a.Name = utils.UpperFirstLetter(name)

	// Add pointer to reference if there is one
	if a.Reference != "" {
		a.ReferenceTo = spec.ReferenceAny(a.Reference)
	}

	// Process Properties
	for n, p := range a.Properties {
		p.Process(n, spec, utils.IsInSlice(a.Required, n))
	}

	// Process Items
	if a.Items != nil {
		a.Items.Process(name+"Items", spec, false)
	}

	// Process AnyOf
	for _, v := range a.AnyOf {
		v.Process(name+"AnyOf", spec, false)

		// Merge with other fields as one struct (invalidate references)
		a.MergeWith(spec, *v)
	}

	// Process OneOf
	for _, v := range a.OneOf {
		v.Process(name+"OneOf", spec, false)

		// Merge with other fields as one struct (invalidate references)
		a.MergeWith(spec, *v)
	}

	// Process AllOf
	for _, v := range a.AllOf {
		v.Process(name+"AllOf", spec, false)

		// Merge with other fields as one struct (invalidate references)
		a.MergeWith(spec, *v)
	}

	// Set IsRequired
	a.IsRequired = isRequired
}

func (a Any) IsFieldRequired(field string) bool {
	return utils.IsInSlice(a.Required, field)
}

func (a *Any) referenceFrom(ref []string) *Any {
	if len(ref) == 0 {
		return a
	}

	return a.Properties[ref[0]].referenceFrom(ref[1:])
}

func (a *Any) MergeWith(spec Specification, a2 Any) {
	a.Type = "object"

	// Getting merged with reference
	if a2.Reference != "" {
		refAny2 := spec.ReferenceAny(a2.Reference)
		a2.MergeWith(spec, *refAny2)
	}

	// Merge AnyOf
	if a2.AnyOf != nil {
		if a.AnyOf == nil {
			copy(a2.AnyOf, a.AnyOf)
		} else {
			a.AnyOf = append(a.AnyOf, a2.AnyOf...)
		}
	}

	// Merge OneOf
	if a2.OneOf != nil {
		if a.OneOf == nil {
			copy(a2.OneOf, a.OneOf)
		} else {
			a.OneOf = append(a.OneOf, a2.OneOf...)
		}
	}

	// Merge properties
	if a2.Properties != nil {
		if a.Properties == nil {
			a.Properties = make(map[string]*Any)
		}

		for k, v := range a2.Properties {
			_, exists := a.Properties[k]
			if !exists {
				a.Properties[k] = v
			}
		}
	}

	// Merge requirements
	a.Required = append(a.Required, a2.Required...)
	a.Required = utils.RemoveDuplicateFromSlice(a.Required)
}
