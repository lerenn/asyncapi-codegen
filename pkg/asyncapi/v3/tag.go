package asyncapiv3

import "github.com/lerenn/asyncapi-codegen/pkg/utils/template"

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
func (t *Tag) Process(path string, spec Specification) error {
	// Prevent modification if nil
	if t == nil {
		return nil
	}

	// Set name
	t.Name = template.Namify(path)

	// Add pointer to reference if there is one
	if t.Reference != "" {
		refTo, err := spec.ReferenceTag(t.Reference)
		if err != nil {
			return err
		}
		t.ReferenceTo = refTo
	}

	return nil
}

// RemoveDuplicateTags removes the tags that have the same name, keeping the first occurrence.
func RemoveDuplicateTags(tags []*Tag) []*Tag {
	newList := make([]*Tag, 0)
	for _, t := range tags {
		present := false
		for _, pt := range newList {
			if pt.Name == t.Name {
				present = true
				break
			}
		}

		if !present {
			newList = append(newList, t)
		}
	}
	return newList
}
