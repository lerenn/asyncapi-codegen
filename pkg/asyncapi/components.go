package asyncapi

type Components struct {
	Messages map[string]Message `json:"messages"`
	Schemas  map[string]Any     `json:"schemas"`
}

func (c *Components) setMapsValuesName() {
	for name, msg := range c.Messages {
		msg.Name = name
		c.Messages[name] = msg
	}

	for name, schema := range c.Schemas {
		schema.Name = name
		c.Schemas[name] = schema
	}
}

func (c *Components) Process(spec Specification) {
	c.setMapsValuesName()

	for key, msg := range c.Messages {
		msg.Process(spec)
		c.Messages[key] = msg
	}
}
