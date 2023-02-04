package generators

import (
	"bytes"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
)

type ClientSubscriberGenerator struct {
	Specification asyncapi.Specification

	// CorrelationIDLocation will indicate where the correlation id is
	// According to this: https://www.asyncapi.com/docs/reference/specification/v2.5.0#correlationIDObject
	CorrelationIDLocation map[string]string
}

func NewClientSubscriberGenerator(spec asyncapi.Specification) ClientSubscriberGenerator {
	return ClientSubscriberGenerator{
		Specification:         spec,
		CorrelationIDLocation: getCorrelationIDsLocationsByChannel(spec),
	}
}

func (asg ClientSubscriberGenerator) Generate() (string, error) {
	tmplt, err := loadTemplate(
		clientSubscriberTemplatePath,
		anyTemplatePath,
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
