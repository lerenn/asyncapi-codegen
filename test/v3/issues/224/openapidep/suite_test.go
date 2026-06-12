//go:generate go run ../../../../../cmd/asyncapi-codegen -p openapidep -g types -i ./asyncapi.yaml -i ./openapi.yaml -o ./asyncapi.gen.go

package openapidep

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

// TestOpenAPIDependencySchemaIsGenerated ensures an OpenAPI document can be used
// as a schema-only dependency: its components/schemas are parsed (ignoring the
// array-typed servers, paths, etc. that AsyncAPI does not allow) and the
// referenced schema is generated and usable as a message payload (issue #224).
func (suite *Suite) TestOpenAPIDependencySchemaIsGenerated() {
	msg := TestMessageMessage{
		Payload: PetSchema{
			Name: utils.ToPointer("Rex"),
			Age:  utils.ToPointer(int64(3)),
		},
	}

	suite.Require().Equal("Rex", *msg.Payload.Name)
	suite.Require().Equal(int64(3), *msg.Payload.Age)
}
