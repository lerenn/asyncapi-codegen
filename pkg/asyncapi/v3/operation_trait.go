package asyncapiv3

import (
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/utils"
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
func (ot *OperationTrait) Process(name string, spec Specification) {
	// Prevent modification if nil
	if ot == nil {
		return
	}

	// Set name
	ot.Name = utils.UpperFirstLetter(name)

	// Add pointer to reference if there is one
	if ot.Reference != "" {
		ot.ReferenceTo = spec.ReferenceOperationTrait(ot.Reference)
	}

	// Process securities
	for i, s := range ot.Security {
		s.Process(fmt.Sprintf("%sSecurity%d", name, i), spec)
	}

	// Process external doc if there is one
	ot.ExternalDocs.Process(name+ExternalDocsNameSuffix, spec)

	// Process tags
	for i, t := range ot.Tags {
		t.Process(fmt.Sprintf("%sTag%d", ot.Name, i), spec)
	}

	// Process bindings if there is one
	ot.Bindings.Process(name+BindingsSuffix, spec)
}

// Follow returns referenced MessageTrait if specified or the actual MessageTrait.
func (ot *OperationTrait) Follow() *OperationTrait {
	if ot.ReferenceTo != nil {
		return ot.ReferenceTo
	}
	return ot
}
