package asyncapiv3

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

// Process processes the Components structure to make it ready for code generation.
func (c *Components) Process(spec Specification) error {
	// Prevent modification if nil
	if c == nil {
		return nil
	}

	// Process schemas
	for name, schema := range c.Schemas {
		if err := schema.Process(name+"Schema", spec, false); err != nil {
			return err
		}
	}

	// Process mapped structured
	if err := c.processMaps(spec); err != nil {
		return err
	}

	// Process reply operations
	for name, reply := range c.Replies {
		if err := reply.Process(name+"Reply", &Operation{}, spec); err != nil {
			return err
		}
	}

	// Process reply addresses
	for name, repAddr := range c.ReplyAddresses {
		if err := repAddr.Process(name+"ReplyAddress", &Operation{}, spec); err != nil {
			return err
		}
	}

	return nil
}

//nolint:cyclop,funlen
func (c *Components) processMaps(spec Specification) error {
	if err := processMap(spec, c.Servers, "Server"); err != nil {
		return err
	}

	if err := processMap(spec, c.Channels, "Channel"); err != nil {
		return err
	}

	if err := processMap(spec, c.Operations, "Operation"); err != nil {
		return err
	}

	if err := processMap(spec, c.Messages, "Message"); err != nil {
		return err
	}

	if err := processMap(spec, c.SecuritySchemes, "SecurityScheme"); err != nil {
		return err
	}

	if err := processMap(spec, c.ServerVariables, "ServerVariable"); err != nil {
		return err
	}

	if err := processMap(spec, c.Parameters, "Parameter"); err != nil {
		return err
	}

	if err := processMap(spec, c.CorrelationIDs, "CorrelationID"); err != nil {
		return err
	}

	if err := processMap(spec, c.ExternalDocs, ExternalDocsNameSuffix); err != nil {
		return err
	}

	if err := processMap(spec, c.Tags, "Tag"); err != nil {
		return err
	}

	if err := processMap(spec, c.OperationTraits, "OperationTrait"); err != nil {
		return err
	}

	if err := processMap(spec, c.MessageTraits, "MessageTrait"); err != nil {
		return err
	}

	if err := processMap(spec, c.ServerBindings, "ServerBinding"); err != nil {
		return err
	}

	if err := processMap(spec, c.ChannelBindings, "ChannelBinding"); err != nil {
		return err
	}

	if err := processMap(spec, c.OperationBindings, "OperationBinding"); err != nil {
		return err
	}

	return processMap(spec, c.MessageBindings, "MessageBinding")
}
