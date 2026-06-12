//go:generate go run ../../../../cmd/asyncapi-codegen -p issue301 -g types -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue301

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

func (suite *Suite) TestArrayExamplePayloadParses() {
	// The specification used to generate this package declares a message whose
	// payload is an array and provides an example payload that is a JSON array.
	// Before the fix, parsing such a spec failed because MessageExample.Payload
	// was typed as map[string]any (issue #301). The mere fact that this package
	// is generated and compiles proves the array example payload is now parsed.
	msg := TestMessageMessage{
		Payload: []string{"first", "second"},
	}

	suite.Require().Len(msg.Payload, 2)
	suite.Require().Equal("first", msg.Payload[0])
}
