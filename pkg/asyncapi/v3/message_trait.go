package asyncapiv3

import (
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/utils/template"
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
func (mt *MessageTrait) Process(name string, spec Specification) error {
	// Prevent modification if nil
	if mt == nil {
		return nil
	}

	// Set name
	if mt.Name == "" {
		mt.Name = template.Namify(name)
	} else {
		mt.Name = template.Namify(mt.Name)
	}

	// Process reference
	if err := mt.processReference(spec); err != nil {
		return err
	}

	// Process Headers and Payload
	if err := mt.Headers.Process(name+MessageHeadersSuffix, spec, false); err != nil {
		return err
	}
	if err := mt.Payload.Process(name+MessagePayloadSuffix, spec, false); err != nil {
		return err
	}

	// Process tags
	if err := mt.processTags(spec); err != nil {
		return err
	}

	// Process external documentation
	if err := mt.ExternalDocs.Process(mt.Name+ExternalDocsNameSuffix, spec); err != nil {
		return err
	}

	// Process Bindings
	if err := mt.Bindings.Process(mt.Name+BindingsSuffix, spec); err != nil {
		return err
	}

	// Process Message Examples
	if err := mt.processExamples(spec); err != nil {
		return err
	}

	return nil
}

func (mt *MessageTrait) processExamples(spec Specification) error {
	for i, e := range mt.Examples {
		if err := e.Process(fmt.Sprintf("%sExample%d", mt.Name, i), spec); err != nil {
			return err
		}
	}

	return nil
}

func (mt *MessageTrait) processTags(spec Specification) error {
	for i, t := range mt.Tags {
		if err := t.Process(fmt.Sprintf("%sTag%d", mt.Name, i), spec); err != nil {
			return err
		}
	}

	return nil
}

func (mt *MessageTrait) processReference(spec Specification) error {
	// check reference exists
	if mt.Reference == "" {
		return nil
	}

	// Add pointer to reference if there is one
	refTo, err := spec.ReferenceMessageTrait(mt.Reference)
	if err != nil {
		return err
	}
	mt.ReferenceTo = refTo

	return nil
}

// Follow returns referenced MessageTrait if specified or the actual MessageTrait.
func (mt *MessageTrait) Follow() *MessageTrait {
	if mt.ReferenceTo != nil {
		return mt.ReferenceTo
	}
	return mt
}
