package asyncapiv3

// Parameter is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#parameterObject
type Parameter struct {
	Enum        []string `json:"enum"`
	Default     string   `json:"default"`
	Description string   `json:"description"`
	Examples    []string `json:"examples"`
	Location    string   `json:"location"`
	Reference   string   `json:"$ref"`

	// Non AsyncAPI fields
	Name        string     `json:"-"`
	ReferenceTo *Parameter `json:"-"`
}

// Process processes the Parameter structure to make it ready for code generation.
func (p *Parameter) Process(name string, spec Specification) {
	// Prevent modification if nil
	if p == nil {
		return
	}

	// Set name
	p.Name = name

	// Add pointer to reference if there is one
	if p.Reference != "" {
		p.ReferenceTo = spec.ReferenceParameter(p.Reference)
	}
}
