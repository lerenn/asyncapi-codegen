package asyncapi

type Parameter struct {
	Description string `json:"description"`
	Schema      *Any   `json:"schema"`
	Location    string `json:"location"`
	Reference   string `json:"$ref"`

	// Non AsyncAPI fields
	Name        string     `json:"-"`
	ReferenceTo *Parameter `json:"-"`
}

func (p *Parameter) Process(name string, spec Specification) {
	// Add parameter name
	p.Name = name

	// Add pointer to reference if there is one
	if p.Reference != "" {
		p.ReferenceTo = spec.ReferenceParameter(p.Reference)
	}
}
