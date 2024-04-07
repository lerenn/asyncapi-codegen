package asyncapiv3

import (
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/utils/template"
)

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
type Operation struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Action       OperationAction        `json:"action"`
	Channel      *Channel               `json:"channel"` // Reference only
	Title        string                 `json:"title"`
	Summary      string                 `json:"summary"`
	Description  string                 `json:"string"`
	Security     []*SecurityScheme      `json:"security"`
	Tags         []*Tag                 `json:"tags"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs"`
	Bindings     *OperationBindings     `json:"bindings"`
	Traits       []*OperationTrait      `json:"traits"`
	Messages     []*Message             `json:"messages"` // References only
	Reply        *OperationReply        `json:"reply"`
	Reference    string                 `json:"$ref"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string     `json:"-"`
	ReplyIs     *Operation `json:"-"`
	ReplyOf     *Operation `json:"-"`
	ReferenceTo *Operation `json:"-"`
}

// Process processes the Channel to make it ready for code generation.
func (op *Operation) Process(name string, spec Specification) error {
	// Prevent modification if nil
	if op == nil {
		return nil
	}

	// Set name
	op.Name = template.Namify(name)

	// Process reference
	if err := op.processReference(spec); err != nil {
		return err
	}

	// Process channel if there is one
	if err := op.Channel.Process(op.Name+ChannelSuffix, spec); err != nil {
		return err
	}

	// Process securities
	if err := op.processSecurities(spec); err != nil {
		return err
	}

	// Process external doc if there is one
	if err := op.ExternalDocs.Process(op.Name+ExternalDocsNameSuffix, spec); err != nil {
		return err
	}

	// Process bindings if there is one
	if err := op.Bindings.Process(op.Name+BindingsSuffix, spec); err != nil {
		return err
	}

	// Process traits and apply them
	if err := op.processTraits(spec); err != nil {
		return err
	}

	// Process messages
	if err := op.processMessages(spec); err != nil {
		return err
	}

	// Process reply if there is one
	if err := op.Reply.Process(op.Name+"Reply", op, spec); err != nil {
		return err
	}

	// Generate reply
	op.generateReply()

	return nil
}

func (op *Operation) processSecurities(spec Specification) error {
	for i, sec := range op.Security {
		if err := sec.Process(fmt.Sprintf("%sSecurity%d", op.Name, i), spec); err != nil {
			return err
		}
	}

	return nil
}

func (op *Operation) processMessages(spec Specification) error {
	for i, msg := range op.Messages {
		if err := msg.Process(fmt.Sprintf("%sMessage%d", op.Name, i), spec); err != nil {
			return err
		}
	}

	return nil
}

func (op *Operation) processReference(spec Specification) error {
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

func (op *Operation) processTraits(spec Specification) error {
	for i, t := range op.Traits {
		if err := t.Process(fmt.Sprintf("%sTrait%d", op.Name, i), spec); err != nil {
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
