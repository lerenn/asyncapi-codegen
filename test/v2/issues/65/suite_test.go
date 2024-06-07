package testutil

import (
	"regexp"
	"testing"

	"github.com/TheSadlig/asyncapi-codegen/pkg/asyncapi"
	asyncapiv2 "github.com/TheSadlig/asyncapi-codegen/pkg/asyncapi/v2"
	generatorv2 "github.com/TheSadlig/asyncapi-codegen/pkg/codegen/generators/v2"
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
	spec := asyncapiv2.Specification{
		Components: asyncapiv2.Components{
			Schemas: map[string]*asyncapiv2.Schema{
				"flag": {
					Name:       "FlagSchema",
					Type:       asyncapiv2.SchemaTypeIsInteger.String(),
					Extensions: asyncapiv2.Extensions{ExtGoType: "uint8"},
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
	spec := asyncapiv2.Specification{
		Components: asyncapiv2.Components{
			Schemas: map[string]*asyncapiv2.Schema{
				asyncapiv2.SchemaTypeIsObject.String(): {
					Type: asyncapiv2.SchemaTypeIsObject.String(),
					Properties: map[string]*asyncapiv2.Schema{
						"flag": {
							Name:       "Flag",
							Type:       asyncapiv2.SchemaTypeIsInteger.String(),
							Extensions: asyncapiv2.Extensions{ExtGoType: "uint8"},
						},
					},
					Validations: asyncapi.Validations[asyncapiv2.Schema]{
						Required: []string{"flag"},
					},
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
	spec := asyncapiv2.Specification{
		Components: asyncapiv2.Components{
			Schemas: map[string]*asyncapiv2.Schema{
				"flags": {
					Name: "FlagsSchema",
					Type: asyncapiv2.SchemaTypeIsArray.String(),
					Items: &asyncapiv2.Schema{
						Name:       "FlagSchema",
						Type:       asyncapiv2.SchemaTypeIsInteger.String(),
						Extensions: asyncapiv2.Extensions{ExtGoType: "uint8"},
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
	spec := asyncapiv2.Specification{
		Components: asyncapiv2.Components{
			Schemas: map[string]*asyncapiv2.Schema{
				asyncapiv2.SchemaTypeIsObject.String(): {
					Name: "ObjectSchema",
					Type: asyncapiv2.SchemaTypeIsObject.String(),
					Properties: map[string]*asyncapiv2.Schema{
						"flag": {
							Name:       "Flag",
							Type:       asyncapiv2.SchemaTypeIsInteger.String(),
							Extensions: asyncapiv2.Extensions{ExtGoType: "mypackage.Flag"},
						},
					},
					Validations: asyncapi.Validations[asyncapiv2.Schema]{
						Required: []string{"flag"},
					},
				},
			},
		},
	}

	// Generate code and test result
	res, err := generatorv2.TypesGenerator{Specification: spec}.Generate()
	suite.Require().NoError(err)
	suite.Require().True(regexp.MustCompile(`Flag +mypackage.Flag`).Match([]byte(res)))
}
