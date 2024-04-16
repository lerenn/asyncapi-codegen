package asyncapiv3

import "github.com/lerenn/asyncapi-codegen/pkg/utils/template"

// ServerVariable is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#serverVariableObject
type ServerVariable struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Enum        []string `json:"enum"`
	Default     string   `json:"default"`
	Description string   `json:"description"`
	Examples    []string `json:"examples"`
	Reference   string   `json:"$ref"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string          `json:"-"`
	ReferenceTo *ServerVariable `json:"-"`
}

// generateMetadata generates metadata for the ServerVariable.
func (sv *ServerVariable) generateMetadata(path string) {
	// Prevent modification if nil
	if sv == nil {
		return
	}

	// Set name
	sv.Name = template.Namify(path)
}

// setDependencies sets dependencies between the different elements of the ServerVariable.
func (sv *ServerVariable) setDependencies(spec Specification) error {
	// Prevent modification if nil
	if sv == nil {
		return nil
	}

	// Add pointer to reference if there is one
	if sv.Reference != "" {
		refTo, err := spec.ReferenceServerVariable(sv.Reference)
		if err != nil {
			return err
		}
		sv.ReferenceTo = refTo
	}

	return nil
}
