package asyncapiv3

import (
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/utils"
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
	Traits       *OperationTrait        `json:"traits"`
	Messages     []*Message             `json:"messages"` // References only
	Reply        *OperationReply        `json:"reply"`
	Reference    string                 `json:"$ref"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string     `json:"-"`
	IsReplyTo   *Operation `json:"-"`
	ReferenceTo *Operation `json:"-"`
}

// Process processes the Channel to make it ready for code generation.
func (op *Operation) Process(name string, spec Specification) {
	// Prevent modification if nil
	if op == nil {
		return
	}

	// Set name
	op.Name = utils.UpperFirstLetter(name)

	// Add pointer to reference if there is one
	if op.Reference != "" {
		op.ReferenceTo = spec.ReferenceOperation(op.Reference)
	}

	// Process channel if there is one
	op.Channel.Process(name+"Channel", spec)

	// Process securities
	for i, s := range op.Security {
		s.Process(fmt.Sprintf("%sSecurity%d", name, i), spec)
	}

	// Process external doc if there is one
	op.ExternalDocs.Process(name+ExternalDocsNameSuffix, spec)

	// Process bindings if there is one
	op.Bindings.Process(name+BindingsSuffix, spec)

	// Process traits if there is one
	op.Traits.Process(name+"Traits", spec)

	// Process messages
	for i, msg := range op.Messages {
		msg.Process(fmt.Sprintf("%s%d", name, i), spec)
	}

	// Process reply if there is one
	op.Reply.Process(name+"Reply", spec)
}

// GetMessage will return the operation message.
func (op Operation) GetMessage() *Message {
	if len(op.Messages) > 0 {
		return op.Messages[0] // TODO: change
	} else {
		return op.Channel.GetMessage()
	}
}
