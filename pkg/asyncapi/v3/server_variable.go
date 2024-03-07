package asyncapiv3

import "github.com/lerenn/asyncapi-codegen/pkg/utils"

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

// Process processes the ServerVariable to make it ready for code generation.
func (sv *ServerVariable) Process(path string, spec Specification) {
	// Prevent modification if nil
	if sv == nil {
		return
	}

	// Set name
	sv.Name = utils.UpperFirstLetter(path)

	// Add pointer to reference if there is one
	if sv.Reference != "" {
		sv.ReferenceTo = spec.ReferenceServerVariable(sv.Reference)
	}
}
