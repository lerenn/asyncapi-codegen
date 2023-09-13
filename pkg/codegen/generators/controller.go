package generators

import (
	"bytes"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
)

// ControllerGenerator is a code generator for controllers that will turn an
// asyncapi specification into controller golang code.
type ControllerGenerator struct {
	MethodCount       uint
	SubscribeChannels map[string]*asyncapi.Channel
	PublishChannels   map[string]*asyncapi.Channel
	Prefix            string
}

// NewControllerGenerator will create a new controller code generator.
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
	gen.SubscribeChannels, gen.PublishChannels = getChannelsBasedOnSide(side, spec)

	// Set generation name
	if side == SideIsApplication {
		gen.Prefix = "App"
	} else {
		gen.Prefix = "User"
	}

	return gen
}

func getChannelsBasedOnSide(side Side, spec asyncapi.Specification) (
	subscribeChannels map[string]*asyncapi.Channel,
	publishChannels map[string]*asyncapi.Channel,
) {
	subscribeChannels = make(map[string]*asyncapi.Channel)
	publishChannels = make(map[string]*asyncapi.Channel)
	for k, v := range spec.Channels {
		// Channels are reverse on application side
		if (side == SideIsApplication && v.Publish != nil) || (side == SideIsUser && v.Subscribe != nil) {
			subscribeChannels[k] = v
		} else if (side == SideIsApplication && v.Subscribe != nil) || (side == SideIsUser && v.Publish != nil) {
			publishChannels[k] = v
		} else {
			panic("this should never happen")
		}
	}
	return subscribeChannels, publishChannels
}

// Generate will generate the controller code.
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
