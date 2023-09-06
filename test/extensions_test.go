package asyncapi_test

import (
	"regexp"
	"sort"
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
	"github.com/lerenn/asyncapi-codegen/pkg/codegen/generators"
	"github.com/matryer/is"
)

func TestExtGoType(t *testing.T) {
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

	is := is.New(t)

	for _, test := range tests {
		spec := asyncapi.Specification{
			Components: asyncapi.Components{
				Schemas: map[string]*asyncapi.Any{test.name: test.schema},
			},
		}
		res, err := generators.TypesGenerator{Specification: spec}.Generate()

		is.NoErr(err)
		is.True(test.expected.Match([]byte(res)))
	}
}

func TestExtGoTypeImport(t *testing.T) {
	tests := []struct {
		name     string
		schema   *asyncapi.Any
		expected []string
		error    bool
	}{
		// Schema
		{
			name: "flag",
			schema: &asyncapi.Any{
				Type: "integer",
				Extensions: asyncapi.Extensions{
					ExtGoType: "mypackage.Flag",
					ExtGoTypeImport: &asyncapi.GoTypeImportExtension{
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
			name: "object",
			schema: &asyncapi.Any{
				Type: "object",
				Properties: map[string]*asyncapi.Any{
					"flag": {
						Type: "integer",
						Extensions: asyncapi.Extensions{
							ExtGoType: "mypackage.Flag",
							ExtGoTypeImport: &asyncapi.GoTypeImportExtension{
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
			schema: &asyncapi.Any{
				Type: "array",
				Items: &asyncapi.Any{
					Type: "integer",
					Extensions: asyncapi.Extensions{
						ExtGoType: "mypackage.Flag",
						ExtGoTypeImport: &asyncapi.GoTypeImportExtension{
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
			name: "object",
			schema: &asyncapi.Any{
				Type: "object",
				Properties: map[string]*asyncapi.Any{
					"flag": {
						Type: "integer",
						Extensions: asyncapi.Extensions{
							ExtGoType: "alias.Flag",
							ExtGoTypeImport: &asyncapi.GoTypeImportExtension{
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
			name: "object",
			schema: &asyncapi.Any{
				Type: "object",
				Properties: map[string]*asyncapi.Any{
					"flag": {
						Type: "integer",
						Extensions: asyncapi.Extensions{
							ExtGoType: "alias.Flag",
							ExtGoTypeImport: &asyncapi.GoTypeImportExtension{
								Path: "abc.xyz/repo/mypackage",
								Name: "alias",
							},
						},
					},
					"id": {
						Type: "string",
						Extensions: asyncapi.Extensions{
							ExtGoType: "xid.ID",
							ExtGoTypeImport: &asyncapi.GoTypeImportExtension{
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
			name: "object",
			schema: &asyncapi.Any{
				Type: "object",
				Properties: map[string]*asyncapi.Any{
					"start_flag": {
						Type: "integer",
						Extensions: asyncapi.Extensions{
							ExtGoType: "alias2.Flag",
							ExtGoTypeImport: &asyncapi.GoTypeImportExtension{
								Path: "abc.xyz/repo/mypackage",
								Name: "alias1",
							},
						},
					},
					"end_flag": {
						Type: "integer",
						Extensions: asyncapi.Extensions{
							ExtGoType: "alias2.Flag",
							ExtGoTypeImport: &asyncapi.GoTypeImportExtension{
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

	is := is.New(t)
	for _, test := range tests {
		spec := asyncapi.Specification{
			Components: asyncapi.Components{
				Schemas: map[string]*asyncapi.Any{test.name: test.schema},
			},
		}

		res, err := spec.CustomImports()

		if test.error {
			is.True(err != nil)
		} else {
			is.NoErr(err)

			sort.Strings(test.expected)
			sort.Strings(res)
			is.Equal(res, test.expected)
		}
	}
}
