package template

import (
	"testing"

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
	// new line
	{In: "\n", Out: ""},
	// just numbers
	{In: "00", Out: ""},
	// Remove leading digits
	{In: "0name0", Out: "Name0"},
	// Remove non alphanumerics
	{In: "?#!name", Out: "Name"},
	// Capitalize
	{In: "name", Out: "Name"},
	// Snake Case
	{In: "eh_oh__ah", Out: "EhOhAh"},
	// Weird delimiters
	{In: "eh.oh__ah", Out: "EhOhAh"},
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
