package asyncapiv3

// Parameter is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#parameterObject
type Parameter struct {
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default,omitempty"`
	Description string   `json:"description,omitempty"`
	Examples    []string `json:"examples,omitempty"`
	Location    string   `json:"location,omitempty"`
	Reference   string   `json:"$ref,omitempty"`

	// Non AsyncAPI fields
	Name        string     `json:"-"`
	ReferenceTo *Parameter `json:"-"`
}

// generateMetadata generates metadata for the Parameter.
func (p *Parameter) generateMetadata(parentName, name string) {
	// Prevent modification if nil
	if p == nil {
		return
	}

	// Set name
	p.Name = generateFullName(parentName, name, "Parameter", nil)
}

// setDependencies sets dependencies between the different elements of the Parameter.
func (p *Parameter) setDependencies(spec Specification) error {
	// Prevent modification if nil
	if p == nil {
		return nil
	}

	// Add pointer to reference if there is one
	if p.Reference != "" {
		refTo, err := spec.ReferenceParameter(p.Reference)
		if err != nil {
			return err
		}
		p.ReferenceTo = refTo
	}

	return nil
}
