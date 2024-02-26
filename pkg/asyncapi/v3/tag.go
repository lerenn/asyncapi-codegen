package asyncapiv3

import "github.com/lerenn/asyncapi-codegen/pkg/utils"

// Tag is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#tagsObject
type Tag struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs"`
	Reference    string                 `json:"$ref"`

	// --- Non AsyncAPI fields -------------------------------------------------

	ReferenceTo *Tag `json:"-"`
}

// Process processes the Tag to make it ready for code generation.
func (t *Tag) Process(path string, spec Specification) {
	// Prevent modification if nil
	if t == nil {
		return
	}

	// Set name
	t.Name = utils.UpperFirstLetter(path)

	// Add pointer to reference if there is one
	if t.Reference != "" {
		t.ReferenceTo = spec.ReferenceTag(t.Reference)
	}
}
