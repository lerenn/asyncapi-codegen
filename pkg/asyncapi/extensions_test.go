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

func (suite *ExtensionsSuite) TestExtGoTypeImport() {
	tests := []struct {
		name     string
		schema   *Any
		expected []string
		error    bool
	}{
		// Schema
		{
			name: "flag",
			schema: &Any{
				Type: TypeIsInteger.String(),
				Extensions: Extensions{
					ExtGoType: "mypackage.Flag",
					ExtGoTypeImport: &GoTypeImportExtension{
						Path: "abc.xyz/repo/mypackage",
					},
				},
			},
			expected: []string{
				`"abc.xyz/repo/mypackage"`,
			},
		},

		// Object property, default name
		{
			name: TypeIsObject.String(),
			schema: &Any{
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
			expected: []string{
				`"abc.xyz/repo/mypackage"`,
			},
		},

		// Array item
		{
			name: "flags",
			schema: &Any{
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
			expected: []string{
				`"abc.xyz/repo/mypackage"`,
			},
		},

		// Object property, change package name
		{
			name: TypeIsObject.String(),
			schema: &Any{
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
			expected: []string{
				`alias "abc.xyz/repo/mypackage"`,
			},
		},

		// Object property, multiple imports
		{
			name: TypeIsObject.String(),
			schema: &Any{
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
			expected: []string{
				`alias "abc.xyz/repo/mypackage"`,
				`"github.com/rs/xid"`,
			},
		},

		// Conflicting import package name will error
		{
			name: TypeIsObject.String(),
			schema: &Any{
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
			error: true,
		},
	}

	for _, test := range tests {
		spec := Specification{
			Components: Components{
				Schemas: map[string]*Any{test.name: test.schema},
			},
		}

		res, err := spec.CustomImports()

		if test.error {
			suite.Require().True(err != nil)
		} else {
			suite.Require().NoError(err)

			sort.Strings(test.expected)
			sort.Strings(res)
			suite.Require().Equal(res, test.expected)
		}
	}
}
