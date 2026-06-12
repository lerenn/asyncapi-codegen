//go:generate go run ../../../../cmd/asyncapi-codegen -p issue296 -g types -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue296

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

func (suite *Suite) TestHeadersFromReferenceAreGenerated() {
	// Headers defined as a reference object must be resolved and generated as
	// their own header type instead of failing parsing (issue #296).
	msg := TestMessageMessage{
		Headers: CommonHeadersSchema{
			ReplyTo: utils.ToPointer("reply-to-address"),
		},
	}

	suite.Require().NotNil(msg.Headers.ReplyTo)
	suite.Require().Equal("reply-to-address", *msg.Headers.ReplyTo)
}
