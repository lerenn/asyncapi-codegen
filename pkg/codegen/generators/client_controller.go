package generators

import (
	"bytes"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
)

type ClientControllerGenerator struct {
	Specification asyncapi.Specification

	PublishCount   uint
	SubscribeCount uint
}

func NewClientControllerGenerator(spec asyncapi.Specification) ClientControllerGenerator {
	publishCount, subscribeCount := spec.GetPublishSubscribeCount()

	return ClientControllerGenerator{
		Specification: spec,

		PublishCount:   publishCount,
		SubscribeCount: subscribeCount,
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
