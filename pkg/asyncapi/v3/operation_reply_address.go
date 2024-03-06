package asyncapiv3

import "github.com/lerenn/asyncapi-codegen/pkg/utils"

// OperationReplyAddress is a representation of the corresponding asyncapi object
// filled from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#operationReplyAddressObject
type OperationReplyAddress struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Description string `json:"description"`
	Location    string `json:"location"`
	Reference   string `json:"$ref"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string                 `json:"-"`
	ReferenceTo *OperationReplyAddress `json:"-"`
}

// Process processes the OperationReplyAddress to make it ready for code generation.
func (ora *OperationReplyAddress) Process(path string, spec Specification) {
	// Prevent modification if nil
	if ora == nil {
		return
	}

	// Set name
	ora.Name = utils.UpperFirstLetter(path)

	// Add pointer to reference if there is one
	if ora.Reference != "" {
		ora.ReferenceTo = spec.ReferenceOperationReplyAddress(ora.Reference)
	}
}
