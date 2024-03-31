package asyncapiv2

import "github.com/lerenn/asyncapi-codegen/pkg/utils/template"

// Parameter is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v2.6.0#parameterObject
type Parameter struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Description string  `json:"description"`
	Schema      *Schema `json:"schema-name"`
	Location    string  `json:"location"`
	Reference   string  `json:"$ref"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string     `json:"-"`
	ReferenceTo *Parameter `json:"-"`
}

// Process processes the Parameter structure to make it ready for code generation.
func (p *Parameter) Process(name string, spec Specification) error {
	// Add parameter name
	p.Name = template.Namify(name)

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
