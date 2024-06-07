package templates

import (
	"testing"

	"github.com/TheSadlig/asyncapi-codegen/pkg/asyncapi"
	asyncapiv2 "github.com/TheSadlig/asyncapi-codegen/pkg/asyncapi/v2"
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
		Schema asyncapiv2.Schema
		Field  string
		Result bool
	}{
		// Is required
		{
			Schema: asyncapiv2.Schema{
				Validations: asyncapi.Validations[asyncapiv2.Schema]{
					Required: []string{"field"},
				},
			},
			Field:  "field",
			Result: true,
		},
		// Is not required
		{
			Schema: asyncapiv2.Schema{
				Validations: asyncapi.Validations[asyncapiv2.Schema]{
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

func (suite *HelpersSuite) TestOperationName() {
	cases := []struct {
		Channel asyncapiv2.Channel
		Result  string
	}{
		// By default
		{
			Channel: asyncapiv2.Channel{Name: "Default"},
			Result:  "Default",
		},
		// With subscribe but no ooperation ID
		{
			Channel: asyncapiv2.Channel{
				Subscribe: &asyncapiv2.Operation{},
				Name:      "Default",
			},
			Result: "Default",
		},
		// With subscribe and operation ID
		{
			Channel: asyncapiv2.Channel{
				Subscribe: &asyncapiv2.Operation{
					OperationID: "Subscribe",
				},
				Name: "Default",
			},
			Result: "Subscribe",
		},
		// With publish but no ooperation ID
		{
			Channel: asyncapiv2.Channel{
				Publish: &asyncapiv2.Operation{},
				Name:    "Default",
			},
			Result: "Default",
		},
		// With publish and operation ID
		{
			Channel: asyncapiv2.Channel{
				Publish: &asyncapiv2.Operation{
					OperationID: "Publish",
				},
				Name: "Default",
			},
			Result: "Publish",
		},
	}

	for i, c := range cases {
		suite.Require().Equal(c.Result, OperationName(c.Channel), i)
	}
}

func (suite *HelpersSuite) TestGetChildrenObjectSchemas() {
	// TODO
}
