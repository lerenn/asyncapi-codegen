package asyncapiv3

// OperationReply is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#operationReplyObject
type OperationReply struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Address   *OperationReplyAddress `json:"address,omitempty"`
	Channel   *Channel               `json:"channel,omitempty"`  // Reference only
	Messages  []*Message             `json:"messages,omitempty"` // References only
	Reference string                 `json:"$ref,omitempty"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string          `json:"-"`
	ReferenceTo *OperationReply `json:"-"`
}

// generateMetadata generates metadata for the OperationReply.
func (or *OperationReply) generateMetadata(parentName, name string) error {
	// Prevent modification if nil
	if or == nil {
		return nil
	}

	// Set name
	or.Name = generateFullName(parentName, name, "Reply", nil)

	// Generate address metadata
	or.Address.generateMetadata(or.Name, "")

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
