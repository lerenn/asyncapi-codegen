package asyncapiv3

import (
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/utils"
)

// MessageTrait is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#messageTraitObject
type MessageTrait struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Headers       *Schema                `json:"headers"`
	Payload       *Schema                `json:"payload"`
	CorrelationID *CorrelationID         `json:"correlationID"`
	ContentType   string                 `json:"contentType"`
	Name          string                 `json:"name"`
	Title         string                 `json:"title"`
	Summary       string                 `json:"summary"`
	Description   string                 `json:"description"`
	Tags          []*Tag                 `json:"tags"`
	ExternalDocs  *ExternalDocumentation `json:"externalDocs"`
	Bindings      *MessageBindings       `json:"bindings"`
	Examples      []*MessageExample      `json:"examples"`
	Reference     string                 `json:"$ref"`

	// --- Non AsyncAPI fields -------------------------------------------------

	ReferenceTo *MessageTrait `json:"-"`

	// CorrelationIDLocation will indicate where the correlation id is
	// According to: https://www.asyncapi.com/docs/reference/specification/v3.0.0#correlationIdObject
	CorrelationIDLocation string `json:"-"`
	CorrelationIDRequired bool   `json:"-"`
}

// Process processes the MessageTrait to make it ready for code generation.
func (mt *MessageTrait) Process(name string, spec Specification) {
	// Prevent modification if nil
	if mt == nil {
		return
	}

	// Set name
	if mt.Name == "" {
		mt.Name = utils.UpperFirstLetter(name)
	} else {
		mt.Name = utils.UpperFirstLetter(mt.Name)
	}

	// Add pointer to reference if there is one
	if mt.Reference != "" {
		mt.ReferenceTo = spec.ReferenceMessageTrait(mt.Reference)
	}

	// Process Headers and Payload
	mt.Headers.Process(name+"Headers", spec, false)
	mt.Payload.Process(name+"Payload", spec, false)

	// Process tags
	for i, t := range mt.Tags {
		t.Process(fmt.Sprintf("%sTag%d", mt.Name, i), spec)
	}

	// Process external documentation
	mt.ExternalDocs.Process(mt.Name+ExternalDocsNameSuffix, spec)

	// Process Bindings
	mt.Bindings.Process(mt.Name+BindingsSuffix, spec)

	// Process Message Examples
	for i, ex := range mt.Examples {
		ex.Process(fmt.Sprintf("%sExample%d", mt.Name, i), spec)
	}
}
