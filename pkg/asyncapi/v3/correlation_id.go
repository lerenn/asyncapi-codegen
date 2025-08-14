package asyncapiv3

// CorrelationID is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#correlationIdObject
type CorrelationID struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Description string `json:"description,omitempty"`
	Location    string `json:"location,omitempty"`
	Reference   string `json:"$ref,omitempty"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string   `json:"-"`
	ReferenceTo *Channel `json:"-"`
}

// generateMetadata generates metadata for the CorrelationID.
func (c *CorrelationID) generateMetadata(parentName, name string) {
	// Prevent modification if nil
	if c == nil {
		return
	}

	// Set name
	c.Name = generateFullName(parentName, name, "Correlation_ID", nil)
}

// setDependencies sets dependencies between the different elements of the CorrelationID.
func (c *CorrelationID) setDependencies(spec Specification) error {
	// Prevent modification if nil
	if c == nil {
		return nil
	}

	// Add pointer to reference if there is one
	if c.Reference != "" {
		refTo, err := spec.ReferenceChannel(c.Reference)
		if err != nil {
			return err
		}
		c.ReferenceTo = refTo
	}

	return nil
}

// Exists checks that the correlation exists (and that the location is set).
func (c *CorrelationID) Exists() bool {
	return c != nil && c.Location != ""
}
