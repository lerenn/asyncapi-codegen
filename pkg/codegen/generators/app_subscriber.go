package generators

import (
	"bytes"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
)

type AppSubscriberGenerator struct {
	Specification asyncapi.Specification

	PublishCount uint
}

func NewAppSubscriberGenerator(spec asyncapi.Specification) AppSubscriberGenerator {
	publishCount, _ := spec.GetPublishSubscribeCount()

	return AppSubscriberGenerator{
		Specification: spec,
		PublishCount:  publishCount,
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
