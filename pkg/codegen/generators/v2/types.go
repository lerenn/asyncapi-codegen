package generatorv2

import (
	"bytes"

	asyncapi "github.com/TheSadlig/asyncapi-codegen/pkg/asyncapi/v2"
)

// TypesGenerator is a code generator for types that will generate all schemas
// contained in an asyncapi specification to golang structures code.
type TypesGenerator struct {
	asyncapi.Specification
	IgnoreStringFormat bool
}

// Generate will create a new types code generator.
func (tg TypesGenerator) Generate() (string, error) {
	tmplt, err := loadTemplate(
		typesTemplatePath,
		schemaDefinitionTemplatePath,
		schemaNameTemplatePath,
		messageTemplatePath,
		parameterTemplatePath,

		marshalingAdditionalPropertiesTemplatePath,
		marshalingTimeTemplatePath,
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
