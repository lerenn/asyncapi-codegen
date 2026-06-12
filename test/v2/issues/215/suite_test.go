//go:generate go run ../../../../cmd/asyncapi-codegen -p issue215 -g types -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue215

import (
	"os"
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

// TestTypesOnlyHasNoControllerBoilerplate ensures that `-g types` output does
// not contain controller/application boilerplate (issue #215). The internal
// controller struct and its options now live with the generated controllers,
// so a types-only package can be combined with others without dead code or
// symbol collisions.
func (suite *Suite) TestTypesOnlyHasNoControllerBoilerplate() {
	content, err := os.ReadFile("asyncapi.gen.go")
	suite.Require().NoError(err)

	generated := string(content)
	suite.Require().NotContains(generated, "type controller struct")
	suite.Require().NotContains(generated, "type ControllerOption")
	suite.Require().NotContains(generated, "func WithLogger(")
	suite.Require().NotContains(generated, "func WithMiddlewares(")
	suite.Require().NotContains(generated, "func WithErrorHandler(")
}
