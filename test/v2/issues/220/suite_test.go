//go:generate go run ../../../../cmd/asyncapi-codegen -p camel -n camel -i ./asyncapi.yaml -o ./camel/asyncapi.gen.go
//go:generate go run ../../../../cmd/asyncapi-codegen -p none -n none -i ./asyncapi.yaml -o ./none/asyncapi.gen.go

package issue220

import (
	"testing"

	"github.com/TheSadlig/asyncapi-codegen/test/v2/issues/220/camel"
	"github.com/TheSadlig/asyncapi-codegen/test/v2/issues/220/none"
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

func (suite *Suite) TestCamel() {
	// Checking if the schema has been generated with the correct naming scheme
	_ = camel.TestSchema{
		AnotherProp2: nil,
		AProp1:       nil,
	}
}

func (suite *Suite) TestNone() {
	// Checking if the schema has been generated with the correct naming scheme
	_ = none.TESTSchema{
		ANOTHERPROP2: nil,
		APROP1:       nil,
	}
}
