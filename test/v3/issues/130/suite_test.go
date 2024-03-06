// More info: https://github.com/lerenn/asyncapi-codegen/issues/130

package issue130

import (
	"testing"

	asyncapi_test "github.com/lerenn/asyncapi-codegen/test"
	"github.com/lerenn/asyncapi-codegen/test/v3/issues/130/decoupling"
	"github.com/lerenn/asyncapi-codegen/test/v3/issues/130/requestreply"
	"github.com/lerenn/asyncapi-codegen/test/v3/issues/130/trait"
	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	brokers, cleanup := asyncapi_test.BrokerControllers(t)
	defer cleanup()

	for _, b := range brokers {
		suite.Run(t, decoupling.NewSuite(b))
		suite.Run(t, requestreply.NewSuite(b))
	}
	suite.Run(t, trait.NewSuite())
}
