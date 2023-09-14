package generators

import (
	"bytes"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
)

// SubscriberGenerator is a code generator for subscribers that will turn an
// asyncapi specification into subscriber golang code.
type SubscriberGenerator struct {
	MethodCount uint
	Channels    map[string]*asyncapi.Channel
	Prefix      string
}

// NewSubscriberGenerator will create a new subscriber code generator.
func NewSubscriberGenerator(side Side, spec asyncapi.Specification) SubscriberGenerator {
	var gen SubscriberGenerator

	// Get subscription methods count based on publish/subscribe count
	publishCount, subscribeCount := spec.GetPublishSubscribeCount()
	if side == SideIsApplication {
		gen.MethodCount = publishCount
	} else {
		gen.MethodCount = subscribeCount
	}

	// Get channels based on publish/subscribe
	gen.Channels = make(map[string]*asyncapi.Channel)
	for k, v := range spec.Channels {
		// Channels are reverse on application side
		if v.Publish != nil && side == SideIsApplication {
			gen.Channels[k] = v
		} else if v.Subscribe != nil && side == SideIsUser {
			gen.Channels[k] = v
		}
	}

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
		schemaTemplatePath,
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
