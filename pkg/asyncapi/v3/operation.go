package asyncapiv3

// OperationAction represents an OperationAction.
type OperationAction string

const (
	// OperationActionIsSend represents a send action.
	OperationActionIsSend OperationAction = "send"
	// OperationActionIsReceive represents a receive action.
	OperationActionIsReceive OperationAction = "receive"
)

// IsSend returns true if the operation action is send.
func (oa OperationAction) IsSend() bool {
	return oa == OperationActionIsSend
}

// IsReceive returns true if the operation action is receive.
func (oa OperationAction) IsReceive() bool {
	return oa == OperationActionIsReceive
}

// Operation is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#operationObject
//
// NOTE: From AsyncAPI specification on the "messages" field:
//
//	Excluding this property from the Operation implies that all messages from the channel will be included. Explicitly set the messages property to [] if this operation should contain no messages.
//
// Because of this caveat, it must be possible to both serialize the messages slice as an empty JSON array or exclude the key entirely.
// The omitzero tag is required to enable this.
type Operation struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Action       OperationAction        `json:"action"`
	Channel      *Channel               `json:"channel"` // Reference only
	Title        string                 `json:"title,omitempty"`
	Summary      string                 `json:"summary,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Security     []*SecurityScheme      `json:"security,omitempty"`
	Tags         []*Tag                 `json:"tags,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
	Bindings     *OperationBindings     `json:"bindings,omitempty"`
	Traits       []*OperationTrait      `json:"traits,omitempty"`
	Messages     []*Message             `json:"messages,omitzero"` // References only
	Reply        *OperationReply        `json:"reply,omitempty"`
	Reference    string                 `json:"$ref,omitempty"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string     `json:"-"`
	ReplyIs     *Operation `json:"-"`
	ReplyOf     *Operation `json:"-"`
	ReferenceTo *Operation `json:"-"`
}

func (op *Operation) generateMetadata(parentName, name string) error {
	// Prevent modification if nil
	if op == nil {
		return nil
	}

	// Set name
	op.Name = generateFullName(parentName, name, "Operation", nil)

	// Generate securities metadata
	for i, sec := range op.Security {
		sec.generateMetadata(op.Name, "", &i)
	}

	// Generate external doc metadata if there is one
	op.ExternalDocs.generateMetadata(op.Name, ExternalDocsNameSuffix)

	// Generate bindings metadata if there is one
	op.Bindings.generateMetadata(op.Name, "")

	// Generate traits metadata
	for i, t := range op.Traits {
		t.generateMetadata(op.Name, "", &i)
	}

	// Generate reply metadata if there is one
	if err := op.Reply.generateMetadata(op.Name, ""); err != nil {
		return err
	}

	return nil
}

// setDependencies sets dependencies between the different elements of the Operation.
//
//nolint:cyclop
func (op *Operation) setDependencies(spec Specification) error {
	// Prevent modification if nil
	if op == nil {
		return nil
	}

	// Set reference
	if err := op.setReference(spec); err != nil {
		return err
	}

	// Set channel dependencies if there is one
	if err := op.Channel.setDependencies(spec); err != nil {
		return err
	}

	// Set securities dependencies
	for _, sec := range op.Security {
		if err := sec.setDependencies(spec); err != nil {
			return err
		}
	}

	// Set external doc dependencies if there is one
	if err := op.ExternalDocs.setDependencies(spec); err != nil {
		return err
	}

	// Set bindings dependencies if there is one
	if err := op.Bindings.setDependencies(spec); err != nil {
		return err
	}

	// Set traits dependencies and apply them
	if err := op.setTraitsDependenciesAndApply(spec); err != nil {
		return err
	}

	// Set messages dependencies
	for _, msg := range op.Messages {
		if err := msg.setDependencies(spec); err != nil {
			return err
		}
	}

	// Set reply dependencies if there is one
	if err := op.Reply.setDependencies(op, spec); err != nil {
		return err
	}

	// Generate reply
	op.generateReply()

	return nil
}

func (op *Operation) setReference(spec Specification) error {
	if op.Reference == "" {
		return nil
	}

	// Add pointer to reference if there is one
	refTo, err := spec.ReferenceOperation(op.Reference)
	if err != nil {
		return err
	}
	op.ReferenceTo = refTo

	return nil
}

func (op *Operation) setTraitsDependenciesAndApply(spec Specification) error {
	for _, t := range op.Traits {
		if err := t.setDependencies(spec); err != nil {
			return err
		}

		op.ApplyTrait(t.Follow(), spec)
	}

	return nil
}

func (op *Operation) generateReply() {
	// Return if there is no reply
	if op == nil || op.Reply == nil {
		return
	}

	// Generate reply
	ch := op.Reply.Channel.Follow()
	op.ReplyIs = &Operation{
		Name:    "ReplyTo" + op.Name,
		Channel: ch,
		ReplyOf: op,
	}
}

// GetMessage will return the operation message.
func (op Operation) GetMessage() (*Message, error) {
	if len(op.Messages) > 0 {
		return op.Messages[0], nil // TODO: change
	}

	return op.Channel.GetMessage()
}

// ApplyTrait applies a trait to the operation.
func (op *Operation) ApplyTrait(ot *OperationTrait, spec Specification) {
	// Check operation is not nil
	if op == nil {
		return
	}

	// Override title if not set
	if op.Title == "" {
		op.Title = ot.Title
	}

	// Override summary if not set
	if op.Summary == "" {
		op.Summary = ot.Summary
	}

	// Override description if not set
	if op.Description == "" {
		op.Description = ot.Description
	}

	// Merge security scheme
	op.Security = append(op.Security, ot.Security...)
	op.Security = RemoveDuplicateSecuritySchemes(op.Security)

	// Merge tags
	op.Tags = append(op.Tags, ot.Tags...)
	op.Tags = RemoveDuplicateTags(op.Tags)

	// Override external docs if not set
	if op.ExternalDocs == nil && ot.ExternalDocs != nil {
		extDoc := *ot.ExternalDocs
		op.ExternalDocs = &extDoc
	}

	// Override bindings if not set
	if op.Bindings == nil && ot.Bindings != nil {
		bindings := *ot.Bindings
		op.Bindings = &bindings
	}
}

// Follow returns referenced operation if specified or the actual operation.
func (op *Operation) Follow() *Operation {
	if op.ReferenceTo != nil {
		return op.ReferenceTo
	}
	return op
}
