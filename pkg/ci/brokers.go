package ci

import (
	"fmt"

	"dagger.io/dagger"

	testutil "github.com/lerenn/asyncapi-codegen/pkg/utils/test"
)

// BindBrokers is used as a helper to bind brokers to a container.
func BindBrokers(brokers map[string]*dagger.Service) func(r *dagger.Container) *dagger.Container {
	return func(r *dagger.Container) *dagger.Container {
		// Bind all brokers to the container.
		for n, b := range brokers {
			r = r.WithServiceBinding(n, b)
		}

		// Set environment variable to indicate that the application is running
		// in a dockerized environment.
		return r.WithEnvVariable("ASYNCAPI_DOCKERIZED", "true")
	}
}

// Brokers returns a map of containers for each broker as service.
func Brokers(client *dagger.Client) map[string]*dagger.Service {
	brokers := make(map[string]*dagger.Service)

	// Kafka
	brokers["kafka"] = BrokerKafka(client)
	brokers["kafka-tls"] = BrokerKafkaSecure(client)
	brokers["kafka-tls-basic-auth"] = BrokerKafkaSecureBasicAuth(client)

	// NATS
	brokers["nats"] = BrokerNATS(client)
	brokers["nats-tls"] = BrokerNATSSecure(client)
	brokers["nats-tls-basic-auth"] = BrokerNATSSecureBasicAuth(client)

	// NATS Jetstream
	brokers["nats-jetstream"] = BrokerNATSJetstream(client)
	brokers["nats-jetstream-tls"] = BrokerNATSJetstreamSecure(client)
	brokers["nats-jetstream-tls-basic-auth"] = BrokerNATSJetstreamSecureBasicAuth(client)

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

// BrokerKafkaSecure returns a service for the Kafka broker secured with TLS.
func BrokerKafkaSecure(client *dagger.Client) *dagger.Service {
	key, cert, cacert, err := testutil.GenerateSelfSignedCertificateWithCA("kafka-tls")
	if err != nil {
		panic(fmt.Errorf("failed to generate self signed certificate: %w", err))
	}

	tlsDir := client.Directory().
		WithNewFile("kafka.keystore.key", string(key)).
		WithNewFile("kafka.keystore.pem", string(cert)).
		WithNewFile("kafka.truststore.pem", string(cacert))

	return client.Container().
		//	Set container image
		From(KafkaImage).

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

// BrokerKafkaSecureBasicAuth returns a service for the Kafka broker secured with TLS and basic auth.
func BrokerKafkaSecureBasicAuth(client *dagger.Client) *dagger.Service {
	key, cert, cacert, err := testutil.GenerateSelfSignedCertificateWithCA("kafka-tls-basic-auth")
	if err != nil {
		panic(fmt.Errorf("failed to generate self signed certificate: %w", err))
	}

	tlsDir := client.Directory().
		WithNewFile("kafka.keystore.key", string(key)).
		WithNewFile("kafka.keystore.pem", string(cert)).
		WithNewFile("kafka.truststore.pem", string(cacert))

	return client.Container().
		//	Set container image
		From(KafkaImage).

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

// BrokerNATSSecure returns a service for the NATS broker secured with TLS.
func BrokerNATSSecure(client *dagger.Client) *dagger.Service {
	key, cert, err := testutil.GenerateSelfSignedCertificate("nats-tls")
	if err != nil {
		panic(fmt.Errorf("failed to generate self signed certificate: %w", err))
	}
	tlsDir := client.Directory().WithNewFile("server-key.pem", string(key)).WithNewFile("server-cert.pem", string(cert))

	return client.Container().
		// Add base image
		From(NATSImage).
		// Add exposed ports
		WithExposedPort(4222).
		// Add server cert and key directory
		WithDirectory("./tls", tlsDir).
		// Start NATS with tls
		WithExec([]string{"--tls", "--tlscert=/tls/server-cert.pem", "--tlskey=/tls/server-key.pem"}).
		// Return container as a service
		AsService()
}

// BrokerNATSSecureBasicAuth returns a service for the NATS broker secured with TLS
// and basic auth user: user password: password.
func BrokerNATSSecureBasicAuth(client *dagger.Client) *dagger.Service {
	key, cert, err := testutil.GenerateSelfSignedCertificate("nats-tls-basic-auth")
	if err != nil {
		panic(fmt.Errorf("failed to generate self signed certificate: %w", err))
	}
	tlsDir := client.Directory().WithNewFile("server-key.pem", string(key)).WithNewFile("server-cert.pem", string(cert))

	return client.Container().
		// Add base image
		From(NATSImage).
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

// BrokerNATSJetstream returns a service for the NATS broker.
func BrokerNATSJetstream(client *dagger.Client) *dagger.Service {
	return client.Container().
		// Add base image
		From(NATSImage).
		// Add exposed ports
		WithExposedPort(4222).
		// Add command
		WithExec([]string{"-js"}).
		// Return container as a service
		AsService()
}

// BrokerNATSJetstreamSecure returns a service for the NATS broker secured with TLS.
func BrokerNATSJetstreamSecure(client *dagger.Client) *dagger.Service {
	key, cert, err := testutil.GenerateSelfSignedCertificate("nats-jetstream-tls-basic-auth")
	if err != nil {
		panic(fmt.Errorf("failed to generate self signed certificate: %w", err))
	}
	tlsDir := client.Directory().WithNewFile("server-key.pem", string(key)).WithNewFile("server-cert.pem", string(cert))

	return client.Container().
		// Add base image
		From(NATSImage).
		// Add exposed ports
		WithExposedPort(4222).
		// Add server cert and key directory
		WithDirectory("./tls", tlsDir).
		// Start NATS jetstream with tls
		WithExec([]string{"-js", "--tls", "--tlscert=/tls/server-cert.pem", "--tlskey=/tls/server-key.pem"}).
		// Return container as a service
		AsService()
}

// BrokerNATSJetstreamSecureBasicAuth returns a service for the NATS broker secured with TLS
// and basic auth user: user password: password.
func BrokerNATSJetstreamSecureBasicAuth(client *dagger.Client) *dagger.Service {
	key, cert, err := testutil.GenerateSelfSignedCertificate("nats-jetstream-tls")
	if err != nil {
		panic(fmt.Errorf("failed to generate self signed certificate: %w", err))
	}
	tlsDir := client.Directory().WithNewFile("server-key.pem", string(key)).WithNewFile("server-cert.pem", string(cert))

	return client.Container().
		// Add base image
		From(NATSImage).
		// Add exposed ports
		WithExposedPort(4222).
		// Add server cert and key directory
		WithDirectory("./tls", tlsDir).
		// Start NATS jetstream with tls and credentials
		WithExec([]string{"-js", "--tls", "--tlscert=/tls/server-cert.pem", "--tlskey=/tls/server-key.pem", "--user", "user", "--pass", "password"}). //nolint:lll
		// Return container as a service
		AsService()
}
