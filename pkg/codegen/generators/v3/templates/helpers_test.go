package templates

import (
	asyncapi2 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
	"testing"

	asyncapi "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
	"github.com/stretchr/testify/suite"
)

func TestHelpersSuite(t *testing.T) {
	suite.Run(t, new(HelpersSuite))
}

type HelpersSuite struct {
	suite.Suite
}

func (suite *HelpersSuite) TestIsRequired() {
	cases := []struct {
		Schema asyncapi.Schema
		Field  string
		Result bool
	}{
		// Is required
		{Schema: asyncapi.Schema{Validations: asyncapi2.Validations[asyncapi.Schema]{Required: []string{"field"}}}, Field: "field", Result: true},
		// Is not required
		{Schema: asyncapi.Schema{Validations: asyncapi2.Validations[asyncapi.Schema]{Required: []string{"another_field"}}}, Field: "field", Result: false},
	}

	for i, c := range cases {
		suite.Require().Equal(c.Result, IsRequired(c.Schema, c.Field), i)
	}
}

func (suite *HelpersSuite) TestGetChildrenObjectSchemas() {
	// TODO
}
