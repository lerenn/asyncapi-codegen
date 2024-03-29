package asyncapiv3

import "github.com/lerenn/asyncapi-codegen/pkg/utils/template"

// MessageExample is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#messageExampleObject
type MessageExample struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Headers   map[string]any `json:"headers"`
	Payload   map[string]any `json:"payload"`
	Name      string         `json:"name"`
	Summary   string         `json:"summary"`
	Reference string         `json:"$ref"`

	// --- Non AsyncAPI fields -------------------------------------------------

	ReferenceTo *MessageExample `json:"-"`
}

// Process processes the MessageExample to make it ready for code generation.
func (me *MessageExample) Process(path string, spec Specification) error {
	// Prevent modification if nil
	if me == nil {
		return nil
	}

	// Set name
	me.Name = template.Namify(path)

	// Add pointer to reference if there is one
	if me.Reference != "" {
		refTo, err := spec.ReferenceMessageExample(me.Reference)
		if err != nil {
			return err
		}
		me.ReferenceTo = refTo
	}

	return nil
}
