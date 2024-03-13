package asyncapiv3

import (
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/utils"
)

// OperationReply is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#operationReplyObject
type OperationReply struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Address   *OperationReplyAddress `json:"address"`
	Channel   *Channel               `json:"channel"`  // Reference only
	Messages  []*Message             `json:"messages"` // References only
	Reference string                 `json:"$ref"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string          `json:"-"`
	ReferenceTo *OperationReply `json:"-"`
}

// Process processes the OperationReply to make it ready for code generation.
func (or *OperationReply) Process(name string, op *Operation, spec Specification) {
	// Prevent modification if nil
	if or == nil {
		return
	}

	// Set name
	or.Name = utils.UpperFirstLetter(name)

	// Add pointer to reference if there is one
	if or.Reference != "" {
		or.ReferenceTo = spec.ReferenceOperationReply(or.Reference)
	}

	// Process channel if there is one
	or.Channel.Process(name+ChannelSuffix, spec)

	// Process messages
	for i, msg := range or.Messages {
		msg.Process(fmt.Sprintf("%s%d", name, i), spec)
	}

	// Process address
	or.Address.Process(name+"Address", op, spec)
}
