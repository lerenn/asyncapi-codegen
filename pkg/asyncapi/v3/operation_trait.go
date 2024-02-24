package asyncapiv3

import "github.com/lerenn/asyncapi-codegen/pkg/utils"

// OperationTrait is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#operationTraitObject
type OperationTrait struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Reference string `json:"$ref"`
	// TODO

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string          `json:"-"`
	ReferenceTo *OperationTrait `json:"-"`
}

// Process processes the OperationTrait to make it ready for code generation.
func (msg *OperationTrait) Process(name string, spec Specification) {
	msg.Name = utils.UpperFirstLetter(name)

	// Add pointer to reference if there is one
	if msg.Reference != "" {
		msg.ReferenceTo = spec.ReferenceOperationTrait(msg.Reference)
	}
}
