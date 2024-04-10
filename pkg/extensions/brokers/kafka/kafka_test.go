package kafka

import (
	"crypto/tls"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSecureConnectionToKafka(t *testing.T) {

	t.Run("test connection is not successfully to TLS secured kafka broker without TLS config", func(t *testing.T) {
		_, err := NewController([]string{"kafka-tls:9092"}, WithGroupID("secureConnectTestWithoutTLS"))
		assert.Error(t, err, "new connection to TLS secured kafka broker without TLS config should return a error")
	})

	t.Run("test connection is successfully to TLS secured kafka broker with TLS config", func(t *testing.T) {
		_, err := NewController([]string{"kafka-tls:9092"}, WithGroupID("secureConnectTestWithTLS"),
			// just for testing use tls.Config with InsecureSkipVerify: true to skip server certificate validation for our self signed certificate
			WithTLS(&tls.Config{InsecureSkipVerify: true}),
		)
		assert.NoError(t, err, "new connection to TLS secured kafka broker with TLS config should return no error")
	})

	t.Run("test connection is not successfully to TLS secured kafka broker with TLS config and missing credentials", func(t *testing.T) {
		_, err := NewController([]string{"kafka-tls-basic-auth:9092"}, WithGroupID("secureConnectTestWithTLSAndWithoutCredentials"),
			// just for testing use tls.Config with InsecureSkipVerify: true to skip server certificate validation for our self signed certificate
			WithTLS(&tls.Config{InsecureSkipVerify: true}),
		)
		assert.Error(t, err, "new connection to TLS secured kafka broker with TLS config and missing credentials should return a error")
	})

	t.Run("test connection is successfully to TLS secured kafka broker with TLS config and credentials", func(t *testing.T) {
		sha512Mechanism, err := scram.Mechanism(scram.SHA512, "user", "password")
		assert.NoError(t, err, "new scram.SHA512 should not return a error")

		_, err = NewController([]string{"kafka-tls-basic-auth:9092"}, WithGroupID("secureConnectTestWithTLSAndWithCredentials"),
			// just for testing use tls.Config with InsecureSkipVerify: true to skip server certificate validation for our self signed certificate
			WithTLS(&tls.Config{InsecureSkipVerify: true}),
			WithSasl(sha512Mechanism),
		)
		assert.NoError(t, err, "new connection to TLS secured kafka broker with TLS config and basic credentials should return no error")
	})
}
