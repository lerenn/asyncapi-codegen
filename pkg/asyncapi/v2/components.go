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

// Process processes the Components structure to make it ready for code generation.
func (c *Components) Process(spec Specification) error {
	// For all schemas, process schema
	for name, schema := range c.Schemas {
		if err := schema.Process(name+"Schema", spec, false); err != nil {
			return err
		}
	}

	// For all messages, process message
	for name, msg := range c.Messages {
		if err := msg.Process(name+MessageSuffix, spec); err != nil {
			return err
		}
	}

	// For all parameters, process param
	for name, param := range c.Parameters {
		if err := param.Process(name+"Parameter", spec); err != nil {
			return err
		}
	}

	return nil
}
