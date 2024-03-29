package asyncapiv3

import (
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/utils/template"
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
func (or *OperationReply) Process(name string, op *Operation, spec Specification) error {
	// Prevent modification if nil
	if or == nil {
		return nil
	}

	// Set name
	or.Name = template.Namify(name)

	// Add pointer to reference if there is one
	if or.Reference != "" {
		refTo, err := spec.ReferenceOperationReply(or.Reference)
		if err != nil {
			return err
		}
		or.ReferenceTo = refTo
	}

	// Process channel if there is one
	if err := or.Channel.Process(name+ChannelSuffix, spec); err != nil {
		return err
	}

	// Process messages
	for i, msg := range or.Messages {
		if err := msg.Process(fmt.Sprintf("%s%d", name, i), spec); err != nil {
			return err
		}
	}

	// Process address
	if err := or.Address.Process(name+"Address", op, spec); err != nil {
		return err
	}

	return nil
}

// Follow returns referenced operation if specified or the actual operation.
func (or *OperationReply) Follow() *OperationReply {
	if or.ReferenceTo != nil {
		return or.ReferenceTo
	}
	return or
}
