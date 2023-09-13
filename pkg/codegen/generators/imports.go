package generators

import (
	"bytes"
)

// ImportsGenerator is a code generator for imports that will add needed imports
// to the code, being asyncapi-codegen packages, standard library packages or
// custom packages.
type ImportsGenerator struct {
	PackageName   string
	ModuleVersion string
	ModuleName    string
	CustomImports []string
}

// Generate will generate the imports code
func (ig ImportsGenerator) Generate() (string, error) {
	tmplt, err := loadTemplate(importsTemplatePath)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := tmplt.Execute(buf, ig); err != nil {
		return "", err
	}

	return buf.String(), nil
}
