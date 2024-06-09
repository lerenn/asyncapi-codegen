package asyncapiv3

const (
	// ExternalDocsNameSuffix is the suffix that is added to the name of external docs.
	ExternalDocsNameSuffix = "External_Docs"
)

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

// generateMetadata generates metadata for the ExternalDocumentation.
func (doc *ExternalDocumentation) generateMetadata(parentName, name string) {
	// Return if empty
	if doc == nil {
		return
	}

	// Set name
	doc.Name = generateFullName(parentName, name, ExternalDocsNameSuffix, nil)
}

// setDependencies sets dependencies between the different elements of the ExternalDocumentation.
func (doc *ExternalDocumentation) setDependencies(spec Specification) error {
	// Return if empty
	if doc == nil {
		return nil
	}

	// Add pointer to reference if there is one
	if doc.Reference != "" {
		refTo, err := spec.ReferenceExternalDocumentation(doc.Reference)
		if err != nil {
			return err
		}
		doc.ReferenceTo = refTo
	}

	return nil
}
