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

type namifyCases struct {
	In  string
	Out string
}

var namifyBaseCases = []namifyCases{
	// Nothing
	{In: "", Out: ""},
	// Remove leading digits
	{In: "0name0", Out: "Name0"},
	// Remove non alphanumerics
	{In: "?#!name", Out: "Name"},
	// Capitalize
	{In: "name", Out: "Name"},
	// Snake Case
	{In: "eh_oh__ah", Out: "EhOhAh"},
	// With acronym in middle
	{In: "IDTata", Out: "IDTata"},
	// With acronym in middle
	{In: "TotoIDLala", Out: "TotoIDLala"},
	{In: "Toto_IDLala", Out: "TotoIDLala"},
	{In: "TotoSMALala", Out: "TotoSMALala"},
	// With acronym at the end
	{In: "TotoID", Out: "TotoID"},
	{In: "Toto_ID", Out: "TotoID"},
	// Without acronym, but still the same letters as the acronym
	{In: "identity", Out: "Identity"},
	{In: "Identity", Out: "Identity"},
	{In: "covid", Out: "Covid"},
}

func (suite *HelpersSuite) TestNamify() {
	for i, c := range namifyBaseCases {
		suite.Require().Equal(c.Out, Namify(c.In), i)
	}
}

func (suite *HelpersSuite) TestNamifyWithoutParams() {
	cases := []namifyCases{
		// With argument
		{In: "name.{id}", Out: "Name"},
	}

	for i, c := range cases {
		suite.Require().Equal(c.Out, NamifyWithoutParams(c.In), i)
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
