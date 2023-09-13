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
	gen.SubscribeChannels = make(map[string]*asyncapi.Channel)
	gen.PublishChannels = make(map[string]*asyncapi.Channel)
	for name, channel := range spec.Channels {
		if isSubscribeChannel(side, channel) {
			gen.SubscribeChannels[name] = channel
		} else {
			gen.PublishChannels[name] = channel
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

func isSubscribeChannel(side Side, channel *asyncapi.Channel) bool {
	switch {
	case side == SideIsApplication && channel.Publish != nil:
		return true
	case side == SideIsUser && channel.Subscribe != nil:
		return true
	case side == SideIsApplication && channel.Subscribe != nil:
		return false
	case side == SideIsUser && channel.Publish != nil:
		return false
	default:
		panic("this should never happen")
	}
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
