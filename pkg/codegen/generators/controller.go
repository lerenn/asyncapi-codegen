package generators

import (
	"bytes"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
)

type ControllerGenerator struct {
	MethodCount       uint
	SubscribeChannels map[string]asyncapi.Channel
	PublishChannels   map[string]asyncapi.Channel
	Prefix            string
}

func NewControllerGenerator(side Side, spec asyncapi.Specification) ControllerGenerator {
	var gen ControllerGenerator

	// Get subscription methods count based on publish/subscribe count
	publishCount, subscribeCount := spec.GetPublishSubscribeCount()
	if side == SideIsApplication {
		gen.MethodCount = publishCount
	} else {
		gen.MethodCount = subscribeCount
	}

	// Get channels based on publish/subscribe
	gen.SubscribeChannels = make(map[string]asyncapi.Channel)
	gen.PublishChannels = make(map[string]asyncapi.Channel)
	for k, v := range spec.Channels {
		// Channels are reverse on application side
		if side == SideIsApplication {
			if v.Publish != nil {
				gen.SubscribeChannels[k] = v
			} else if v.Subscribe != nil {
				gen.PublishChannels[k] = v
			}
		} else {
			if v.Publish != nil {
				gen.PublishChannels[k] = v
			} else if v.Subscribe != nil {
				gen.SubscribeChannels[k] = v
			}
		}
	}

	// Set generation name
	if side == SideIsApplication {
		gen.Prefix = "App"
	} else {
		gen.Prefix = "Client"
	}

	return gen
}

func (asg ControllerGenerator) Generate() (string, error) {
	tmplt, err := loadTemplate(
		controllerTemplatePath,
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
