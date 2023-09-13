package asyncapi

import (
	"fmt"
	"strings"
)

// Extensions holds additional properties defined for asyncapi-codegen
// that are out of the AsyncAPI spec.
type Extensions struct {
	// Setting custom Go type when generating schemas
	ExtGoType string `json:"x-go-type"`

	// Setting custom import statements for ExtGoType
	ExtGoTypeImport *GoTypeImportExtension `json:"x-go-type-import"`
}

// GoTypeImportExtension specifies the required import statement
// for the x-go-type extension.
// For example, GoTypeImportExtension{Name: "myuuid", Path: "github.com/google/uuid"}
// will generate `import myuuid github.com/google/uuid`.
type GoTypeImportExtension struct {
	Name GoTypeImportName `json:"name"` // Package name for import, optional
	Path GoTypeImportPath `json:"path"` // Path to package to import
}

// GoTypeImportPath is the import path type for x-go-type-import.
type GoTypeImportPath string

// GoTypeImportName is the import name type for x-go-type-import.
type GoTypeImportName string

// goTypeImports collects custom imports in this Schema Object set by x-go-type-import key
// into the imports map.
// Reports error when the same import path is assigned multiple import names.
func (a *Any) goTypeImports(imports map[GoTypeImportPath]GoTypeImportName) error {
	// Process Properties
	for _, p := range a.Properties {
		if err := p.goTypeImports(imports); err != nil {
			return err
		}
	}

	// Process Items
	if a.Items != nil {
		if err := a.Items.goTypeImports(imports); err != nil {
			return err
		}
	}

	if a.ExtGoTypeImport != nil {
		name, exists := imports[a.ExtGoTypeImport.Path]
		if exists && name != a.ExtGoTypeImport.Name {
			return fmt.Errorf(
				"x-go-type-import name conflict for item %s: %s and %s for %s",
				a.Name, name, a.ExtGoTypeImport.Name, a.ExtGoTypeImport.Path,
			)
		}

		if !exists {
			imports[a.ExtGoTypeImport.Path] = a.ExtGoTypeImport.Name
		}
	}
	return nil
}

// CustomImports collects all custom import paths set by x-go-type-imports
// in all Schema Objects in the Specification.
// Returns import strings like `alias "abc.xyz/repo/package"` for code generation.
// Returns error when import name conflicts.
func (s Specification) CustomImports() ([]string, error) { //nolint:cyclop
	importsSet := make(map[GoTypeImportPath]GoTypeImportName)

	for _, v := range s.Components.Schemas {
		if err := v.goTypeImports(importsSet); err != nil {
			return nil, fmt.Errorf("/components/schemas custom import error: %w", err)
		}
	}

	for _, v := range s.Components.Messages {
		if v.Payload != nil {
			if err := v.Payload.goTypeImports(importsSet); err != nil {
				return nil, fmt.Errorf("/components/messages payload custom import error: %w", err)
			}
		}

		if v.Headers != nil {
			if err := v.Headers.goTypeImports(importsSet); err != nil {
				return nil, fmt.Errorf("/components/messages headers custom import error: %w", err)
			}
		}
	}

	for _, v := range s.Components.Parameters {
		if v.Schema != nil {
			if err := v.Schema.goTypeImports(importsSet); err != nil {
				return nil, fmt.Errorf("/components/parameters custom import error: %w", err)
			}
		}
	}

	return importsMapToList(importsSet), nil
}

func importsMapToList(importsSet map[GoTypeImportPath]GoTypeImportName) []string {
	imports := make([]string, 0, len(importsSet))
	for path, name := range importsSet {
		if path == "" {
			continue
		}
		imports = append(imports, strings.TrimSpace(fmt.Sprintf(`%s "%s"`, name, path)))
	}
	return imports
}
