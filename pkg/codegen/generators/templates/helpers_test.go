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
		// With acronym in middle
		{In: "IdTata", Out: "IDTata"},
		// With acronym in middle
		{In: "TotoIdLala", Out: "TotoIDLala"},
		// With acronym at the end
		{In: "TotoId", Out: "TotoID"},
		// Without acronym, but still the same letters as the acronym
		{In: "identity", Out: "Identity"},
		{In: "Identity", Out: "Identity"},
		{In: "covid", Out: "Covid"},
	}

	for i, c := range cases {
		suite.Require().Equal(c.Out, Namify(c.In), i)
	}
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
