package templates

import (
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
	"github.com/stretchr/testify/suite"
)

func TestHelpersSuite(t *testing.T) {
	suite.Run(t, new(HelpersSuite))
}

type HelpersSuite struct {
	suite.Suite
}

func (suite *HelpersSuite) TestNamify() {
	cases := []struct {
		In  string
		Out string
	}{
		// Remove leading digits
		{In: "0name0", Out: "Name0"},
		// Remove non alphanumerics
		{In: "?#!name", Out: "Name"},
		// Capitalize
		{In: "name", Out: "Name"},
		// Snake Case
		{In: "eh_oh__ah", Out: "EhOhAh"},
		// With acronym
		{In: "TotoIdLala", Out: "TotoIDLala"},
		{In: "TotoId", Out: "TotoID"},
	}

	for i, c := range cases {
		suite.Require().Equal(c.Out, Namify(c.In), i)
	}
}

func (suite *HelpersSuite) TestIsRequired() {
	cases := []struct {
		Any    asyncapi.Any
		Field  string
		Result bool
	}{
		// Is required
		{Any: asyncapi.Any{Required: []string{"field"}}, Field: "field", Result: true},
		// Is not required
		{Any: asyncapi.Any{Required: []string{"another_field"}}, Field: "field", Result: false},
	}

	for i, c := range cases {
		suite.Require().Equal(c.Result, IsRequired(c.Any, c.Field), i)
	}
}
