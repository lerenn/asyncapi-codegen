package generators

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestTemplatesSuite(t *testing.T) {
	suite.Run(t, new(TemplatesSuite))
}

type TemplatesSuite struct {
	suite.Suite
}

func (suite *TemplatesSuite) TestNamify() {
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
	}

	for i, c := range cases {
		suite.Require().Equal(c.Out, namify(c.In), i)
	}
}
