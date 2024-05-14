package templates

import (
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
	asyncapiv3 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
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
		Schema asyncapiv3.Schema
		Field  string
		Result bool
	}{
		// Is required
		{
			Schema: asyncapiv3.Schema{
				Validations: asyncapi.Validations[asyncapiv3.Schema]{
					Required: []string{"field"},
				},
			},
			Field:  "field",
			Result: true,
		},
		// Is not required
		{
			Schema: asyncapiv3.Schema{
				Validations: asyncapi.Validations[asyncapiv3.Schema]{
					Required: []string{"another_field"},
				},
			},
			Field:  "field",
			Result: false,
		},
	}

	for i, c := range cases {
		suite.Require().Equal(c.Result, IsRequired(c.Schema, c.Field), i)
	}
}

func (suite *HelpersSuite) TestGetChildrenObjectSchemas() {
	// TODO
}
