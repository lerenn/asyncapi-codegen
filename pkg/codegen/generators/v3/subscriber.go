package generatorv3

import (
	"bytes"

	asyncapi "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
)

// SubscriberGenerator is a code generator for subscribers that will turn an
// asyncapi specification into a subscriber golang code.
type SubscriberGenerator struct {
	SendOperationsCount    uint
	ReceiveOperationsCount uint
	Operations             map[string]*asyncapi.Operation
	Prefix                 string
}

// NewSubscriberGenerator will create a new subscriber code generator.
func NewSubscriberGenerator(side Side, spec asyncapi.Specification) SubscriberGenerator {
	var gen SubscriberGenerator

	// Get send/receive methods count based on action count
	sendCount, receiveCount := spec.GetOperationCountByAction()
	if side == SideIsApplication {
		gen.SendOperationsCount = sendCount
		gen.ReceiveOperationsCount = receiveCount
	} else {
		gen.SendOperationsCount = receiveCount
		gen.ReceiveOperationsCount = sendCount
	}

	// Get channels based on send/receive
	gen.Operations = make(map[string]*asyncapi.Operation)
	for k, op := range spec.Operations {
		// Channels are reverse on application side
		if op.Action.IsReceive() && side == SideIsApplication {
			gen.Operations[k] = op
		} else if op.Action.IsSend() && side == SideIsUser {
			gen.Operations[k] = op
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
