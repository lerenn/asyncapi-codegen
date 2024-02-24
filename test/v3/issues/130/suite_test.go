// More info: https://github.com/lerenn/asyncapi-codegen/issues/130

package issue130

import (
	"testing"

	asyncapi_test "github.com/lerenn/asyncapi-codegen/test"
	decoupling "github.com/lerenn/asyncapi-codegen/test/v3/issues/130/decoupling"
	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	brokers, cleanup := asyncapi_test.BrokerControllers(t)
	defer cleanup()

	// Only do it with one broker as this is not testing the broker
	suite.Run(t, decoupling.NewSuite(brokers[0]))
}
