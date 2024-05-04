package asyncapiv3

// Info is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#infoObject
type Info struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Title          string                 `json:"title"`
	Version        string                 `json:"version"`
	Description    string                 `json:"description"`
	TermsOfService string                 `json:"termsOfService"`
	Contact        *Contact               `json:"contact"`
	License        *License               `json:"license"`
	Tags           []*Tag                 `json:"tags"`
	ExternalDocs   *ExternalDocumentation `json:"externalDocs"`

	// --- Non AsyncAPI fields -------------------------------------------------
}

// generateMetadata generates metadata for the Info.
func (info *Info) generateMetadata(parentName string) error {
	// Prevent modification if nil
	if info == nil {
		return nil
	}

	// Process tags
	for i, t := range info.Tags {
		t.generateMetadata(parentName, "", &i)
	}

	// Process external documentation
	info.ExternalDocs.generateMetadata(parentName, ExternalDocsNameSuffix)

	return nil
}

// setDependencies sets dependencies between the different elements of the Info.
func (info *Info) setDependencies(spec Specification) error {
	// Prevent modification if nil
	if info == nil {
		return nil
	}

	// Process tags
	for _, t := range info.Tags {
		if err := t.setDependencies(spec); err != nil {
			return err
		}
	}

	// Process external documentation
	if err := info.ExternalDocs.setDependencies(spec); err != nil {
		return err
	}

	return nil
}
