//go:generate go run ../../../../cmd/asyncapi-codegen -p issue138 -g types -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue138

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

// TestSetDefaultsFillsUnsetFields ensures the generated SetDefaults method fills
// optional fields that have a default value defined in the specification, while
// leaving fields that are already set untouched (issue #138).
func (suite *Suite) TestSetDefaultsFillsUnsetFields() {
	s := SettingsSchema{Name: "test"}
	s.SetDefaults()

	suite.Require().NotNil(s.Retries)
	suite.Require().Equal(int64(3), *s.Retries)

	suite.Require().NotNil(s.Ratio)
	suite.Require().Equal(float32(1.5), *s.Ratio)

	suite.Require().NotNil(s.Enabled)
	suite.Require().True(*s.Enabled)

	suite.Require().NotNil(s.Label)
	suite.Require().Equal("hello world", *s.Label)
}

// TestSetDefaultsKeepsExistingValues ensures already-set fields are not
// overwritten by their default value.
func (suite *Suite) TestSetDefaultsKeepsExistingValues() {
	retries := int64(10)
	s := SettingsSchema{Name: "test", Retries: &retries}
	s.SetDefaults()

	suite.Require().Equal(int64(10), *s.Retries)
}
