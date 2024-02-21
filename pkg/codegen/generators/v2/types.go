package v2

import (
	"bytes"

	asyncapi "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v2"
)

// typesGenerator is a code generator for types that will generate all schemas
// contained in an asyncapi specification to golang structures code.
type typesGenerator struct {
	asyncapi.Specification
}

// Generate will create a new types code generator.
func (tg typesGenerator) Generate() (string, error) {
	tmplt, err := loadTemplate(
		typesTemplatePath,
		schemaTemplatePath,
		messageTemplatePath,
		parameterTemplatePath,
	)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := tmplt.Execute(buf, tg); err != nil {
		return "", err
	}

	return buf.String(), nil
}
