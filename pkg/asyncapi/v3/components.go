package asyncapiv3

import "github.com/lerenn/asyncapi-codegen/pkg/utils/template"

// Components is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#componentsObject
type Components struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Schemas           map[string]*Schema                `json:"schemas"`
	Servers           map[string]*Server                `json:"servers"`
	Channels          map[string]*Channel               `json:"channels"`
	Operations        map[string]*Operation             `json:"operations"`
	Messages          map[string]*Message               `json:"messages"`
	SecuritySchemes   map[string]*SecurityScheme        `json:"securitySchemes"`
	ServerVariables   map[string]*ServerVariable        `json:"serverVariables"`
	Parameters        map[string]*Parameter             `json:"parameters"`
	CorrelationIDs    map[string]*CorrelationID         `json:"correlationIds"`
	Replies           map[string]*OperationReply        `json:"replies"`
	ReplyAddresses    map[string]*OperationReplyAddress `json:"replyAddresses"`
	ExternalDocs      map[string]*ExternalDocumentation `json:"externalDocs"`
	Tags              map[string]*Tag                   `json:"tags"`
	OperationTraits   map[string]*OperationTrait        `json:"operationTraits"`
	MessageTraits     map[string]*MessageTrait          `json:"messageTraits"`
	ServerBindings    map[string]*ServerBindings        `json:"serverBindings"`
	ChannelBindings   map[string]*ChannelBindings       `json:"channelBindings"`
	OperationBindings map[string]*OperationBindings     `json:"operationBindings"`
	MessageBindings   map[string]*MessageBindings       `json:"messageBindings"`

	// --- Non AsyncAPI fields -------------------------------------------------
}

// generateMetadata generates metadata for the Components.
func (c *Components) generateMetadata() error {
	// Prevent modification if nil
	if c == nil {
		return nil
	}

	// Generate schemas metadata
	for name, schema := range c.Schemas {
		if err := schema.generateMetadata(name+"Schema", false); err != nil {
			return err
		}
	}

	// Generate mapped structured metadata
	if err := c.generateMetadataFromMaps(); err != nil {
		return err
	}

	// Generate reply operations metadata
	for name, reply := range c.Replies {
		if err := reply.generateMetadata(name + "Reply"); err != nil {
			return err
		}
	}

	// Generate reply addresses metadata
	for name, repAddr := range c.ReplyAddresses {
		repAddr.generateMetadata(name + "ReplyAddress")
	}

	return nil
}

// setDependencies sets dependencies between the different elements of the Components.
func (c *Components) setDependencies(spec Specification) error {
	// Prevent modification if nil
	if c == nil {
		return nil
	}

	// Set schemas dependencies
	for _, schema := range c.Schemas {
		if err := schema.setDependencies(spec); err != nil {
			return err
		}
	}

	// Set mapped structured dependencies
	if err := c.setDependenciesFromMaps(spec); err != nil {
		return err
	}

	// Set reply operations dependencies
	for _, reply := range c.Replies {
		if err := reply.setDependencies(&Operation{}, spec); err != nil {
			return err
		}
	}

	// Set reply addresses dependencies
	for _, repAddr := range c.ReplyAddresses {
		if err := repAddr.setDependencies(&Operation{}, spec); err != nil {
			return err
		}
	}

	return nil
}

//nolint:cyclop,funlen
func (c *Components) generateMetadataFromMaps() error {
	for name, entity := range c.Servers {
		entity.generateMetadata(template.Namify(name) + "Server")
	}

	for name, entity := range c.Channels {
		if err := entity.generateMetadata(template.Namify(name) + "Channel"); err != nil {
			return err
		}
	}

	for name, entity := range c.Operations {
		if err := entity.generateMetadata(template.Namify(name) + "Operation"); err != nil {
			return err
		}
	}

	for name, entity := range c.Messages {
		if err := entity.generateMetadata(template.Namify(name) + "Message"); err != nil {
			return err
		}
	}

	for name, entity := range c.SecuritySchemes {
		entity.generateMetadata(template.Namify(name) + "SecurityScheme")
	}

	for name, entity := range c.ServerVariables {
		entity.generateMetadata(template.Namify(name) + "ServerVariable")
	}

	for name, entity := range c.Parameters {
		entity.generateMetadata(template.Namify(name) + "Parameter")
	}

	for name, entity := range c.CorrelationIDs {
		entity.generateMetadata(template.Namify(name) + "CorrelationID")
	}

	for name, entity := range c.ExternalDocs {
		entity.generateMetadata(template.Namify(name) + ExternalDocsNameSuffix)
	}

	for name, entity := range c.Tags {
		entity.generateMetadata(template.Namify(name) + "Tag")
	}

	for name, entity := range c.OperationTraits {
		entity.generateMetadata(template.Namify(name) + "OperationTrait")
	}

	for name, entity := range c.MessageTraits {
		if err := entity.generateMetadata(template.Namify(name) + "MessageTrait"); err != nil {
			return err
		}
	}

	for name, entity := range c.ServerBindings {
		entity.generateMetadata(template.Namify(name) + "ServerBinding")
	}

	for name, entity := range c.ChannelBindings {
		entity.generateMetadata(template.Namify(name) + "ChannelBinding")
	}

	for name, entity := range c.OperationBindings {
		entity.generateMetadata(template.Namify(name) + "OperationBinding")
	}

	for name, entity := range c.MessageBindings {
		entity.generateMetadata(template.Namify(name) + "MessageBinding")
	}

	return nil
}

//nolint:cyclop,funlen
func (c *Components) setDependenciesFromMaps(spec Specification) error {
	for _, entity := range c.Servers {
		if err := entity.setDependencies(spec); err != nil {
			return err
		}
	}

	for _, entity := range c.Channels {
		if err := entity.setDependencies(spec); err != nil {
			return err
		}
	}

	for _, entity := range c.Operations {
		if err := entity.setDependencies(spec); err != nil {
			return err
		}
	}

	for _, entity := range c.Messages {
		if err := entity.setDependencies(spec); err != nil {
			return err
		}
	}

	for _, entity := range c.SecuritySchemes {
		if err := entity.setDependencies(spec); err != nil {
			return err
		}
	}

	for _, entity := range c.ServerVariables {
		if err := entity.setDependencies(spec); err != nil {
			return err
		}
	}

	for _, entity := range c.Parameters {
		if err := entity.setDependencies(spec); err != nil {
			return err
		}
	}

	for _, entity := range c.CorrelationIDs {
		if err := entity.setDependencies(spec); err != nil {
			return err
		}
	}

	for _, entity := range c.ExternalDocs {
		if err := entity.setDependencies(spec); err != nil {
			return err
		}
	}

	for _, entity := range c.Tags {
		if err := entity.setDependencies(spec); err != nil {
			return err
		}
	}

	for _, entity := range c.OperationTraits {
		if err := entity.setDependencies(spec); err != nil {
			return err
		}
	}

	for _, entity := range c.MessageTraits {
		if err := entity.setDependencies(spec); err != nil {
			return err
		}
	}

	for _, entity := range c.ServerBindings {
		if err := entity.setDependencies(spec); err != nil {
			return err
		}
	}

	for _, entity := range c.ChannelBindings {
		if err := entity.setDependencies(spec); err != nil {
			return err
		}
	}

	for _, entity := range c.OperationBindings {
		if err := entity.setDependencies(spec); err != nil {
			return err
		}
	}

	for _, entity := range c.MessageBindings {
		if err := entity.setDependencies(spec); err != nil {
			return err
		}
	}

	return nil
}
