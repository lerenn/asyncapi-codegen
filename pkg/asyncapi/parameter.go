package asyncapi

type Parameter struct {
	Description string `json:"description"`
	Schema      *Any   `json:"schema"`
	Location    string `json:"location"`

	// Non AsyncAPI fields
	Name string `json:"-"`
}

func (p *Parameter) Process(name string, spec Specification) {
	p.Name = name
}
