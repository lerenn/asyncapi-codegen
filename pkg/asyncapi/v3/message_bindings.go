package asyncapiv3

import "github.com/lerenn/asyncapi-codegen/pkg/utils/template"

// MessageBindings is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#messageBindingsObject
type MessageBindings struct {
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

	Name        string           `json:"-"`
	ReferenceTo *MessageBindings `json:"-"`
}

// generateMetadata generates metadata for the MessageBindings.
func (mb *MessageBindings) generateMetadata(name string) {
	// Prevent modification if nil
	if mb == nil {
		return
	}

	// Set name
	mb.Name = template.Namify(name)
}

// setDependencies sets dependencies between the different elements of the MessageBindings.
func (mb *MessageBindings) setDependencies(spec Specification) error {
	// Prevent modification if nil
	if mb == nil {
		return nil
	}

	// Add pointer to reference if there is one
	if mb.Reference != "" {
		refTo, err := spec.ReferenceMessageBindings(mb.Reference)
		if err != nil {
			return err
		}
		mb.ReferenceTo = refTo
	}

	return nil
}
