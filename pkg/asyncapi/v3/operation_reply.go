package asyncapiv3

import "github.com/lerenn/asyncapi-codegen/pkg/utils"

// OperationReply is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#operationReplyObject
type OperationReply struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Reference string `json:"$ref"`
	// TODO

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string          `json:"-"`
	ReferenceTo *OperationReply `json:"-"`
}

// Process processes the OperationReply to make it ready for code generation.
func (msg *OperationReply) Process(name string, spec Specification) {
	msg.Name = utils.UpperFirstLetter(name)

	// Add pointer to reference if there is one
	if msg.Reference != "" {
		msg.ReferenceTo = spec.ReferenceOperationReply(msg.Reference)
	}
}
