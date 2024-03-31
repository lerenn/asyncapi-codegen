package asyncapiv3

import (
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/utils/template"
)

// OperationTrait is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#msgerationTraitObject
type OperationTrait struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Title        string                 `json:"title"`
	Summary      string                 `json:"summary"`
	Description  string                 `json:"description"`
	Security     []*SecurityScheme      `json:"security"`
	Tags         []*Tag                 `json:"tags"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs"`
	Bindings     *OperationBindings     `json:"bindings"`
	Reference    string                 `json:"$ref"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string          `json:"-"`
	ReferenceTo *OperationTrait `json:"-"`
}

// Process processes the OperationTrait to make it ready for code generation.
func (ot *OperationTrait) Process(name string, spec Specification) error {
	// Prevent modification if nil
	if ot == nil {
		return nil
	}

	// Set name
	ot.Name = template.Namify(name)

	// Add pointer to reference if there is one
	if ot.Reference != "" {
		refTo, err := spec.ReferenceOperationTrait(ot.Reference)
		if err != nil {
			return err
		}
		ot.ReferenceTo = refTo
	}

	// Process securities
	for i, s := range ot.Security {
		if err := s.Process(fmt.Sprintf("%sSecurity%d", name, i), spec); err != nil {
			return err
		}
	}

	// Process external doc if there is one
	if err := ot.ExternalDocs.Process(name+ExternalDocsNameSuffix, spec); err != nil {
		return err
	}

	// Process tags
	for i, t := range ot.Tags {
		if err := t.Process(fmt.Sprintf("%sTag%d", ot.Name, i), spec); err != nil {
			return err
		}
	}

	// Process bindings if there is one
	return ot.Bindings.Process(name+BindingsSuffix, spec)
}

// Follow returns referenced MessageTrait if specified or the actual MessageTrait.
func (ot *OperationTrait) Follow() *OperationTrait {
	if ot.ReferenceTo != nil {
		return ot.ReferenceTo
	}
	return ot
}
