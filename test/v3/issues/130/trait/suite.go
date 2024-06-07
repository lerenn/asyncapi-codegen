//nolint:revive
package trait

import (
	"github.com/TheSadlig/asyncapi-codegen/pkg/asyncapi/parser"
	asyncapiv3 "github.com/TheSadlig/asyncapi-codegen/pkg/asyncapi/v3"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

func NewSuite() *Suite {
	return &Suite{}
}

func (suite *Suite) TestGenerateWithTrait() {
	// Generate codegen from file
	agnosticSpec, err := parser.FromFile(parser.FromFileParams{
		Path: "./trait/asyncapi.yaml",
	})
	suite.Require().NoError(err)

	// Process it to apply traits
	suite.Require().NoError(agnosticSpec.Process())

	// Get spec from codegen
	spec, ok := agnosticSpec.(*asyncapiv3.Specification)
	suite.Require().True(ok)

	// Check description hasn't change
	suite.Require().Equal("A longer description.", spec.Components.Messages["UserSignup"].Description)

	// Check summary has been applied
	suite.Require().Equal("Action to sign a user up.", spec.Components.Messages["UserSignup"].Summary)
}
