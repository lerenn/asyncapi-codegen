package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	testutil "github.com/lerenn/asyncapi-codegen/pkg/utils/test"
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
	keystorePath := filepath.Join(basePath, "kafka.keystore.jks")

	// Check if keystore is missing
	if _, err := os.Stat(keystorePath); os.IsNotExist(err) {
		generateKafkaKeystore(basePath)
	}
}

func generateKafkaKeystore(basePath string) {
	keystorePath := filepath.Join(basePath, "kafka.keystore.jks")
	keystorePasswordPath := filepath.Join(basePath, "kafka.keystore.jks.password")
	truststorePath := filepath.Join(basePath, "kafka.truststore.jks")
	truststorePasswordPath := filepath.Join(basePath, "kafka.truststore.jks.password")

	// Check if Java is available, otherwise we'll use Docker
	useDocker := exec.Command("java", "-version").Run() != nil
	// Create directories
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		panic(err)
	}

	// Create Kafka certs
	key, cert, cacert, err := testutil.GenerateSelfSignedCertificateWithCA("localhost")
	if err != nil {
		panic(fmt.Errorf("failed to generate self signed certificate: %w", err))
	}

	// Write temporary PEM files
	keyPath := filepath.Join(basePath, "key.pem")
	certPath := filepath.Join(basePath, "cert.pem")
	caPath := filepath.Join(basePath, "cacert.pem")
	chainPath := filepath.Join(basePath, "chain.pem")
	p12Path := filepath.Join(basePath, "keystore.p12")

	writePEMFiles(keyPath, certPath, caPath, key, cert, cacert)

	keystorePassword := "changeit"
	truststorePassword := "changeit"

	// Create certificate chain file (cert + CA cert)
	createCertificateChain(chainPath, cert, cacert)

	// Convert PEM to PKCS12 with certificate chain
	convertPEMToPKCS12(chainPath, keyPath, p12Path, keystorePassword)

	// Convert PKCS12 to JKS keystore
	convertPKCS12ToJKS(useDocker, basePath, p12Path, keystorePath, keystorePassword)

	// Import CA cert as trusted certificate into keystore
	importCACert(useDocker, basePath, caPath, keystorePath, keystorePassword)

	// Create symlink from truststore to keystore (use relative path)
	createTruststoreSymlink(truststorePath, keystorePath)

	// Create password files
	writePasswordFiles(keystorePasswordPath, truststorePasswordPath, keystorePassword, truststorePassword)

	// Clean up temporary files
	cleanupTempFiles(keyPath, certPath, caPath, chainPath, p12Path)
}

func writePEMFiles(keyPath, certPath, caPath string, key, cert, cacert []byte) {
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

func createCertificateChain(chainPath string, cert, cacert []byte) {
	chainContent := make([]byte, 0, len(cert)+len(cacert))
	chainContent = append(chainContent, cert...)
	chainContent = append(chainContent, cacert...)
	if err := os.WriteFile(chainPath, chainContent, os.ModePerm); err != nil {
		panic(err)
	}
}

func convertPEMToPKCS12(chainPath, keyPath, p12Path, keystorePassword string) {
	cmd := exec.Command("openssl", "pkcs12", "-export",
		"-in", chainPath,
		"-inkey", keyPath,
		"-out", p12Path,
		"-password", "pass:"+keystorePassword,
		"-name", "kafka")
	var opensslStderr bytes.Buffer
	cmd.Stderr = &opensslStderr
	if err := cmd.Run(); err != nil {
		panic(fmt.Errorf("failed to create PKCS12 keystore: %w\nstderr: %s", err, opensslStderr.String()))
	}
}

func convertPKCS12ToJKS(useDocker bool, basePath, p12Path, keystorePath, keystorePassword string) {
	var p12RelPath, keystoreRelPath string
	if useDocker {
		p12RelPath = filepath.Base(p12Path)
		keystoreRelPath = filepath.Base(keystorePath)
	} else {
		p12RelPath = p12Path
		keystoreRelPath = keystorePath
	}

	if err := runKeytool(useDocker, basePath,
		"-importkeystore",
		"-srckeystore", p12RelPath,
		"-srcstoretype", "PKCS12",
		"-destkeystore", keystoreRelPath,
		"-deststoretype", "JKS",
		"-srcstorepass", keystorePassword,
		"-deststorepass", keystorePassword,
		"-noprompt"); err != nil {
		panic(fmt.Errorf("failed to convert PKCS12 to JKS: %w", err))
	}
}

func importCACert(useDocker bool, basePath, caPath, keystorePath, keystorePassword string) {
	var caRelPath, keystoreRelPath string
	if useDocker {
		caRelPath = filepath.Base(caPath)
		keystoreRelPath = filepath.Base(keystorePath)
	} else {
		caRelPath = caPath
		keystoreRelPath = keystorePath
	}

	if err := runKeytool(useDocker, basePath,
		"-import", "-trustcacerts",
		"-alias", "ca",
		"-file", caRelPath,
		"-keystore", keystoreRelPath,
		"-storepass", keystorePassword,
		"-noprompt"); err != nil {
		panic(fmt.Errorf("failed to import CA cert into keystore: %w", err))
	}
}

func createTruststoreSymlink(truststorePath, keystorePath string) {
	if err := os.Remove(truststorePath); err != nil && !os.IsNotExist(err) {
		panic(fmt.Errorf("failed to remove existing truststore: %w", err))
	}
	keystoreSymlinkPath := filepath.Base(keystorePath)
	if err := os.Symlink(keystoreSymlinkPath, truststorePath); err != nil {
		panic(fmt.Errorf("failed to create truststore symlink: %w", err))
	}
}

func writePasswordFiles(keystorePasswordPath, truststorePasswordPath, keystorePassword, truststorePassword string) {
	if err := os.WriteFile(keystorePasswordPath, []byte(keystorePassword), os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.WriteFile(truststorePasswordPath, []byte(truststorePassword), os.ModePerm); err != nil {
		panic(err)
	}
}

func cleanupTempFiles(keyPath, certPath, caPath, chainPath, p12Path string) {
	os.Remove(keyPath)
	os.Remove(certPath)
	os.Remove(caPath)
	os.Remove(chainPath)
	os.Remove(p12Path)
}

func runKeytool(useDocker bool, workDir string, args ...string) error {
	var cmd *exec.Cmd
	if useDocker {
		// Use Docker to run keytool in a Java container
		absWorkDir, err := filepath.Abs(workDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}
		dockerArgs := []string{"run", "--rm",
			"-v", absWorkDir + ":/work",
			"-w", "/work",
			"eclipse-temurin:17-jre",
			"keytool"}
		dockerArgs = append(dockerArgs, args...)
		cmd = exec.Command("docker", dockerArgs...)
	} else {
		cmd = exec.Command("keytool", args...)
		cmd.Dir = workDir
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w\nstderr: %s", err, stderr.String())
	}
	return nil
}

func checkIfOneOfFilesIsMissing(files ...string) bool {
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return true
		}
	}

	return false
}
