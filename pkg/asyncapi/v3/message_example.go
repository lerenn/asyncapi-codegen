package asyncapiv3

// MessageExample is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#messageExampleObject
type MessageExample struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Headers   map[string]any `json:"headers,omitempty"`
	Payload   map[string]any `json:"payload,omitempty"`
	Name      string         `json:"name,omitempty"`
	Summary   string         `json:"summary,omitempty"`
	Reference string         `json:"$ref,omitempty"`

	// --- Non AsyncAPI fields -------------------------------------------------

	ReferenceTo *MessageExample `json:"-"`
}

// generateMetadata generates metadata for the MessageExample.
func (me *MessageExample) generateMetadata(parentName, name string, number *int) {
	// Prevent modification if nil
	if me == nil {
		return
	}

	// Set name
	me.Name = generateFullName(parentName, name, "Example", number)
}

// setDependencies sets dependencies between the different elements of the MessageExample.
func (me *MessageExample) setDependencies(spec Specification) error {
	// Prevent modification if nil
	if me == nil {
		return nil
	}

	// Add pointer to reference if there is one
	if me.Reference != "" {
		refTo, err := spec.ReferenceMessageExample(me.Reference)
		if err != nil {
			return err
		}
		me.ReferenceTo = refTo
	}

	return nil
}
