package asyncapi

import (
	"github.com/lerenn/asyncapi-codegen/pkg/utils"
)

type Any struct {
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
	a.Name = name

	// Add pointer to reference if there is one
	if a.Reference != "" {
		a.ReferenceTo = spec.ReferenceAny(a.Reference)
	}

	for _, p := range a.Properties {
		p.Process("", spec)
	}

	if a.Items != nil {
		a.Items.Process("", spec)
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
