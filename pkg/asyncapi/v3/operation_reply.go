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

// generateMetadata generates metadata for the OperationReply.
func (or *OperationReply) generateMetadata(name string) error {
	// Prevent modification if nil
	if or == nil {
		return nil
	}

	// Set name
	or.Name = template.Namify(name)

	// Generate channel metadata if there is one
	if err := or.Channel.generateMetadata(name + ChannelSuffix); err != nil {
		return err
	}

	// Generate messages metadata
	for i, msg := range or.Messages {
		if err := msg.generateMetadata(fmt.Sprintf("%s%d", name, i)); err != nil {
			return err
		}
	}

	// Generate address metadata
	or.Address.generateMetadata(name + "Address")

	return nil
}

// setDependencies sets dependencies between the different elements of the OperationReply.
func (or *OperationReply) setDependencies(op *Operation, spec Specification) error {
	// Prevent modification if nil
	if or == nil {
		return nil
	}

	// Add pointer to reference if there is one
	if or.Reference != "" {
		refTo, err := spec.ReferenceOperationReply(or.Reference)
		if err != nil {
			return err
		}
		or.ReferenceTo = refTo
	}

	// Set channel dependencies if there is one
	if err := or.Channel.setDependencies(spec); err != nil {
		return err
	}

	// Set messages dependencies
	for _, msg := range or.Messages {
		if err := msg.setDependencies(spec); err != nil {
			return err
		}
	}

	// Set address dependencies
	if err := or.Address.setDependencies(op, spec); err != nil {
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
