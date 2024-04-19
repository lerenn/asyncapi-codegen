package main

import (
	"fmt"
	"os"
	"path/filepath"

	testutil "github.com/lerenn/asyncapi-codegen/pkg/utils/test"
)

func main() {
	createKafkaCerts()
	createNATSCerts()
}

func createNATSCerts() {
	// Create NATS certs
	key, cert, err := testutil.GenerateSelfSignedCertificate("localhost")
	if err != nil {
		panic(fmt.Errorf("failed to generate self signed certificate: %w", err))
	}

	// Export NATS certs
	basePath := filepath.Join(".", "tmp", "certs", "nats")
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.WriteFile(filepath.Join(basePath, "server-key.pem"), key, os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.WriteFile(filepath.Join(basePath, "server-cert.pem"), cert, os.ModePerm); err != nil {
		panic(err)
	}
}

func createKafkaCerts() {
	// Create Kafka certs
	key, cert, cacert, err := testutil.GenerateSelfSignedCertificateWithCA("localhost")
	if err != nil {
		panic(fmt.Errorf("failed to generate self signed certificate: %w", err))
	}

	// Export Kafka certs
	basePath := filepath.Join(".", "tmp", "certs", "kafka")
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.WriteFile(filepath.Join(basePath, "kafka.keystore.key"), key, os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.WriteFile(filepath.Join(basePath, "kafka.keystore.pem"), cert, os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.WriteFile(filepath.Join(basePath, "kafka.truststore.pem"), cacert, os.ModePerm); err != nil {
		panic(err)
	}
}
