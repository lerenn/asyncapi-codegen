package asyncapiv2

import (
	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
	"sort"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestExtensionsSuite(t *testing.T) {
	suite.Run(t, new(ExtensionsSuite))
}

type ExtensionsSuite struct {
	suite.Suite
}

func (suite *ExtensionsSuite) TextExtGoTypeImportWithSchema() {
	// Set specification
	spec := Specification{
		Components: Components{
			Schemas: map[string]*Schema{
				"flag": {
					Type: SchemaTypeIsInteger.String(),
					Extensions: Extensions{
						ExtGoType: "mypackage.Flag",
						ExtGoTypeImport: &GoTypeImportExtension{
							Path: "abc.xyz/repo/mypackage",
						},
					},
				},
			},
		},
	}

	// Process custom imports
	res, err := spec.CustomImports()
	suite.Require().NoError(err)

	// Check result
	sort.Strings(res)
	suite.Require().Equal([]string{`"abc.xyz/repo/mypackage"`}, res)
}

func (suite *ExtensionsSuite) TextExtGoTypeImportWithObjectProperty() {
	// Set specification
	spec := Specification{
		Components: Components{
			Schemas: map[string]*Schema{
				SchemaTypeIsObject.String(): {
					Type: SchemaTypeIsObject.String(),
					Properties: map[string]*Schema{
						"flag": {
							Type: SchemaTypeIsInteger.String(),
							Extensions: Extensions{
								ExtGoType: "mypackage.Flag",
								ExtGoTypeImport: &GoTypeImportExtension{
									Path: "abc.xyz/repo/mypackage",
								},
							},
						},
					},
					Validations: asyncapi.Validations[Schema]{
						Required: []string{"flag"},
					},
				},
			},
		},
	}

	// Process custom imports
	res, err := spec.CustomImports()
	suite.Require().NoError(err)

	// Check result
	sort.Strings(res)
	suite.Require().Equal([]string{`"abc.xyz/repo/mypackage"`}, res)
}

func (suite *ExtensionsSuite) TextExtGoTypeImportWithArrayItem() {
	// Set specification
	spec := Specification{
		Components: Components{
			Schemas: map[string]*Schema{
				"flag": {
					Type: SchemaTypeIsArray.String(),
					Items: &Schema{
						Type: SchemaTypeIsInteger.String(),
						Extensions: Extensions{
							ExtGoType: "mypackage.Flag",
							ExtGoTypeImport: &GoTypeImportExtension{
								Path: "abc.xyz/repo/mypackage",
							},
						},
					},
				},
			},
		},
	}

	// Process custom imports
	res, err := spec.CustomImports()
	suite.Require().NoError(err)

	// Check result
	sort.Strings(res)
	suite.Require().Equal([]string{`"abc.xyz/repo/mypackage"`}, res)
}

func (suite *ExtensionsSuite) TextExtGoTypeImportWithObjectPropertyAndDifferentPackageName() {
	// Set specification
	spec := Specification{
		Components: Components{
			Schemas: map[string]*Schema{
				SchemaTypeIsObject.String(): {
					Type: SchemaTypeIsObject.String(),
					Properties: map[string]*Schema{
						"flag": {
							Type: SchemaTypeIsInteger.String(),
							Extensions: Extensions{
								ExtGoType: "alias.Flag",
								ExtGoTypeImport: &GoTypeImportExtension{
									Path: "abc.xyz/repo/mypackage",
									Name: "alias",
								},
							},
						},
					},
					Validations: asyncapi.Validations[Schema]{
						Required: []string{"flag"},
					},
				},
			},
		},
	}

	// Process custom imports
	res, err := spec.CustomImports()
	suite.Require().NoError(err)

	// Check result
	sort.Strings(res)
	suite.Require().Equal([]string{`alias "abc.xyz/repo/mypackage"`}, res)
}

func (suite *ExtensionsSuite) TextExtGoTypeImportWithObjectPropertyAndMultipleImports() {
	// Set specification
	spec := Specification{
		Components: Components{
			Schemas: map[string]*Schema{
				SchemaTypeIsObject.String(): {
					Type: SchemaTypeIsObject.String(),
					Properties: map[string]*Schema{
						"flag": {
							Type: SchemaTypeIsInteger.String(),
							Extensions: Extensions{
								ExtGoType: "alias.Flag",
								ExtGoTypeImport: &GoTypeImportExtension{
									Path: "abc.xyz/repo/mypackage",
									Name: "alias",
								},
							},
						},
						"id": {
							Type: SchemaTypeIsString.String(),
							Extensions: Extensions{
								ExtGoType: "xid.ID",
								ExtGoTypeImport: &GoTypeImportExtension{
									Path: "github.com/rs/xid",
								},
							},
						},
					},
					Validations: asyncapi.Validations[Schema]{
						Required: []string{"flag"},
					},
				},
			},
		},
	}

	// Process custom imports
	res, err := spec.CustomImports()
	suite.Require().NoError(err)

	// Check result
	sort.Strings(res)
	suite.Require().Equal([]string{`alias "abc.xyz/repo/mypackage"`, `"github.com/rs/xid"`}, res)
}

func (suite *ExtensionsSuite) TextExtGoTypeImportWithConflictingPackageNames() {
	// Set specification
	spec := Specification{
		Components: Components{
			Schemas: map[string]*Schema{
				SchemaTypeIsObject.String(): {
					Type: SchemaTypeIsObject.String(),
					Properties: map[string]*Schema{
						"start_flag": {
							Type: SchemaTypeIsInteger.String(),
							Extensions: Extensions{
								ExtGoType: "alias2.Flag",
								ExtGoTypeImport: &GoTypeImportExtension{
									Path: "abc.xyz/repo/mypackage",
									Name: "alias1",
								},
							},
						},
						"end_flag": {
							Type: SchemaTypeIsInteger.String(),
							Extensions: Extensions{
								ExtGoType: "alias2.Flag",
								ExtGoTypeImport: &GoTypeImportExtension{
									Path: "abc.xyz/repo/mypackage",
									Name: "alias2",
								},
							},
						},
					},
					Validations: asyncapi.Validations[Schema]{
						Required: []string{"start_flag", "end_flag"},
					},
				},
			},
		},
	}

	// Process custom imports
	_, err := spec.CustomImports()

	// It should be an error
	suite.Require().Error(err)
}
