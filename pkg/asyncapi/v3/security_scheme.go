package asyncapiv3

import "github.com/lerenn/asyncapi-codegen/pkg/utils"

// SecurityScheme is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#securitySchemeObject
type SecurityScheme struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Reference string `json:"$ref"`
	// TODO

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string          `json:"-"`
	ReferenceTo *SecurityScheme `json:"-"`
}

// Process processes the SecurityScheme to make it ready for code generation.
func (s *SecurityScheme) Process(name string, spec Specification) {
	// Prevent modification if nil
	if s == nil {
		return
	}

	// Set name
	s.Name = utils.UpperFirstLetter(name)

	// Add pointer to reference if there is one
	if s.Reference != "" {
		s.ReferenceTo = spec.ReferenceSecurity(s.Reference)
	}
}
