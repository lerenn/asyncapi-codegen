package test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

// GenerateSelfSignedCertificate generates a self-signed certificate.
func GenerateSelfSignedCertificate(name string) ([]byte, []byte, error) {
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

// GenerateSelfSignedCertificateWithCA generates a self-signed certificate with a CA certificate.
func GenerateSelfSignedCertificateWithCA(name string) ([]byte, []byte, []byte, error) {
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
