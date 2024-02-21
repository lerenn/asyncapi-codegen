//go:generate go run ../../../../cmd/asyncapi-codegen -g types -p issue114 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue114

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

func (suite *Suite) TestCorrectVersion() {
	suite.Require().Equal("1.2.3", AsyncAPIVersion)
}
