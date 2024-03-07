package templates

import (
	"testing"

	asyncapi "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v2"
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
		{Schema: asyncapi.Schema{Required: []string{"field"}}, Field: "field", Result: true},
		// Is not required
		{Schema: asyncapi.Schema{Required: []string{"another_field"}}, Field: "field", Result: false},
	}

	for i, c := range cases {
		suite.Require().Equal(c.Result, IsRequired(c.Schema, c.Field), i)
	}
}

func (suite *HelpersSuite) TestOperationName() {
	cases := []struct {
		Channel asyncapi.Channel
		Result  string
	}{
		// By default
		{
			Channel: asyncapi.Channel{Name: "Default"},
			Result:  "Default",
		},
		// With subscribe but no ooperation ID
		{
			Channel: asyncapi.Channel{
				Subscribe: &asyncapi.Operation{},
				Name:      "Default",
			},
			Result: "Default",
		},
		// With subscribe and operation ID
		{
			Channel: asyncapi.Channel{
				Subscribe: &asyncapi.Operation{
					OperationID: "Subscribe",
				},
				Name: "Default",
			},
			Result: "Subscribe",
		},
		// With publish but no ooperation ID
		{
			Channel: asyncapi.Channel{
				Publish: &asyncapi.Operation{},
				Name:    "Default",
			},
			Result: "Default",
		},
		// With publish and operation ID
		{
			Channel: asyncapi.Channel{
				Publish: &asyncapi.Operation{
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
