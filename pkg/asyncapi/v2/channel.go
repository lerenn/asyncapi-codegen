package asyncapiv2

import "github.com/lerenn/asyncapi-codegen/pkg/utils/template"

// Channel is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v2.6.0#channelItemObject
type Channel struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Parameters map[string]*Parameter `json:"parameters"`

	Subscribe *Operation `json:"subscribe"`
	Publish   *Operation `json:"publish"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name string `json:"-"`
	Path string `json:"-"`
}

// Process processes the Channel to make it ready for code generation.
func (c *Channel) Process(path string, spec Specification) error {
	// Set channel name and path
	c.Name = template.Namify(path)
	c.Path = path

	// If there is publish and subscribe, add suffix to avoid duplicate names
	suffixPublish, suffixSubscribe := "", ""
	if c.Subscribe != nil && c.Publish != nil {
		suffixPublish = "Publish"
		suffixSubscribe = "Subscribe"
	}

	// Process subscribe operation
	if c.Subscribe != nil {
		if err := c.Subscribe.Process(c.Name+suffixSubscribe, spec); err != nil {
			return err
		}
	}

	// Process publish operation
	if c.Publish != nil {
		if err := c.Publish.Process(c.Name+suffixPublish, spec); err != nil {
			return err
		}
	}

	// Process parameters
	for n, p := range c.Parameters {
		if err := p.Process(n, spec); err != nil {
			return err
		}
	}

	return nil
}
