package generators

import (
	"bytes"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
)

type AppSubscriberGenerator struct {
	Specification asyncapi.Specification

	// CorrelationIdLocation will indicate where the correlation id is
	// According to this: https://www.asyncapi.com/docs/reference/specification/v2.5.0#correlationIdObject
	CorrelationIdLocation map[string]string
}

func NewAppSubscriberGenerator(spec asyncapi.Specification) AppSubscriberGenerator {
	return AppSubscriberGenerator{
		Specification:         spec,
		CorrelationIdLocation: getCorrelationIdsLocationsByChannel(spec),
	}
}

func (asg AppSubscriberGenerator) Generate() (string, error) {
	tmplt, err := loadTemplate(
		appSubscriberTemplatePath,
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
