package asyncapiv2

// Components is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v2.6.0#componentsObject
type Components struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Messages   map[string]*Message   `json:"messages"`
	Schemas    map[string]*Schema    `json:"schemas"`
	Parameters map[string]*Parameter `json:"parameters"`

	// --- Non AsyncAPI fields -------------------------------------------------
}

// generateMetadata generate metadata for the components and its children.
func (c *Components) generateMetadata() error {
	// For all schemas, generate schema metadata
	for name, schema := range c.Schemas {
		if err := schema.generateMetadata(name+"_Schema", false); err != nil {
			return err
		}
	}

	// For all messages, generate message metadata
	for name, msg := range c.Messages {
		if err := msg.generateMetadata(name + MessageSuffix); err != nil {
			return err
		}
	}

	// For all parameters, generate param metadata
	for name, param := range c.Parameters {
		param.generateMetadata(name + "_Parameter")
	}

	return nil
}

// setDependencies set dependencies for the components and its children from specification.
func (c *Components) setDependencies(spec Specification) error {
	// For all schemas, set schema dependencies
	for _, schema := range c.Schemas {
		if err := schema.setDependencies(spec); err != nil {
			return err
		}
	}

	// For all messages, set message dependencies
	for _, msg := range c.Messages {
		if err := msg.setDependencies(spec); err != nil {
			return err
		}
	}

	// For all parameters, set param dependencies
	for _, param := range c.Parameters {
		if err := param.setDependencies(spec); err != nil {
			return err
		}
	}

	return nil
}
