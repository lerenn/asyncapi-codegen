package generators

import (
	"bytes"
)

type ImportsGenerator struct {
	PackageName   string
	ModuleVersion string
	ModuleName    string
	CustomImports []string
}

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
