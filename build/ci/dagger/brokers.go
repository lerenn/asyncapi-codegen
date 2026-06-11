package main

import (
	"asyncapi-codegen/ci/dagger/internal/dagger"
	"fmt"

	testutil "github.com/lerenn/asyncapi-codegen/pkg/utils/test"
)

const (
	// kafkaImage is the image used for kafka.
	// NOTE: the official Apache Kafka image is used (the Bitnami catalog was
	// deprecated). It maps KAFKA_<PROPERTY> environment variables directly to
	// server.properties.
	kafkaImage = "apache/kafka:3.9.0"
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

// kafkaSingleNodeEnv applies the environment variables shared by every Kafka
// broker: a single combined controller/broker node with replication factors set
// to 1 so the internal topics can be created on a one-node cluster.
func kafkaSingleNodeEnv(c *dagger.Container, advertisedHost, protocolMap string) *dagger.Container {
	return c.
		WithEnvVariable("KAFKA_NODE_ID", "0").
		WithEnvVariable("KAFKA_PROCESS_ROLES", "controller,broker").
		WithEnvVariable("KAFKA_LISTENERS", "INTERNAL://:9092,CONTROLLER://:9093").
		WithEnvVariable("KAFKA_ADVERTISED_LISTENERS", "INTERNAL://"+advertisedHost+":9092").
		WithEnvVariable("KAFKA_LISTENER_SECURITY_PROTOCOL_MAP", protocolMap).
		WithEnvVariable("KAFKA_CONTROLLER_QUORUM_VOTERS", "0@localhost:9093").
		WithEnvVariable("KAFKA_CONTROLLER_LISTENER_NAMES", "CONTROLLER").
		WithEnvVariable("KAFKA_INTER_BROKER_LISTENER_NAME", "INTERNAL").
		WithEnvVariable("KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR", "1").
		WithEnvVariable("KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR", "1").
		WithEnvVariable("KAFKA_TRANSACTION_STATE_LOG_MIN_ISR", "1").
		WithEnvVariable("KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS", "0")
}

// kafkaTLSDirectory builds the PEM material expected by the Apache Kafka image:
// a combined keystore (certificate chain + private key in a single file) and a
// CA truststore so the broker trusts its own certificate over SSL/SASL_SSL.
func kafkaTLSDirectory(host string) *dagger.Directory {
	key, cert, cacert, err := testutil.GenerateSelfSignedCertificateWithCA(host)
	if err != nil {
		panic(fmt.Errorf("failed to generate self signed certificate: %w", err))
	}

	return dag.Directory().
		WithNewFile("kafka.keystore.combined.pem", string(cert)+string(key)).
		WithNewFile("kafka.truststore.pem", string(cacert))
}

// brokerKafka returns a container for the Kafka broker.
func brokerKafka() *dagger.Container {
	c := dag.Container().From(kafkaImage)
	c = kafkaSingleNodeEnv(c, "kafka", "CONTROLLER:PLAINTEXT,INTERNAL:PLAINTEXT")
	return c.
		WithExposedPort(9092).
		WithExposedPort(9093)
}

// brokerKafkaSecure returns a container for the Kafka broker secured with TLS.
func brokerKafkaSecure() *dagger.Container {
	c := dag.Container().From(kafkaImage)
	c = kafkaSingleNodeEnv(c, "kafka-tls", "CONTROLLER:PLAINTEXT,INTERNAL:SSL")
	return c.
		// PEM keystore + CA truststore, client certificate auth disabled.
		WithEnvVariable("KAFKA_SSL_KEYSTORE_TYPE", "PEM").
		WithEnvVariable("KAFKA_SSL_KEYSTORE_LOCATION", "/etc/kafka/secrets/kafka.keystore.combined.pem").
		WithEnvVariable("KAFKA_SSL_TRUSTSTORE_TYPE", "PEM").
		WithEnvVariable("KAFKA_SSL_TRUSTSTORE_LOCATION", "/etc/kafka/secrets/kafka.truststore.pem").
		WithEnvVariable("KAFKA_SSL_CLIENT_AUTH", "none").
		WithExposedPort(9092).
		WithExposedPort(9093).
		WithDirectory("/etc/kafka/secrets", kafkaTLSDirectory("kafka-tls"))
}

// kafkaSASLProperties is the server.properties used by the SASL_SSL broker. The
// SCRAM credentials must be present in the KRaft metadata log, which is only
// possible at format time via `--add-scram` (see brokerKafkaSecureBasicAuth).
const kafkaSASLProperties = `node.id=0
process.roles=controller,broker
listeners=INTERNAL://:9092,CONTROLLER://:9093
advertised.listeners=INTERNAL://kafka-tls-basic-auth:9092
listener.security.protocol.map=CONTROLLER:PLAINTEXT,INTERNAL:SASL_SSL
controller.quorum.voters=0@localhost:9093
controller.listener.names=CONTROLLER
inter.broker.listener.name=INTERNAL
log.dirs=/tmp/kraft-combined-logs
offsets.topic.replication.factor=1
transaction.state.log.replication.factor=1
transaction.state.log.min.isr=1
group.initial.rebalance.delay.ms=0
sasl.enabled.mechanisms=SCRAM-SHA-512
sasl.mechanism.inter.broker.protocol=SCRAM-SHA-512
listener.name.internal.sasl.enabled.mechanisms=SCRAM-SHA-512
listener.name.internal.scram-sha-512.sasl.jaas.config=org.apache.kafka.common.security.scram.ScramLoginModule required username="user" password="password";
ssl.keystore.type=PEM
ssl.keystore.location=/etc/kafka/secrets/kafka.keystore.combined.pem
ssl.truststore.type=PEM
ssl.truststore.location=/etc/kafka/secrets/kafka.truststore.pem
ssl.client.auth=none
`

// brokerKafkaSecureBasicAuth returns a container for the Kafka broker secured with TLS and basic auth.
func brokerKafkaSecureBasicAuth() *dagger.Container {
	// Format the storage with the SCRAM credentials, then start the broker.
	startScript := "set -e\n" +
		"CID=$(/opt/kafka/bin/kafka-storage.sh random-uuid)\n" +
		"/opt/kafka/bin/kafka-storage.sh format -t \"$CID\" -c /etc/kafka/secrets/server.properties " +
		"--add-scram SCRAM-SHA-512=[name=user,password=password] --ignore-formatted\n" +
		"exec /opt/kafka/bin/kafka-server-start.sh /etc/kafka/secrets/server.properties\n"

	tlsDir := kafkaTLSDirectory("kafka-tls-basic-auth").
		WithNewFile("server.properties", kafkaSASLProperties)

	return dag.Container().
		From(kafkaImage).
		WithDirectory("/etc/kafka/secrets", tlsDir).
		WithExposedPort(9092).
		WithExposedPort(9093).
		WithoutEntrypoint().
		WithExec([]string{"bash", "-c", startScript})
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
