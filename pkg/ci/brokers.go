package ci

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	"dagger.io/dagger"
)

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

	brokers["kafka"] = BrokerRedPandaAsKafka(client)
	brokers["kafka-tls"] = BrokerKafkaSecure(client)
	brokers["kafka-tls-basic-auth"] = BrokerKafkaSecureBasicAuth(client)
	brokers["nats"] = BrokerNATS(client)
	brokers["nats-tls"] = BrokerNATSSecure(client)
	brokers["nats-tls-basic-auth"] = BrokerNATSSecureBasicAuth(client)
	brokers["nats-jetstream"] = BrokerNATSJetstream(client)
	brokers["nats-jetstream-tls"] = BrokerNATSJetstreamSecure(client)
	brokers["nats-jetstream-tls-basic-auth"] = BrokerNATSJetstreamSecureBasicAuth(client)

	return brokers
}

func BrokerRedPandaAsKafka(client *dagger.Client) *dagger.Service {
	return client.Container().
		//	Set container image
		From(RedPandaImage).

		// Add Command
		WithExec([]string{
			"redpanda",
			"start",
			"--kafka-addr internal://0.0.0.0:9092,controller://0.0.0.0:9093",
			"--advertise-kafka-addr internal://kafka:9092,external://localhost:9093",
			// Mode dev-container uses well-known configuration properties
			// for development in containers.
			"--mode dev-container",
			// Tells Seastar (the framework Redpanda uses under the hood) to
			// use 1 core on the system.
			"--smp 1",
			"--default-log-level=info",
		}).

		// Return container as a service
		AsService()
}

// BrokerKafkaSecure returns a service for the Kafka broker secured with TLS.
func BrokerKafkaSecure(client *dagger.Client) *dagger.Service {
	key, cert, cacert, err := generateSelfSignedCertificateWithCA("kafka-tls")
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
	key, cert, cacert, err := generateSelfSignedCertificateWithCA("kafka-tls-basic-auth")
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
	key, cert, err := generateSelfSignedCertificate("nats-tls")
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
	key, cert, err := generateSelfSignedCertificate("nats-tls-basic-auth")
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
		WithExec([]string{"--tls", "--tlscert=/tls/server-cert.pem", "--tlskey=/tls/server-key.pem", "--user", "user", "--pass", "password"}). //nolint:lll
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
	key, cert, err := generateSelfSignedCertificate("nats-jetstream-tls-basic-auth")
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
	key, cert, err := generateSelfSignedCertificate("nats-jetstream-tls")
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

func certificateTemplateForHost(name string) x509.Certificate {
	return x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:       []string{"asyncapi-codegen"},
			OrganizationalUnit: []string{"localtest"},
			CommonName:         name,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, 1), // Valid for 1 day
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{name, "localhost", "127.0.0.1"},
	}
}

func generateSelfSignedCertificate(name string) ([]byte, []byte, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	template := certificateTemplateForHost(name)

	// Generate self-signed certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	// Encode private key to PEM format
	keyBytes := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	// Encode certificate to PEM format
	certBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	return keyBytes, certBytes, nil
}

func generateSelfSignedCertificateWithCA(name string) ([]byte, []byte, []byte, error) {
	// Generate private key for CA
	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create a self-signed CA certificate
	caCertTemplate := x509.Certificate{
		SerialNumber: big.NewInt(2), // Use a different serial number for the CA certificate
		Subject: pkix.Name{
			Organization: []string{"asyncapi-codegen"},
			CommonName:   "CA asyncapi-codegen",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // Valid for 10 years
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	// Generate self-signed CA certificate
	caDERBytes, err := x509.CreateCertificate(rand.Reader, &caCertTemplate, &caCertTemplate,
		&caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return nil, nil, nil, err
	}

	// Generate private key for server
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create server certificate signed by CA
	certTemplate := certificateTemplateForHost(name)

	derBytes, err := x509.CreateCertificate(rand.Reader, &certTemplate, &caCertTemplate,
		&privateKey.PublicKey, caPrivateKey)
	if err != nil {
		return nil, nil, nil, err
	}

	// Convert private key to PKCS #8
	privatKeyPKC8Bytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, nil, nil, err
	}

	// Encode server private key to PEM format
	keyBytes := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privatKeyPKC8Bytes})

	// Encode server certificate to PEM format
	certBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	// Encode CA certificate to PEM format
	caCertBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDERBytes})

	return keyBytes, certBytes, caCertBytes, nil
}
