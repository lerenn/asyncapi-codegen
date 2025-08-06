package asyncapiv3

// OperationTrait is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#msgerationTraitObject
type OperationTrait struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Title        string                 `json:"title,omitempty"`
	Summary      string                 `json:"summary,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Security     []*SecurityScheme      `json:"security,omitempty"`
	Tags         []*Tag                 `json:"tags,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
	Bindings     *OperationBindings     `json:"bindings,omitempty"`
	Reference    string                 `json:"$ref,omitempty"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string          `json:"-"`
	ReferenceTo *OperationTrait `json:"-"`
}

// generateMetadata generates metadata for the OperationTrait.
func (ot *OperationTrait) generateMetadata(parentName, name string, number *int) {
	// Prevent modification if nil
	if ot == nil {
		return
	}

	// Set name
	ot.Name = generateFullName(parentName, name, "Trait", number)

	// Generate securities metadata
	for i, s := range ot.Security {
		s.generateMetadata(ot.Name, "", &i)
	}

	// Generate external doc metadata if there is one
	ot.ExternalDocs.generateMetadata(ot.Name, ExternalDocsNameSuffix)

	// Generate tags metadata
	for i, t := range ot.Tags {
		t.generateMetadata(ot.Name, "", &i)
	}

	// Generate bindings metadata if there is one
	ot.Bindings.generateMetadata(ot.Name, "")
}

// setDependencies sets dependencies for the OperationTrait.
func (ot *OperationTrait) setDependencies(spec Specification) error {
	// Prevent modification if nil
	if ot == nil {
		return nil
	}

	// Add pointer to reference if there is one
	if ot.Reference != "" {
		refTo, err := spec.ReferenceOperationTrait(ot.Reference)
		if err != nil {
			return err
		}
		ot.ReferenceTo = refTo
	}

	// Set securities dependencies
	for _, s := range ot.Security {
		if err := s.setDependencies(spec); err != nil {
			return err
		}
	}

	// Set external doc dependencies if there is one
	if err := ot.ExternalDocs.setDependencies(spec); err != nil {
		return err
	}

	// Set tags dependencies
	for _, t := range ot.Tags {
		if err := t.setDependencies(spec); err != nil {
			return err
		}
	}

	// Set bindings dependencies if there is one
	return ot.Bindings.setDependencies(spec)
}

// Follow returns referenced MessageTrait if specified or the actual MessageTrait.
func (ot *OperationTrait) Follow() *OperationTrait {
	if ot.ReferenceTo != nil {
		return ot.ReferenceTo
	}
	return ot
}
