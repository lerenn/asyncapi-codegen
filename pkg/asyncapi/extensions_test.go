package asyncapi

import (
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
			Schemas: map[string]*Any{
				"flag": {
					Type: TypeIsInteger.String(),
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
			Schemas: map[string]*Any{
				TypeIsObject.String(): {
					Type: TypeIsObject.String(),
					Properties: map[string]*Any{
						"flag": {
							Type: TypeIsInteger.String(),
							Extensions: Extensions{
								ExtGoType: "mypackage.Flag",
								ExtGoTypeImport: &GoTypeImportExtension{
									Path: "abc.xyz/repo/mypackage",
								},
							},
						},
					},
					Required: []string{"flag"},
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
			Schemas: map[string]*Any{
				"flag": {
					Type: TypeIsArray.String(),
					Items: &Any{
						Type: TypeIsInteger.String(),
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
			Schemas: map[string]*Any{
				TypeIsObject.String(): {
					Type: TypeIsObject.String(),
					Properties: map[string]*Any{
						"flag": {
							Type: TypeIsInteger.String(),
							Extensions: Extensions{
								ExtGoType: "alias.Flag",
								ExtGoTypeImport: &GoTypeImportExtension{
									Path: "abc.xyz/repo/mypackage",
									Name: "alias",
								},
							},
						},
					},
					Required: []string{"flag"},
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
			Schemas: map[string]*Any{
				TypeIsObject.String(): {
					Type: TypeIsObject.String(),
					Properties: map[string]*Any{
						"flag": {
							Type: TypeIsInteger.String(),
							Extensions: Extensions{
								ExtGoType: "alias.Flag",
								ExtGoTypeImport: &GoTypeImportExtension{
									Path: "abc.xyz/repo/mypackage",
									Name: "alias",
								},
							},
						},
						"id": {
							Type: TypeIsString.String(),
							Extensions: Extensions{
								ExtGoType: "xid.ID",
								ExtGoTypeImport: &GoTypeImportExtension{
									Path: "github.com/rs/xid",
								},
							},
						},
					},
					Required: []string{"flag"},
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
			Schemas: map[string]*Any{
				TypeIsObject.String(): {
					Type: TypeIsObject.String(),
					Properties: map[string]*Any{
						"start_flag": {
							Type: TypeIsInteger.String(),
							Extensions: Extensions{
								ExtGoType: "alias2.Flag",
								ExtGoTypeImport: &GoTypeImportExtension{
									Path: "abc.xyz/repo/mypackage",
									Name: "alias1",
								},
							},
						},
						"end_flag": {
							Type: TypeIsInteger.String(),
							Extensions: Extensions{
								ExtGoType: "alias2.Flag",
								ExtGoTypeImport: &GoTypeImportExtension{
									Path: "abc.xyz/repo/mypackage",
									Name: "alias2",
								},
							},
						},
					},
					Required: []string{"start_flag", "end_flag"},
				},
			},
		},
	}

	// Process custom imports
	_, err := spec.CustomImports()

	// It should be an error
	suite.Require().Error(err)
}
