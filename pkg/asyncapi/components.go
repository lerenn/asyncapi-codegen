package asyncapi

type Components struct {
	Messages map[string]*Message `json:"messages"`
	Schemas  map[string]*Any     `json:"schemas"`
}

func (c *Components) Process(spec Specification) {
	// For all messages, process message
	for name, msg := range c.Messages {
		msg.Process(name, spec)
	}

	// For all schemas, process schema
	for name, schema := range c.Schemas {
		schema.Process(name, spec)
	}
}
