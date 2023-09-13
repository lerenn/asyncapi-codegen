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

func (suite *Suite) TestExtensionsWithSchema() {
	// Set specification
	spec := asyncapi.Specification{
		Components: asyncapi.Components{
			Schemas: map[string]*asyncapi.Any{
				"flag": {
					Type:       asyncapi.TypeIsInteger.String(),
					Extensions: asyncapi.Extensions{ExtGoType: "uint8"},
				},
			},
		},
	}

	// Generate code and test result
	res, err := generators.TypesGenerator{Specification: spec}.Generate()
	suite.Require().NoError(err)
	suite.Require().True(regexp.MustCompile("FlagSchema +uint8").Match([]byte(res)))
}

func (suite *Suite) TestExtensionsWithObjectProperty() {
	// Set specification
	spec := asyncapi.Specification{
		Components: asyncapi.Components{
			Schemas: map[string]*asyncapi.Any{
				asyncapi.TypeIsObject.String(): {
					Type: asyncapi.TypeIsObject.String(),
					Properties: map[string]*asyncapi.Any{
						"flag": {
							Type:       asyncapi.TypeIsInteger.String(),
							Extensions: asyncapi.Extensions{ExtGoType: "uint8"},
						},
					},
					Required: []string{"flag"},
				},
			},
		},
	}

	// Generate code and test result
	res, err := generators.TypesGenerator{Specification: spec}.Generate()
	suite.Require().NoError(err)
	suite.Require().True(regexp.MustCompile("Flag +uint8").Match([]byte(res)))
}

func (suite *Suite) TestExtensionsWithArrayItem() {
	// Set specification
	spec := asyncapi.Specification{
		Components: asyncapi.Components{
			Schemas: map[string]*asyncapi.Any{
				"flags": {
					Type: asyncapi.TypeIsArray.String(),
					Items: &asyncapi.Any{
						Type:       asyncapi.TypeIsInteger.String(),
						Extensions: asyncapi.Extensions{ExtGoType: "uint8"},
					},
				},
			},
		},
	}

	// Generate code and test result
	res, err := generators.TypesGenerator{Specification: spec}.Generate()
	suite.Require().NoError(err)
	suite.Require().True(regexp.MustCompile(`FlagsSchema +\[\]uint8`).Match([]byte(res)))
}

func (suite *Suite) TestExtensionsWithObjectPropertyAndTypeFromPackage() {
	// Set specification
	spec := asyncapi.Specification{
		Components: asyncapi.Components{
			Schemas: map[string]*asyncapi.Any{
				asyncapi.TypeIsObject.String(): {
					Type: asyncapi.TypeIsObject.String(),
					Properties: map[string]*asyncapi.Any{
						"flag": {
							Type:       asyncapi.TypeIsInteger.String(),
							Extensions: asyncapi.Extensions{ExtGoType: "mypackage.Flag"},
						},
					},
					Required: []string{"flag"},
				},
			},
		},
	}

	// Generate code and test result
	res, err := generators.TypesGenerator{Specification: spec}.Generate()
	suite.Require().NoError(err)
	suite.Require().True(regexp.MustCompile(`Flag +mypackage.Flag`).Match([]byte(res)))
}
