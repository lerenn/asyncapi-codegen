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
func (c *Components) Process(spec Specification) {
	// Prevent modification if nil
	if c == nil {
		return
	}

	// Process schemas
	for name, schema := range c.Schemas {
		schema.Process(name, spec, false)
	}

	// Process mapped structured
	processMap(spec, c.Servers)
	processMap(spec, c.Channels)
	processMap(spec, c.Operations)
	processMap(spec, c.Messages)
	processMap(spec, c.SecuritySchemes)
	processMap(spec, c.ServerVariables)
	processMap(spec, c.Parameters)
	processMap(spec, c.CorrelationIDs)
	processMap(spec, c.Replies)
	processMap(spec, c.ReplyAddresses)
	processMap(spec, c.ExternalDocs)
	processMap(spec, c.Tags)
	processMap(spec, c.OperationTraits)
	processMap(spec, c.MessageTraits)
	processMap(spec, c.ServerBindings)
	processMap(spec, c.ChannelBindings)
	processMap(spec, c.OperationBindings)
	processMap(spec, c.MessageBindings)
}
