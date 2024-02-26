package asyncapiv3

const (
	// BindingsSuffix is the suffix added to the bindings name.
	BindingsSuffix = "Bindings"
)

// HTTPBinding represents protocol-specific information for an HTTP channel.
type HTTPBinding any

// WsBinding represents protocol-specific information for a WebSockets channel.
type WsBinding any

// KafkaBinding represents protocol-specific information for a Kafka channel.
type KafkaBinding any

// AnyPointMqBinding represents protocol-specific information for an Anypoint MQ channel.
type AnyPointMqBinding any

// AMQPBinding represents protocol-specific information for an AMQP 0-9-1 channel.
type AMQPBinding any

// AMQP1Binding represents protocol-specific information for an AMQP 1.0 channel.
type AMQP1Binding any

// MQTTBinding represents protocol-specific information for an MQTT channel.
type MQTTBinding any

// MQTT5Binding represents protocol-specific information for an MQTT 5 channel.
type MQTT5Binding any

// NATSBinding represents protocol-specific information for a NATS channel.
type NATSBinding any

// JMSBinding represents protocol-specific information for a JMS channel.
type JMSBinding any

// SNSBinding represents protocol-specific information for an SNS channel.
type SNSBinding any

// SolaceBinding represents protocol-specific information for a Solace channel.
type SolaceBinding any

// SQSBinding represents protocol-specific information for an SQS channel.
type SQSBinding any

// StompBinding represents protocol-specific information for a STOMP channel.
type StompBinding any

// RedisBinding represents protocol-specific information for a Redis channel.
type RedisBinding any

// MercureBinding represents protocol-specific information for a Mercure channel.
type MercureBinding any

// IBMMQBinding represents protocol-specific information for an IBM MQ channel.
type IBMMQBinding any

// GooglePubSubBinding represents protocol-specific information for a Google Cloud Pub/Sub channel.
type GooglePubSubBinding any

// PulsarBinding represents protocol-specific information for a Pulsar channel.
type PulsarBinding any
