package asyncapiv3

import (
	"fmt"
	"strings"
)

// Extensions holds additional properties defined for asyncapi-codegen
// that are out of the AsyncAPI spec.
type Extensions struct {
	// Setting custom Go type when generating schemas
	ExtGoType string `json:"x-go-type,omitempty"`

	// Setting custom import statements for ExtGoType
	ExtGoTypeImport *GoTypeImportExtension `json:"x-go-type-import,omitempty"`

	// Controls whether to include omitempty in JSON tags
	// If false, omitempty will be removed from JSON tags even if the field can be null
	ExtOmitEmpty *bool `json:"x-omitempty,omitempty"`
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
func (s *Schema) goTypeImports(imports map[GoTypeImportPath]GoTypeImportName) error {
	// Process Properties
	for _, p := range s.Properties {
		if err := p.goTypeImports(imports); err != nil {
			return err
		}
	}

	// Process Items
	if s.Items != nil {
		if err := s.Items.goTypeImports(imports); err != nil {
			return err
		}
	}

	if s.ExtGoTypeImport != nil {
		name, exists := imports[s.ExtGoTypeImport.Path]
		if exists && name != s.ExtGoTypeImport.Name {
			return fmt.Errorf(
				"x-go-type-import name conflict for item %s: %s and %s for %s",
				s.Name, name, s.ExtGoTypeImport.Name, s.ExtGoTypeImport.Path,
			)
		}

		if !exists {
			imports[s.ExtGoTypeImport.Path] = s.ExtGoTypeImport.Name
		}
	}
	return nil
}

// CustomImports collects all custom import paths set by x-go-type-imports
// in all Schema Objects in the Specification.
// Returns import strings like `alias "abc.xyz/repo/package"` for code generation.
// Returns error when import name conflicts.
func (s Specification) CustomImports() ([]string, error) {
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

	// TODO: support Parameters

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
