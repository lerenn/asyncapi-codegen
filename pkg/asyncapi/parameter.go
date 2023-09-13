package asyncapi

// Parameter is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v2.6.0#parameterObject
type Parameter struct {
	Description string `json:"description"`
	Schema      *Any   `json:"schema"`
	Location    string `json:"location"`
	Reference   string `json:"$ref"`

	// Non AsyncAPI fields
	Name        string     `json:"-"`
	ReferenceTo *Parameter `json:"-"`
}

// Process processes the Parameter structure to make it ready for code generation
func (p *Parameter) Process(name string, spec Specification) {
	// Add parameter name
	p.Name = name

	// Add pointer to reference if there is one
	if p.Reference != "" {
		p.ReferenceTo = spec.ReferenceParameter(p.Reference)
	}
}
