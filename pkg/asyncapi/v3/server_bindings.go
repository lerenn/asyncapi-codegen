package asyncapiv3

import "github.com/lerenn/asyncapi-codegen/pkg/utils/template"

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

// Process processes the ServerBindings to make it ready for code generation.
func (ob *ServerBindings) Process(name string, spec Specification) error {
	// Prevent modification if nil
	if ob == nil {
		return nil
	}

	// Set name
	ob.Name = template.Namify(name)

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
