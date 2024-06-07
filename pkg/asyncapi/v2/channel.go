package asyncapiv2

import "github.com/TheSadlig/asyncapi-codegen/pkg/utils/template"

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

// generateMetadata generate metadata for the channel and its children.
func (c *Channel) generateMetadata(path string) error {
	// Set channel name and path
	c.Name = template.Namify(path)
	c.Path = path

	// If there is publish and subscribe, add suffix to avoid duplicate names
	suffixPublish, suffixSubscribe := "", ""
	if c.Subscribe != nil && c.Publish != nil {
		suffixPublish = "Publish"
		suffixSubscribe = "Subscribe"
	}

	// Generate subscribe operation metadata
	if c.Subscribe != nil {
		if err := c.Subscribe.generateMetadata(c.Name + suffixSubscribe); err != nil {
			return err
		}
	}

	// Generate publish operation metadata
	if c.Publish != nil {
		if err := c.Publish.generateMetadata(c.Name + suffixPublish); err != nil {
			return err
		}
	}

	// Generate parameters metadata
	for n, p := range c.Parameters {
		p.generateMetadata(n)
	}

	return nil
}

// setDependencies set dependencies for the channel and its children from specification.
func (c *Channel) setDependencies(spec Specification) error {
	// Set subscribe operation dependencies if present
	if c.Subscribe != nil {
		if err := c.Subscribe.setDependencies(spec); err != nil {
			return err
		}
	}

	// Set publish operation dependencies if present
	if c.Publish != nil {
		if err := c.Publish.setDependencies(spec); err != nil {
			return err
		}
	}

	// Set parameters dependencies
	for _, p := range c.Parameters {
		if err := p.setDependencies(spec); err != nil {
			return err
		}
	}

	return nil
}
