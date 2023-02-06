package generators

import (
	"bytes"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
)

type ClientSubscriberGenerator struct {
	Specification asyncapi.Specification

	SubscribeCount uint
}

func NewClientSubscriberGenerator(spec asyncapi.Specification) ClientSubscriberGenerator {
	_, subscribeCount := spec.GetPublishSubscribeCount()

	return ClientSubscriberGenerator{
		Specification:  spec,
		SubscribeCount: subscribeCount,
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
