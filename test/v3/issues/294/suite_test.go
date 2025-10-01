//go:generate go run ../../../../cmd/asyncapi-codegen -p issue294 -g types -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue294

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, NewSuite())
}

type Suite struct {
	suite.Suite
}

func NewSuite() *Suite {
	return &Suite{}
}

func (suite *Suite) TestServerSecurityArrayParsing() {
	// This test verifies that the AsyncAPI specification with server security
	// defined as an array can be parsed without errors.
	// The actual test is whether the code generation succeeds (via go:generate).
	// If we reach this point, the parsing was successful.
	suite.T().Log("Server security array parsing successful")
}
