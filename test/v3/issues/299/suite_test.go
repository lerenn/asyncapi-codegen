//go:generate go run ../../../../cmd/asyncapi-codegen -p issue299 -g types -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue299

import (
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/utils"
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

func (suite *Suite) TestSpecVersion310IsSupported() {
	// This package is generated from an `asyncapi: 3.1.0` specification.
	// Before support was added, generation failed with
	// `unsupported/invalid version: "3.1.0"` (issue #299). The package being
	// generated and compiling proves the 3.1.0 spec is parsed.
	msg := TestMessageMessage{}
	msg.Payload.Event = utils.ToPointer("ping")

	suite.Require().NotNil(msg.Payload.Event)
	suite.Require().Equal("ping", *msg.Payload.Event)
}
