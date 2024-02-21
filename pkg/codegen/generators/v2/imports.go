package v2

import (
	"bytes"
)

// importsGenerator is a code generator for imports that will add needed imports
// to the code, being asyncapi-codegen packages, standard library packages or
// custom packages.
type importsGenerator struct {
	PackageName   string
	ModuleVersion string
	ModuleName    string
	CustomImports []string
}

// Generate will generate the imports code.
func (ig importsGenerator) Generate() (string, error) {
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
