package generators

import (
	"bytes"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
)

type AppControllerGenerator struct {
	Specification asyncapi.Specification

	// CorrelationIDLocation will indicate where the correlation id is
	// According to this: https://www.asyncapi.com/docs/reference/specification/v2.5.0#correlationIDObject
	CorrelationIDLocation map[string]string

	PublishCount   uint
	SubscribeCount uint
}

func NewAppControllerGenerator(spec asyncapi.Specification) AppControllerGenerator {
	publishCount, subscribeCount := spec.GetPublishSubscribeCount()

	return AppControllerGenerator{
		Specification: spec,

		CorrelationIDLocation: getCorrelationIDsLocationsByChannel(spec),

		PublishCount:   publishCount,
		SubscribeCount: subscribeCount,
	}
}

func (acg AppControllerGenerator) Generate() (string, error) {
	tmplt, err := loadTemplate(
		appControllerTemplatePath,
		anyTemplatePath,
		messageTemplatePath,
	)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := tmplt.Execute(buf, acg); err != nil {
		return "", err
	}

	return buf.String(), nil
}
