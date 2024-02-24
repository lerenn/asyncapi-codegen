package asyncapiv3

import "github.com/lerenn/asyncapi-codegen/pkg/utils"

// ExternalDocumentation is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#externalDocumentationObject
type ExternalDocumentation struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Description string `json:"description"`
	URL         string `json:"url"`
	Reference   string `json:"$ref"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string                 `json:"-"`
	ReferenceTo *ExternalDocumentation `json:"-"`
}

// Process processes the ExternalDocumentation to make it ready for code generation.
func (doc *ExternalDocumentation) Process(name string, spec Specification) {
	doc.Name = utils.UpperFirstLetter(name)

	// Add pointer to reference if there is one
	if doc.Reference != "" {
		doc.ReferenceTo = spec.ReferenceExternalDocumentation(doc.Reference)
	}
}
