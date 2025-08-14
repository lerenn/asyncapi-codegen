package asyncapiv3

// Tag is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#tagsObject
type Tag struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Name         string                 `json:"name"`
	Description  string                 `json:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
	Reference    string                 `json:"$ref,omitempty"`

	// --- Non AsyncAPI fields -------------------------------------------------

	ReferenceTo *Tag `json:"-"`
}

// generateMetadata generates metadata for the Tag.
func (t *Tag) generateMetadata(parentName, name string, number *int) {
	// Prevent modification if nil
	if t == nil {
		return
	}

	// Set name
	t.Name = generateFullName(parentName, name, "Tag", number)

	// Generate ExternalDocs metadata
	t.ExternalDocs.generateMetadata(name, ExternalDocsNameSuffix)
}

// setDependencies sets dependencies between the different elements of the Tag.
func (t *Tag) setDependencies(spec Specification) error {
	// Prevent modification if nil
	if t == nil {
		return nil
	}

	// Process external documentation
	if err := t.ExternalDocs.setDependencies(spec); err != nil {
		return err
	}

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
