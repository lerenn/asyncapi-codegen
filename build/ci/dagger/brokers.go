package main

import (
	"fmt"

	"asyncapi-codegen/ci/dagger/internal/dagger"
	testutil "github.com/lerenn/asyncapi-codegen/pkg/utils/test"
)

const (
	// kafkaImage is the image used for kafka.
	kafkaImage = "bitnami/kafka:3.5.1"
	// natsImage is the image used for NATS.
	natsImage = "nats:2.10"
	// rabbitmqImage is the image used for RabbitMQ.
	rabbitmqImage = "rabbitmq:4.0.6"
)

func bindBrokers(brokers map[string]*dagger.Service) func(r *dagger.Container) *dagger.Container {
	return func(r *dagger.Container) *dagger.Container {
		for n, b := range brokers {
			r = r.WithServiceBinding(n, b)
		}

		// Set environment variable to indicate that the application is running
		// in a dockerized environment.
		return r.WithEnvVariable("ASYNCAPI_DOCKERIZED", "true")
	}
}

// brokerServices returns a map of containers for each broker as service.
func brokerServices() map[string]*dagger.Service {
	brokers := make(map[string]*dagger.Service)

	// Kafka
	brokers["kafka"] = brokerKafka().AsService()
	brokers["kafka-tls"] = brokerKafkaSecure().AsService()
	brokers["kafka-tls-basic-auth"] = brokerKafkaSecureBasicAuth().AsService()

	// NATS
	brokers["nats"] = brokerNATS().AsService()
	brokers["nats-tls"] = brokerNATSSecure().AsService()
	brokers["nats-tls-basic-auth"] = brokerNATSSecureBasicAuth().AsService()

	// NATS Jetstream
	brokers["nats-jetstream"] = brokerNATSJetstream().AsService()
	brokers["nats-jetstream-tls"] = brokerNATSJetstreamSecure().AsService()
	brokers["nats-jetstream-tls-basic-auth"] = brokerNATSJetstreamSecureBasicAuth().AsService()

	// RabbitMQ
	brokers["rabbitmq"] = brokerRabbitMQ().AsService()

	return brokers
}

// brokerKafka returns a container for the Kafka broker.
func brokerKafka() *dagger.Container {
	return dag.Container().
		//	Set container image
		From(kafkaImage).

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
		WithExposedPort(9093)
}

// brokerKafkaSecure returns a container for the Kafka broker secured with TLS.
func brokerKafkaSecure() *dagger.Container {
	key, cert, cacert, err := testutil.GenerateSelfSignedCertificateWithCA("kafka-tls")
	if err != nil {
		panic(fmt.Errorf("failed to generate self signed certificate: %w", err))
	}

	tlsDir := dag.Directory().
		WithNewFile("kafka.keystore.key", string(key)).
		WithNewFile("kafka.keystore.pem", string(cert)).
		WithNewFile("kafka.truststore.pem", string(cacert))

	return dag.Container().
		//	Set container image
		From(kafkaImage).

		// Add environment variables
		WithEnvVariable("KAFKA_CFG_NODE_ID", "0").
		WithEnvVariable("KAFKA_CFG_PROCESS_ROLES", "controller,broker").
		WithEnvVariable("KAFKA_CFG_LISTENERS", "INTERNAL://:9092,CONTROLLER://:9093").
		WithEnvVariable("KAFKA_CFG_ADVERTISED_LISTENERS", "INTERNAL://kafka-tls:9092").
		WithEnvVariable("KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP",
			"CONTROLLER:PLAINTEXT,INTERNAL:SSL").
		WithEnvVariable("KAFKA_CFG_CONTROLLER_QUORUM_VOTERS", "0@:9093").
		WithEnvVariable("KAFKA_CFG_CONTROLLER_LISTENER_NAMES", "CONTROLLER").
		WithEnvVariable("KAFKA_CFG_INTER_BROKER_LISTENER_NAME", "INTERNAL").

		// Add tls config
		WithEnvVariable("KAFKA_TLS_TYPE", "PEM").
		// disable client cert
		WithEnvVariable("KAFKA_TLS_CLIENT_AUTH", "none").

		// Add exposed ports
		WithExposedPort(9092).
		WithExposedPort(9093).

		// Add server cert and key directory
		WithDirectory("/bitnami/kafka/config/certs/", tlsDir)
}

// brokerKafkaSecureBasicAuth returns a container for the Kafka broker secured with TLS and basic auth.
func brokerKafkaSecureBasicAuth() *dagger.Container {
	key, cert, cacert, err := testutil.GenerateSelfSignedCertificateWithCA("kafka-tls-basic-auth")
	if err != nil {
		panic(fmt.Errorf("failed to generate self signed certificate: %w", err))
	}

	tlsDir := dag.Directory().
		WithNewFile("kafka.keystore.key", string(key)).
		WithNewFile("kafka.keystore.pem", string(cert)).
		WithNewFile("kafka.truststore.pem", string(cacert))

	return dag.Container().
		//	Set container image
		From(kafkaImage).

		// Add environment variables
		WithEnvVariable("KAFKA_CFG_NODE_ID", "0").
		WithEnvVariable("KAFKA_CFG_PROCESS_ROLES", "controller,broker").
		WithEnvVariable("KAFKA_CFG_LISTENERS", "INTERNAL://:9092,CONTROLLER://:9093").
		WithEnvVariable("KAFKA_CFG_ADVERTISED_LISTENERS", "INTERNAL://kafka-tls-basic-auth:9092").
		WithEnvVariable("KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP",
			"CONTROLLER:PLAINTEXT,INTERNAL:SASL_SSL").
		WithEnvVariable("KAFKA_CFG_CONTROLLER_QUORUM_VOTERS", "0@:9093").
		WithEnvVariable("KAFKA_CFG_CONTROLLER_LISTENER_NAMES", "CONTROLLER").
		WithEnvVariable("KAFKA_CFG_INTER_BROKER_LISTENER_NAME", "INTERNAL").
		WithEnvVariable("KAFKA_CFG_SASL_MECHANISM_INTER_BROKER_PROTOCOL", "SCRAM-SHA-512").

		// Add tls config
		WithEnvVariable("KAFKA_TLS_TYPE", "PEM").

		// add basic auth user and pw
		// WithEnvVariable("KAFKA_CLIENT_USERS", "user").
		// WithEnvVariable("KAFKA_CLIENT_PASSWORDS", "password").
		WithEnvVariable("KAFKA_INTER_BROKER_USER", "user").
		WithEnvVariable("KAFKA_INTER_BROKER_PASSWORD", "password").
		// disable client cert
		WithEnvVariable("KAFKA_TLS_CLIENT_AUTH", "none").

		// Add exposed ports
		WithExposedPort(9092).
		WithExposedPort(9093).

		// Add server cert and key directory
		WithDirectory("/bitnami/kafka/config/certs/", tlsDir)
}

// brokerNATS returns a container for the NATS broker.
func brokerNATS() *dagger.Container {
	return dag.Container().
		// Add base image
		From(natsImage).
		// Add exposed ports
		WithExposedPort(4222)
}

// brokerNATSSecure returns a container for the NATS broker secured with TLS.
func brokerNATSSecure() *dagger.Container {
	key, cert, err := testutil.GenerateSelfSignedCertificate("nats-tls")
	if err != nil {
		panic(fmt.Errorf("failed to generate self signed certificate: %w", err))
	}
	tlsDir := dag.Directory().WithNewFile("server-key.pem", string(key)).WithNewFile("server-cert.pem", string(cert))

	return dag.Container().
		// Add base image
		From(natsImage).
		// Add exposed ports
		WithExposedPort(4222).
		// Add server cert and key directory
		WithDirectory("./tls", tlsDir).
		// Start NATS with tls
		WithoutEntrypoint().
		WithExec([]string{"nats-server", "--tls", "--tlscert=/tls/server-cert.pem", "--tlskey=/tls/server-key.pem"})
}

// brokerNATSSecureBasicAuth returns a container for the NATS broker secured with TLS
// and basic auth user: user password: password.
func brokerNATSSecureBasicAuth() *dagger.Container {
	key, cert, err := testutil.GenerateSelfSignedCertificate("nats-tls-basic-auth")
	if err != nil {
		panic(fmt.Errorf("failed to generate self signed certificate: %w", err))
	}
	tlsDir := dag.Directory().WithNewFile("server-key.pem", string(key)).WithNewFile("server-cert.pem", string(cert))

	return dag.Container().
		// Add base image
		From(natsImage).
		// Add exposed ports
		WithExposedPort(4222).
		// Add server cert and key directory
		WithDirectory("./tls", tlsDir).
		// Start NATS with tls and credentials
		WithoutEntrypoint().
		WithExec([]string{
			"nats-server",
			"--tls",
			"--tlscert=/tls/server-cert.pem",
			"--tlskey=/tls/server-key.pem",
			"--user", "user",
			"--pass", "password"})
}

// brokerNATSJetstream returns a container for the NATS broker.
func brokerNATSJetstream() *dagger.Container {
	return dag.Container().
		// Add base image
		From(natsImage).
		// Add exposed ports
		WithExposedPort(4222).
		// Add command
		WithoutEntrypoint().
		WithExec([]string{"nats-server", "-js"})
}

// brokerNATSJetstreamSecure returns a container for the NATS broker secured with TLS.
func brokerNATSJetstreamSecure() *dagger.Container {
	key, cert, err := testutil.GenerateSelfSignedCertificate("nats-jetstream-tls-basic-auth")
	if err != nil {
		panic(fmt.Errorf("failed to generate self signed certificate: %w", err))
	}
	tlsDir := dag.Directory().WithNewFile("server-key.pem", string(key)).WithNewFile("server-cert.pem", string(cert))

	return dag.Container().
		// Add base image
		From(natsImage).
		// Add exposed ports
		WithExposedPort(4222).
		// Add server cert and key directory
		WithDirectory("./tls", tlsDir).
		// Start NATS jetstream with tls
		WithoutEntrypoint().
		WithExec([]string{"nats-server", "-js", "--tls", "--tlscert=/tls/server-cert.pem", "--tlskey=/tls/server-key.pem"})
}

// brokerNATSJetstreamSecureBasicAuth returns a container for the NATS broker secured with TLS
// and basic auth user: user password: password.
func brokerNATSJetstreamSecureBasicAuth() *dagger.Container {
	key, cert, err := testutil.GenerateSelfSignedCertificate("nats-jetstream-tls")
	if err != nil {
		panic(fmt.Errorf("failed to generate self signed certificate: %w", err))
	}
	tlsDir := dag.Directory().WithNewFile("server-key.pem", string(key)).WithNewFile("server-cert.pem", string(cert))

	return dag.Container().
		// Add base image
		From(natsImage).
		// Add exposed ports
		WithExposedPort(4222).
		// Add server cert and key directory
		WithDirectory("./tls", tlsDir).
		// Start NATS jetstream with tls and credentials
		WithoutEntrypoint().
		WithExec([]string{
			"nats-server",
			"-js",
			"--tls",
			"--tlscert=/tls/server-cert.pem",
			"--tlskey=/tls/server-key.pem",
			"--user", "user",
			"--pass", "password",
		})
}

// brokerRabbitMQ returns a container for the RabbitMQ broker.
func brokerRabbitMQ() *dagger.Container {
	return dag.Container().
		// Add base image
		From(rabbitmqImage).
		// Add exposed ports
		WithExposedPort(5672)
}
