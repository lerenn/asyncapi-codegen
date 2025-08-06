package asyncapiv3

// MessageTrait is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#messageTraitObject
type MessageTrait struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Headers *Schema `json:"headers,omitempty"`
	// NOTE: "payload" is explicitly disallowed: https://github.com/asyncapi/spec/blob/c5888f52b70f8eb99f782df9d23a88e1a4dce112/spec/asyncapi.md?plain=1#L1414
	Payload       *Schema                `json:"payload,omitempty"`
	CorrelationID *CorrelationID         `json:"correlationID,omitempty"`
	ContentType   string                 `json:"contentType,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Title         string                 `json:"title,omitempty"`
	Summary       string                 `json:"summary,omitempty"`
	Description   string                 `json:"description,omitempty"`
	Tags          []*Tag                 `json:"tags,omitempty"`
	ExternalDocs  *ExternalDocumentation `json:"externalDocs,omitempty"`
	Bindings      *MessageBindings       `json:"bindings,omitempty"`
	Examples      []*MessageExample      `json:"examples,omitempty"`
	Reference     string                 `json:"$ref,omitempty"`

	// --- Non AsyncAPI fields -------------------------------------------------

	ReferenceTo *MessageTrait `json:"-"`

	// CorrelationIDLocation will indicate where the correlation id is
	// According to: https://www.asyncapi.com/docs/reference/specification/v3.0.0#correlationIdObject
	CorrelationIDLocation string `json:"-"`
	CorrelationIDRequired bool   `json:"-"`
}

// generateMetadata generates metadata for the MessageTrait.
func (mt *MessageTrait) generateMetadata(parentName, name string, number *int) error {
	// Prevent modification if nil
	if mt == nil {
		return nil
	}

	// Set name
	mt.Name = generateFullName(parentName, name, "Trait", number)

	// Generate Headers metadata
	if err := mt.Headers.generateMetadata(mt.Name, "Headers", nil, false); err != nil {
		return err
	}

	// generate Payload metadata
	if err := mt.Payload.generateMetadata(mt.Name, "Payload", nil, false); err != nil {
		return err
	}

	// Generate tags metadata
	for i, t := range mt.Tags {
		t.generateMetadata(mt.Name, "", &i)
	}

	// Generate external documentation metadata
	mt.ExternalDocs.generateMetadata(mt.Name, ExternalDocsNameSuffix)

	// Generate Bindings metadata
	mt.Bindings.generateMetadata(mt.Name, "")

	// Generate Message Examples metadata
	for i, e := range mt.Examples {
		e.generateMetadata(mt.Name, "", &i)
	}

	return nil
}

// setDependencies sets dependencies between the different elements of the MessageTrait.
//
//nolint:cyclop
func (mt *MessageTrait) setDependencies(spec Specification) error {
	// Prevent modification if nil
	if mt == nil {
		return nil
	}

	// Set reference
	if err := mt.setReference(spec); err != nil {
		return err
	}

	// Set Headers and Payload depencies
	if err := mt.Headers.setDependencies(spec); err != nil {
		return err
	}
	if err := mt.Payload.setDependencies(spec); err != nil {
		return err
	}

	// Set tags dependencies
	for _, t := range mt.Tags {
		if err := t.setDependencies(spec); err != nil {
			return err
		}
	}

	// Set external documentation dependencies
	if err := mt.ExternalDocs.setDependencies(spec); err != nil {
		return err
	}

	// Set Bindings dependencies
	if err := mt.Bindings.setDependencies(spec); err != nil {
		return err
	}

	// Set Message Examples dependencies
	for _, e := range mt.Examples {
		if err := e.setDependencies(spec); err != nil {
			return err
		}
	}

	return nil
}

func (mt *MessageTrait) setReference(spec Specification) error {
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
