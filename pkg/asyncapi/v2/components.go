package asyncapiv2

// Components is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v2.6.0#componentsObject
type Components struct {
	Messages   map[string]*Message   `json:"messages"`
	Schemas    map[string]*Schema    `json:"schemas"`
	Parameters map[string]*Parameter `json:"parameters"`
}

// Process processes the Components structure to make it ready for code generation.
func (c *Components) Process(spec Specification) {
	// For all schemas, process schema
	for name, schema := range c.Schemas {
		schema.Process(name, spec, false)
	}

	// For all messages, process message
	for name, msg := range c.Messages {
		msg.Process(name, spec)
	}

	// For all parameters, process param
	for name, param := range c.Parameters {
		param.Process(name, spec)
	}
}
