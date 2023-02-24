package asyncapi

import (
	"github.com/lerenn/asyncapi-codegen/pkg/utils"
)

type Any struct {
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
}

func NewAny() Any {
	return Any{
		Properties: make(map[string]*Any),
		Required:   make([]string, 0),
	}
}

func (a *Any) Process(name string, spec Specification) {
	a.Name = utils.UpperFirstLetter(name)

	// Add pointer to reference if there is one
	if a.Reference != "" {
		a.ReferenceTo = spec.ReferenceAny(a.Reference)
	}

	// Process Properties
	for n, p := range a.Properties {
		p.Process(n, spec)
	}

	// Process Items
	if a.Items != nil {
		a.Items.Process(name+"Items", spec)
	}

	// Process AnyOf
	for _, v := range a.AnyOf {
		v.Process(name+"AnyOf", spec)
	}

	// Process OneOf
	for _, v := range a.OneOf {
		v.Process(name+"OneOf", spec)
	}
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
