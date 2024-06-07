package main

import (
	"dagger/asyncapi-codegen-ci/internal/dagger"
	"fmt"

	testutil "github.com/TheSadlig/asyncapi-codegen/pkg/utils/test"
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

		// Set environment variable to indicate that the application is running
		// in a dockerized environment.
		return r.WithEnvVariable("ASYNCAPI_DOCKERIZED", "true")
	}
}

// brokers returns a map of containers for each broker as service.
func brokers() map[string]*dagger.Service {
	brokers := make(map[string]*dagger.Service)

	// Kafka
	brokers["kafka"] = brokerKafka()
	brokers["kafka-tls"] = brokerKafkaSecure()
	brokers["kafka-tls-basic-auth"] = brokerKafkaSecureBasicAuth()

	// NATS
	brokers["nats"] = brokerNATS()
	brokers["nats-tls"] = brokerNATSSecure()
	brokers["nats-tls-basic-auth"] = brokerNATSSecureBasicAuth()

	// NATS Jetstream
	brokers["nats-jetstream"] = brokerNATSJetstream()
	brokers["nats-jetstream-tls"] = brokerNATSJetstreamSecure()
	brokers["nats-jetstream-tls-basic-auth"] = brokerNATSJetstreamSecureBasicAuth()

	return brokers
}

// brokerKafka returns a service for the Kafka broker.
func brokerKafka() *dagger.Service {
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
		WithExposedPort(9093).

		// Return container as a service
		AsService()
}

// brokerKafkaSecure returns a service for the Kafka broker secured with TLS.
func brokerKafkaSecure() *dagger.Service {
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
		WithDirectory("/bitnami/kafka/config/certs/", tlsDir).

		// Return container as a service
		AsService()
}

// brokerKafkaSecureBasicAuth returns a service for the Kafka broker secured with TLS and basic auth.
func brokerKafkaSecureBasicAuth() *dagger.Service {
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
		WithDirectory("/bitnami/kafka/config/certs/", tlsDir).

		// Return container as a service
		AsService()
}

// brokerNATS returns a service for the NATS broker.
func brokerNATS() *dagger.Service {
	return dag.Container().
		// Add base image
		From(natsImage).
		// Add exposed ports
		WithExposedPort(4222).
		// Return container as a service
		AsService()
}

// brokerNATSSecure returns a service for the NATS broker secured with TLS.
func brokerNATSSecure() *dagger.Service {
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
		WithExec([]string{"--tls", "--tlscert=/tls/server-cert.pem", "--tlskey=/tls/server-key.pem"}).
		// Return container as a service
		AsService()
}

// brokerNATSSecureBasicAuth returns a service for the NATS broker secured with TLS
// and basic auth user: user password: password.
func brokerNATSSecureBasicAuth() *dagger.Service {
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
		WithExec([]string{
			"--tls",
			"--tlscert=/tls/server-cert.pem",
			"--tlskey=/tls/server-key.pem",
			"--user", "user",
			"--pass", "password"}).
		// Return container as a service
		AsService()
}

// brokerNATSJetstream returns a service for the NATS broker.
func brokerNATSJetstream() *dagger.Service {
	return dag.Container().
		// Add base image
		From(natsImage).
		// Add exposed ports
		WithExposedPort(4222).
		// Add command
		WithExec([]string{"-js"}).
		// Return container as a service
		AsService()
}

// brokerNATSJetstreamSecure returns a service for the NATS broker secured with TLS.
func brokerNATSJetstreamSecure() *dagger.Service {
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
		WithExec([]string{"-js", "--tls", "--tlscert=/tls/server-cert.pem", "--tlskey=/tls/server-key.pem"}).
		// Return container as a service
		AsService()
}

// brokerNATSJetstreamSecureBasicAuth returns a service for the NATS broker secured with TLS
// and basic auth user: user password: password.
func brokerNATSJetstreamSecureBasicAuth() *dagger.Service {
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
		WithExec([]string{"-js", "--tls", "--tlscert=/tls/server-cert.pem", "--tlskey=/tls/server-key.pem", "--user", "user", "--pass", "password"}). //nolint:lll
		// Return container as a service
		AsService()
}
