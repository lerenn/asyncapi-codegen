package pipeline

import "dagger.io/dagger"

// BindBrokers is used as a helper to bind brokers to a container.
func BindBrokers(brokers map[string]*dagger.Service) func(r *dagger.Container) *dagger.Container {
	return func(r *dagger.Container) *dagger.Container {
		for n, b := range brokers {
			r = r.WithServiceBinding(n, b)
		}
		return r
	}
}

// Brokers returns a map of containers for each broker as service.
func Brokers(client *dagger.Client) map[string]*dagger.Service {
	brokers := make(map[string]*dagger.Service)

	brokers["kafka"] = BrokerKafka(client)
	brokers["nats"] = BrokerNATS(client)

	return brokers
}

// BrokerKafka returns a service for the Kafka broker.
func BrokerKafka(client *dagger.Client) *dagger.Service {
	return client.Container().
		//	Set container image
		From(KafkaImage).

		// Add environment variables
		WithEnvVariable("KAFKA_CFG_NODE_ID", "0").
		WithEnvVariable("KAFKA_CFG_PROCESS_ROLES", "controller,broker").
		WithEnvVariable("KAFKA_CFG_LISTENERS", "INTERNAL://:9092,CONTROLLER://:9093").
		WithEnvVariable("KAFKA_CFG_ADVERTISED_LISTENERS", "INTERNAL://kafka:9092").
		WithEnvVariable("KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP",
			"CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT,INTERNAL:PLAINTEXT").
		WithEnvVariable("KAFKA_CFG_CONTROLLER_QUORUM_VOTERS", "0@:9093").
		WithEnvVariable("KAFKA_CFG_CONTROLLER_LISTENER_NAMES", "CONTROLLER").
		WithEnvVariable("KAFKA_CFG_INTER_BROKER_LISTENER_NAME", "INTERNAL").

		// Add exposed ports
		WithExposedPort(9092).
		WithExposedPort(9093).

		// Return container as a service
		AsService()
}

// BrokerNATS returns a service for the NATS broker.
func BrokerNATS(client *dagger.Client) *dagger.Service {
	return client.Container().
		// Add base image
		From(NATSImage).
		// Add exposed ports
		WithExposedPort(4222).
		// Return container as a service
		AsService()
}
