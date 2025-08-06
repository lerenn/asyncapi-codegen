package asyncapiv3

// ServerBindings is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#serverBindingsObject
type ServerBindings struct {
	// --- AsyncAPI fields -----------------------------------------------------

	HTTP         HTTPBinding         `json:"http,omitzero"`
	WS           WsBinding           `json:"ws,omitzero"`
	Kafka        KafkaBinding        `json:"kafka,omitzero"`
	AnyPointMQ   AnyPointMqBinding   `json:"anypointmq,omitzero"`
	AMQP         AMQPBinding         `json:"amqp,omitzero"`
	AMQP1        AMQP1Binding        `json:"amqp1,omitzero"`
	MQTT         MQTTBinding         `json:"mqtt,omitzero"`
	MQTT5        MQTT5Binding        `json:"mqtt5,omitzero"`
	NATS         NATSBinding         `json:"nats,omitzero"`
	JMS          JMSBinding          `json:"jms,omitzero"`
	SNS          SNSBinding          `json:"sns,omitzero"`
	Solace       SolaceBinding       `json:"solace,omitzero"`
	SQS          SQSBinding          `json:"sqs,omitzero"`
	Stomp        StompBinding        `json:"stomp,omitzero"`
	Redis        RedisBinding        `json:"redis,omitzero"`
	Mercure      MercureBinding      `json:"mercure,omitzero"`
	IBMMQ        IBMMQBinding        `json:"ibmmq,omitzero"`
	GooglePubSub GooglePubSubBinding `json:"googlepubsub,omitzero"`
	Pulsar       PulsarBinding       `json:"pulsar,omitzero"`
	Reference    string              `json:"$ref,omitempty"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name        string          `json:"-"`
	ReferenceTo *ServerBindings `json:"-"`
}

// generateMetadata generates metadata for the ServerBindings.
func (ob *ServerBindings) generateMetadata(parentName, name string) {
	// Prevent modification if nil
	if ob == nil {
		return
	}

	// Set name
	ob.Name = generateFullName(parentName, name, BindingsSuffix, nil)
}

// setDependencies sets dependencies between the different elements of the ServerBindings.
func (ob *ServerBindings) setDependencies(spec Specification) error {
	// Prevent modification if nil
	if ob == nil {
		return nil
	}

	// Add pointer to reference if there is one
	if ob.Reference != "" {
		refTo, err := spec.ReferenceServerBindings(ob.Reference)
		if err != nil {
			return err
		}
		ob.ReferenceTo = refTo
	}

	return nil
}
