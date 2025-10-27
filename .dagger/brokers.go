package main

import (
	"asyncapi-codegen/dagger/internal/dagger"
	"fmt"
)

const (
	// kafkaImage is the image used for kafka.
	kafkaImage = "confluentinc/cp-kafka:7.5.0"
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

	return brokers
}

// brokerKafka returns a container for the Kafka broker.
func brokerKafka() *dagger.Container {
	return dag.Container().
		// Set container image
		From(kafkaImage).
		// Add environment variables for Confluent Platform
		WithEnvVariable("CLUSTER_ID", "AAAAAAAAAAAAAAAAAAAAAA").
		WithEnvVariable("KAFKA_NODE_ID", "0").
		WithEnvVariable("KAFKA_PROCESS_ROLES", "broker,controller").
		WithEnvVariable("KAFKA_CONTROLLER_QUORUM_VOTERS", "0@localhost:9093").
		WithEnvVariable("KAFKA_LISTENERS", "PLAINTEXT://localhost:29092,CONTROLLER://localhost:9093,PLAINTEXT_HOST://0.0.0.0:9092").
		WithEnvVariable("KAFKA_ADVERTISED_LISTENERS", "PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092").
		WithEnvVariable("KAFKA_LISTENER_SECURITY_PROTOCOL_MAP", "CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT").
		WithEnvVariable("KAFKA_CONTROLLER_LISTENER_NAMES", "CONTROLLER").
		WithEnvVariable("KAFKA_INTER_BROKER_LISTENER_NAME", "PLAINTEXT").
		WithEnvVariable("KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR", "1").
		WithEnvVariable("KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR", "1").
		WithEnvVariable("KAFKA_TRANSACTION_STATE_LOG_MIN_ISR", "1").
		// Add exposed ports
		WithExposedPort(9092).
		WithExposedPort(9093).
		WithExposedPort(29092)
}

// brokerKafkaSecure returns a container for the Kafka broker secured with TLS.
func brokerKafkaSecure() *dagger.Container {
	key, cert, cacert, err := GenerateSelfSignedCertificateWithCA("kafka-tls")
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
		// Add environment variables for Confluent Platform with TLS
		WithEnvVariable("CLUSTER_ID", "AAAAAAAAAAAAAAAAAAAAAA").
		WithEnvVariable("KAFKA_NODE_ID", "0").
		WithEnvVariable("KAFKA_PROCESS_ROLES", "broker,controller").
		WithEnvVariable("KAFKA_CONTROLLER_QUORUM_VOTERS", "0@localhost:9093").
		WithEnvVariable("KAFKA_LISTENERS", "SSL://localhost:29092,CONTROLLER://localhost:9093,SSL_HOST://0.0.0.0:9092").
		WithEnvVariable("KAFKA_ADVERTISED_LISTENERS", "SSL://kafka-tls:29092,SSL_HOST://localhost:9092").
		WithEnvVariable("KAFKA_LISTENER_SECURITY_PROTOCOL_MAP", "CONTROLLER:PLAINTEXT,SSL:SSL,SSL_HOST:SSL").
		WithEnvVariable("KAFKA_CONTROLLER_LISTENER_NAMES", "CONTROLLER").
		WithEnvVariable("KAFKA_INTER_BROKER_LISTENER_NAME", "SSL").
		WithEnvVariable("KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR", "1").
		WithEnvVariable("KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR", "1").
		WithEnvVariable("KAFKA_TRANSACTION_STATE_LOG_MIN_ISR", "1").
		WithEnvVariable("KAFKA_SSL_KEYSTORE_FILENAME", "kafka.keystore.pem").
		WithEnvVariable("KAFKA_SSL_KEYSTORE_LOCATION", "/etc/kafka/secrets/kafka.keystore.pem").
		WithEnvVariable("KAFKA_SSL_TRUSTSTORE_FILENAME", "kafka.truststore.pem").
		WithEnvVariable("KAFKA_SSL_TRUSTSTORE_LOCATION", "/etc/kafka/secrets/kafka.truststore.pem").
		WithEnvVariable("KAFKA_SSL_CLIENT_AUTH", "none").
		WithEnvVariable("KAFKA_SSL_ENDPOINT_IDENTIFICATION_ALGORITHM", "").
		// Add exposed ports
		WithExposedPort(9092).
		WithExposedPort(9093).
		WithExposedPort(29092).
		// Mount certificates
		WithDirectory("/etc/kafka/secrets", tlsDir)
}

// brokerKafkaSecureBasicAuth returns a container for the Kafka broker secured with TLS and basic auth.
func brokerKafkaSecureBasicAuth() *dagger.Container {
	key, cert, cacert, err := GenerateSelfSignedCertificateWithCA("kafka-tls-basic-auth")
	if err != nil {
		panic(fmt.Errorf("failed to generate self signed certificate: %w", err))
	}

	tlsDir := dag.Directory().
		WithNewFile("kafka.keystore.key", string(key)).
		WithNewFile("kafka.keystore.pem", string(cert)).
		WithNewFile("kafka.truststore.pem", string(cacert))

	// Create JAAS config for SASL
	jaasConfig := `
KafkaServer {
  org.apache.kafka.common.security.plain.PlainLoginModule required
  username="admin"
  password="admin-secret"
  user_admin="admin-secret"
  user_user="password";
};
`

	return dag.Container().
		//	Set container image
		From(kafkaImage).
		// Add environment variables for Confluent Platform with TLS and SASL
		WithEnvVariable("CLUSTER_ID", "AAAAAAAAAAAAAAAAAAAAAA").
		WithEnvVariable("KAFKA_NODE_ID", "0").
		WithEnvVariable("KAFKA_PROCESS_ROLES", "broker,controller").
		WithEnvVariable("KAFKA_CONTROLLER_QUORUM_VOTERS", "0@localhost:9093").
		WithEnvVariable("KAFKA_LISTENERS", "SASL_SSL://localhost:29092,CONTROLLER://localhost:9093,SASL_SSL_HOST://0.0.0.0:9092").
		WithEnvVariable("KAFKA_ADVERTISED_LISTENERS", "SASL_SSL://kafka-tls-basic-auth:29092,SASL_SSL_HOST://localhost:9092").
		WithEnvVariable("KAFKA_LISTENER_SECURITY_PROTOCOL_MAP", "CONTROLLER:PLAINTEXT,SASL_SSL:SASL_SSL,SASL_SSL_HOST:SASL_SSL").
		WithEnvVariable("KAFKA_CONTROLLER_LISTENER_NAMES", "CONTROLLER").
		WithEnvVariable("KAFKA_INTER_BROKER_LISTENER_NAME", "SASL_SSL").
		WithEnvVariable("KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR", "1").
		WithEnvVariable("KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR", "1").
		WithEnvVariable("KAFKA_TRANSACTION_STATE_LOG_MIN_ISR", "1").
		WithEnvVariable("KAFKA_SSL_KEYSTORE_FILENAME", "kafka.keystore.pem").
		WithEnvVariable("KAFKA_SSL_KEYSTORE_LOCATION", "/etc/kafka/secrets/kafka.keystore.pem").
		WithEnvVariable("KAFKA_SSL_TRUSTSTORE_FILENAME", "kafka.truststore.pem").
		WithEnvVariable("KAFKA_SSL_TRUSTSTORE_LOCATION", "/etc/kafka/secrets/kafka.truststore.pem").
		WithEnvVariable("KAFKA_SSL_CLIENT_AUTH", "none").
		WithEnvVariable("KAFKA_SSL_ENDPOINT_IDENTIFICATION_ALGORITHM", "").
		WithEnvVariable("KAFKA_SASL_ENABLED_MECHANISMS", "PLAIN").
		WithEnvVariable("KAFKA_SASL_MECHANISM_INTER_BROKER_PROTOCOL", "PLAIN").
		WithEnvVariable("KAFKA_OPTS", "-Djava.security.auth.login.config=/etc/kafka/kafka_server_jaas.conf").
		// Add exposed ports
		WithExposedPort(9092).
		WithExposedPort(9093).
		WithExposedPort(29092).
		// Mount certificates and JAAS config
		WithDirectory("/etc/kafka/secrets", tlsDir).
		WithNewFile("/etc/kafka/kafka_server_jaas.conf", jaasConfig)
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
	key, cert, err := GenerateSelfSignedCertificate("nats-tls")
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
		WithDefaultArgs([]string{"nats-server", "--tls", "--tlscert=/tls/server-cert.pem", "--tlskey=/tls/server-key.pem"})
}

// brokerNATSSecureBasicAuth returns a container for the NATS broker secured with TLS
// and basic auth user: user password: password.
func brokerNATSSecureBasicAuth() *dagger.Container {
	key, cert, err := GenerateSelfSignedCertificate("nats-tls-basic-auth")
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
		WithDefaultArgs([]string{
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
		WithDefaultArgs([]string{"nats-server", "-js"})
}

// brokerNATSJetstreamSecure returns a container for the NATS broker secured with TLS.
func brokerNATSJetstreamSecure() *dagger.Container {
	key, cert, err := GenerateSelfSignedCertificate("nats-jetstream-tls-basic-auth")
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
		WithDefaultArgs([]string{"nats-server", "-js", "--tls", "--tlscert=/tls/server-cert.pem", "--tlskey=/tls/server-key.pem"})
}

// brokerNATSJetstreamSecureBasicAuth returns a container for the NATS broker secured with TLS
// and basic auth user: user password: password.
func brokerNATSJetstreamSecureBasicAuth() *dagger.Container {
	key, cert, err := GenerateSelfSignedCertificate("nats-jetstream-tls")
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
		WithDefaultArgs([]string{
			"nats-server",
			"-js",
			"--tls",
			"--tlscert=/tls/server-cert.pem",
			"--tlskey=/tls/server-key.pem",
			"--user", "user",
			"--pass", "password",
		})
}
