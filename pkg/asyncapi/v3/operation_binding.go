package asyncapiv3

import "github.com/lerenn/asyncapi-codegen/pkg/utils"

// OperationBinding is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#operationBindingsObject
type OperationBinding struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Reference string `json:"$ref"`
	// TODO

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string            `json:"-"`
	ReferenceTo *OperationBinding `json:"-"`
}

// Process processes the OperationBinding to make it ready for code generation.
func (ob *OperationBinding) Process(name string, spec Specification) {
	ob.Name = utils.UpperFirstLetter(name)

	// Add pointer to reference if there is one
	if ob.Reference != "" {
		ob.ReferenceTo = spec.ReferenceOperationBinding(ob.Reference)
	}
}
