package asyncapi_test

import (
	"regexp"
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
	"github.com/lerenn/asyncapi-codegen/pkg/codegen/generators"
	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

type Suite struct {
	suite.Suite
}

func (suite *Suite) TestExtensions() {
	tests := []struct {
		name     string
		schema   *asyncapi.Any
		expected *regexp.Regexp
	}{
		// Schema
		{
			name: "flag",
			schema: &asyncapi.Any{
				Type:       "integer",
				Extensions: asyncapi.Extensions{ExtGoType: "uint8"},
			},
			expected: regexp.MustCompile("FlagSchema +uint8"),
		},

		// Object property
		{
			name: "object",
			schema: &asyncapi.Any{
				Type: "object",
				Properties: map[string]*asyncapi.Any{
					"flag": {
						Type:       "integer",
						Extensions: asyncapi.Extensions{ExtGoType: "uint8"},
					},
				},
				Required: []string{"flag"},
			},
			expected: regexp.MustCompile("Flag +uint8"),
		},

		// Array item
		{
			name: "flags",
			schema: &asyncapi.Any{
				Type: "array",
				Items: &asyncapi.Any{
					Type:       "integer",
					Extensions: asyncapi.Extensions{ExtGoType: "uint8"},
				},
			},
			expected: regexp.MustCompile(`FlagsSchema +\[\]uint8`),
		},

		// Object property, type from package
		{
			name: "object",
			schema: &asyncapi.Any{
				Type: "object",
				Properties: map[string]*asyncapi.Any{
					"flag": {
						Type:       "integer",
						Extensions: asyncapi.Extensions{ExtGoType: "mypackage.Flag"},
					},
				},
				Required: []string{"flag"},
			},
			expected: regexp.MustCompile(`Flag +mypackage.Flag`),
		},
	}

	for _, test := range tests {
		spec := asyncapi.Specification{
			Components: asyncapi.Components{
				Schemas: map[string]*asyncapi.Any{test.name: test.schema},
			},
		}
		res, err := generators.TypesGenerator{Specification: spec}.Generate()

		suite.Require().NoError(err)
		suite.Require().True(test.expected.Match([]byte(res)))
	}
}
