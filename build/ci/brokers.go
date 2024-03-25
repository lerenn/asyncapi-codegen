package main

import (
	"dagger/asyncapi-codegen-ci/internal/dagger"
)

const (
	// kafkaImage is the image used for kafka.
	kafkaImage = "bitnami/kafka:3.5.1"
	// natsImage is the image used for NATS.
	natsImage = "nats:2.10"
)

func bindBrokers(brokers map[string]*dagger.Service) func(r *dagger.Container) *dagger.Container {
	return func(r *dagger.Container) *dagger.Container {
		for n, b := range brokers {
			r = r.WithServiceBinding(n, b)
		}
		return r
	}
}

// brokers returns a map of containers for each broker as service.
func brokers() map[string]*dagger.Service {
	brokers := make(map[string]*dagger.Service)

	brokers["kafka"] = brokerKafka()
	brokers["nats"] = brokerNATS()
	brokers["nats-jetstream"] = brokerNATSJetstream()

	return brokers
}

// brokerKafka returns a service for the Kafka broker.
func brokerKafka() *dagger.Service {
	return dag.Container().
		From(kafkaImage).
		WithEnvVariable("KAFKA_CFG_NODE_ID", "0").
		WithEnvVariable("KAFKA_CFG_PROCESS_ROLES", "controller,broker").
		WithEnvVariable("KAFKA_CFG_LISTENERS", "INTERNAL://:9092,CONTROLLER://:9093").
		WithEnvVariable("KAFKA_CFG_ADVERTISED_LISTENERS", "INTERNAL://kafka:9092").
		WithEnvVariable("KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP", "CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT,INTERNAL:PLAINTEXT").
		WithEnvVariable("KAFKA_CFG_CONTROLLER_QUORUM_VOTERS", "0@:9093").
		WithEnvVariable("KAFKA_CFG_CONTROLLER_LISTENER_NAMES", "CONTROLLER").
		WithEnvVariable("KAFKA_CFG_INTER_BROKER_LISTENER_NAME", "INTERNAL").
		WithExposedPort(9092).
		WithExposedPort(9093).
		AsService()
}

// brokerNATS returns a service for the NATS broker.
func brokerNATS() *dagger.Service {
	return dag.Container().
		From(natsImage).
		WithExposedPort(4222).
		AsService()
}

// brokerNATSJetstream returns a service for the NATS broker.
func brokerNATSJetstream() *dagger.Service {
	return dag.Container().
		From(natsImage).
		WithExposedPort(4222).
		WithExec([]string{"-js"}).
		AsService()
}
