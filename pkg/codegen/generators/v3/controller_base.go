package generatorv3

import (
	"bytes"
)

// ControllerBaseGenerator generates the shared controller boilerplate (the
// internal controller struct and its options) used by both the application and
// user controllers. It is generated once, only when a controller side is
// generated, so that the `types` output stays free of controller boilerplate.
type ControllerBaseGenerator struct{}

// Generate will generate the shared controller base code.
func (ControllerBaseGenerator) Generate() (string, error) {
	tmplt, err := loadTemplate(controllerBaseTemplatePath)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := tmplt.Execute(buf, nil); err != nil {
		return "", err
	}

	return buf.String(), nil
}
