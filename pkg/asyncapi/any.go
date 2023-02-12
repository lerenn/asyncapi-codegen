package asyncapi

import "github.com/lerenn/asyncapi-codegen/pkg/utils"

type Any struct {
	Type        string          `json:"type"`
	Description string          `json:"description"`
	Format      string          `json:"format"`
	Properties  map[string]*Any `json:"properties"`
	Items       Items           `json:"items"`
	Reference   string          `json:"$ref"`
	Required    []string        `json:"required"`

	// Non AsyncAPI fields
	Name string `json:"-"`
}

func NewAny() Any {
	return Any{
		Properties: make(map[string]*Any),
		Required:   make([]string, 0),
	}
}

func (a Any) IsFieldRequired(field string) bool {
	return utils.IsInSlice(a.Required, field)
}
