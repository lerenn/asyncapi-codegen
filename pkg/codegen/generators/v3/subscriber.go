package generatorv3

import (
	"bytes"

	asyncapi "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
)

// SubscriberGenerator is a code generator for subscribers that will turn an
// asyncapi specification into a subscriber golang code.
type SubscriberGenerator struct {
	Operations ActionOperations
	Prefix     string
}

// NewSubscriberGenerator will create a new subscriber code generator.
func NewSubscriberGenerator(side Side, spec asyncapi.Specification) SubscriberGenerator {
	var gen SubscriberGenerator

	// Generate receive send operations
	gen.Operations = NewActionOperations(side, spec)

	// Set generation name
	if side == SideIsApplication {
		gen.Prefix = "App"
	} else {
		gen.Prefix = "User"
	}

	return gen
}

// Generate will generate the subscriber code.
func (asg SubscriberGenerator) Generate() (string, error) {
	tmplt, err := loadTemplate(
		subscriberTemplatePath,
		schemaDefinitionTemplatePath,
		schemaNameTemplatePath,
		messageTemplatePath,
	)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := tmplt.Execute(buf, asg); err != nil {
		return "", err
	}

	return buf.String(), nil
}
