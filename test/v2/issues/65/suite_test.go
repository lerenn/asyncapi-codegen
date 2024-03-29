package asyncapi_test

import (
	"regexp"
	"testing"

	asyncapi "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v2"
	generatorv2 "github.com/lerenn/asyncapi-codegen/pkg/codegen/generators/v2"
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
			Schemas: map[string]*asyncapi.Schema{
				"flag": {
					Name:       "FlagSchema",
					Type:       asyncapi.SchemaTypeIsInteger.String(),
					Extensions: asyncapi.Extensions{ExtGoType: "uint8"},
				},
			},
		},
	}

	// Generate code and test result
	res, err := generatorv2.TypesGenerator{Specification: spec}.Generate()
	suite.Require().NoError(err)
	suite.Require().True(regexp.MustCompile("FlagSchema +uint8").Match([]byte(res)))
}

func (suite *Suite) TestExtensionsWithObjectProperty() {
	// Set specification
	spec := asyncapi.Specification{
		Components: asyncapi.Components{
			Schemas: map[string]*asyncapi.Schema{
				asyncapi.SchemaTypeIsObject.String(): {
					Type: asyncapi.SchemaTypeIsObject.String(),
					Properties: map[string]*asyncapi.Schema{
						"flag": {
							Name:       "Flag",
							Type:       asyncapi.SchemaTypeIsInteger.String(),
							Extensions: asyncapi.Extensions{ExtGoType: "uint8"},
						},
					},
					Required: []string{"flag"},
				},
			},
		},
	}

	// Generate code and test result
	res, err := generatorv2.TypesGenerator{Specification: spec}.Generate()
	suite.Require().NoError(err)
	suite.Require().True(regexp.MustCompile("Flag +uint8").Match([]byte(res)))
}

func (suite *Suite) TestExtensionsWithArrayItem() {
	// Set specification
	spec := asyncapi.Specification{
		Components: asyncapi.Components{
			Schemas: map[string]*asyncapi.Schema{
				"flags": {
					Name: "FlagsSchema",
					Type: asyncapi.SchemaTypeIsArray.String(),
					Items: &asyncapi.Schema{
						Name:       "FlagSchema",
						Type:       asyncapi.SchemaTypeIsInteger.String(),
						Extensions: asyncapi.Extensions{ExtGoType: "uint8"},
					},
				},
			},
		},
	}

	// Generate code and test result
	res, err := generatorv2.TypesGenerator{Specification: spec}.Generate()
	suite.Require().NoError(err)
	suite.Require().True(regexp.MustCompile(`FlagsSchema +\[\]uint8`).Match([]byte(res)))
}

func (suite *Suite) TestExtensionsWithObjectPropertyAndTypeFromPackage() {
	// Set specification
	spec := asyncapi.Specification{
		Components: asyncapi.Components{
			Schemas: map[string]*asyncapi.Schema{
				asyncapi.SchemaTypeIsObject.String(): {
					Name: "ObjectSchema",
					Type: asyncapi.SchemaTypeIsObject.String(),
					Properties: map[string]*asyncapi.Schema{
						"flag": {
							Name:       "Flag",
							Type:       asyncapi.SchemaTypeIsInteger.String(),
							Extensions: asyncapi.Extensions{ExtGoType: "mypackage.Flag"},
						},
					},
					Required: []string{"flag"},
				},
			},
		},
	}

	// Generate code and test result
	res, err := generatorv2.TypesGenerator{Specification: spec}.Generate()
	suite.Require().NoError(err)
	suite.Require().True(regexp.MustCompile(`Flag +mypackage.Flag`).Match([]byte(res)))
}
