package asyncapiv3

// ServerBindings is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#serverBindingsObject
type ServerBindings struct {
	// --- AsyncAPI fields -----------------------------------------------------

	HTTP         HTTPBinding         `json:"http"`
	WS           WsBinding           `json:"ws"`
	Kafka        KafkaBinding        `json:"kafka"`
	AnyPointMQ   AnyPointMqBinding   `json:"anypointmq"`
	AMQP         AMQPBinding         `json:"amqp"`
	AMQP1        AMQP1Binding        `json:"amqp1"`
	MQTT         MQTTBinding         `json:"mqtt"`
	MQTT5        MQTT5Binding        `json:"mqtt5"`
	NATS         NATSBinding         `json:"nats"`
	JMS          JMSBinding          `json:"jms"`
	SNS          SNSBinding          `json:"sns"`
	Solace       SolaceBinding       `json:"solace"`
	SQS          SQSBinding          `json:"sqs"`
	Stomp        StompBinding        `json:"stomp"`
	Redis        RedisBinding        `json:"redis"`
	Mercure      MercureBinding      `json:"mercure"`
	IBMMQ        IBMMQBinding        `json:"ibmmq"`
	GooglePubSub GooglePubSubBinding `json:"googlepubsub"`
	Pulsar       PulsarBinding       `json:"pulsar"`
	Reference    string              `json:"$ref"`

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
