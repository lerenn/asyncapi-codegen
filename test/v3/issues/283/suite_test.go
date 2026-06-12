//go:generate go run ../../../../cmd/asyncapi-codegen -p issue283 -g types -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue283

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

// TestHeadersHaveDistinctTypes ensures messages that define a correlation ID
// without explicit headers each get their own header type instead of two
// types both named "Header" colliding in the same package (issue #283). If the
// bug were present, this package would not compile.
func (suite *Suite) TestHeadersHaveDistinctTypes() {
	ping := PingMessage{
		Headers: HeadersFromPingMessage{CorrelationId: utils.ToPointer("id-1")},
	}
	pong := PongMessage{
		Headers: HeadersFromPongMessage{CorrelationId: utils.ToPointer("id-2")},
	}

	suite.Require().Equal("id-1", *ping.Headers.CorrelationId)
	suite.Require().Equal("id-2", *pong.Headers.CorrelationId)
}
