package main

import (
	"fmt"
	"os"
	"path/filepath"

	testutil "github.com/TheSadlig/asyncapi-codegen/pkg/utils/test"
)

func main() {
	createKafkaCerts()
	createNATSCerts()
}

func createNATSCerts() {
	// Set paths
	basePath := filepath.Join(".", "tmp", "certs", "nats")
	keyPath := filepath.Join(basePath, "server-key.pem")
	certPath := filepath.Join(basePath, "server-cert.pem")

	// Check if one file is missing
	if !checkIfOneOfFilesIsMissing(keyPath, certPath) {
		return
	}

	// Create directories
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		panic(err)
	}

	// Create NATS certs
	key, cert, err := testutil.GenerateSelfSignedCertificate("localhost")
	if err != nil {
		panic(fmt.Errorf("failed to generate self signed certificate: %w", err))
	}

	// Export NATS certs
	if err := os.WriteFile(keyPath, key, os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.WriteFile(certPath, cert, os.ModePerm); err != nil {
		panic(err)
	}
}

func createKafkaCerts() {
	basePath := filepath.Join(".", "tmp", "certs", "kafka")
	keyPath := filepath.Join(basePath, "kafka.keystore.key")
	certPath := filepath.Join(basePath, "kafka.keystore.pem")
	caPath := filepath.Join(basePath, "kafka.truststore.pem")

	// Check if one file is missing
	if !checkIfOneOfFilesIsMissing(keyPath, certPath, caPath) {
		return
	}

	// Create directories
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		panic(err)
	}

	// Create Kafka certs
	key, cert, cacert, err := testutil.GenerateSelfSignedCertificateWithCA("localhost")
	if err != nil {
		panic(fmt.Errorf("failed to generate self signed certificate: %w", err))
	}

	// Export Kafka certs
	if err := os.WriteFile(keyPath, key, os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.WriteFile(certPath, cert, os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.WriteFile(caPath, cacert, os.ModePerm); err != nil {
		panic(err)
	}
}

func checkIfOneOfFilesIsMissing(files ...string) bool {
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return true
		}
	}

	return false
}
