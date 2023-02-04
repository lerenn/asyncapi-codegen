package generators

import (
	"bytes"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
)

type ClientControllerGenerator struct {
	Specification asyncapi.Specification

	// CorrelationIDLocation will indicate where the correlation id is
	// According to this: https://www.asyncapi.com/docs/reference/specification/v2.5.0#correlationIDObject
	CorrelationIDLocation map[string]string
}

func NewClientControllerGenerator(spec asyncapi.Specification) ClientControllerGenerator {
	return ClientControllerGenerator{
		Specification:         spec,
		CorrelationIDLocation: getCorrelationIDsLocationsByChannel(spec),
	}
}

func (ccg ClientControllerGenerator) Generate() (string, error) {
	tmplt, err := loadTemplate(
		clientControllerTemplatePath,
		anyTemplatePath,
		messageTemplatePath,
	)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := tmplt.Execute(buf, ccg); err != nil {
		return "", err
	}

	return buf.String(), nil
}
