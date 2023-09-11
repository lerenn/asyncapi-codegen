//go:generate go run ../../../cmd/asyncapi-codegen -p issue49 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue49

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

type Suite struct {
	suite.Suite
}

func (suite *Suite) TestUserSubscriberGenerated() {
	// Check that the Subscriber is indeed generated
	var _ UserSubscriber
}
