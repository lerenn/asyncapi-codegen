package generatorv3

import (
	"bytes"

	asyncapi "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
)

// ListenerGenerator is a code generator for listeners that will turn an
// asyncapi specification into a listener golang code.
type ListenerGenerator struct {
	MethodCount uint
	Operations  map[string]*asyncapi.Operation
	Prefix      string
}

// NewListenerGenerator will create a new listener code generator.
func NewListenerGenerator(side Side, spec asyncapi.Specification) ListenerGenerator {
	var gen ListenerGenerator

	// Get send/receive methods count based on action count
	sendCount, receiveCount := spec.GetByActionCount()
	if side == SideIsApplication {
		gen.MethodCount = sendCount
	} else {
		gen.MethodCount = receiveCount
	}

	// Get channels based on send/receive
	gen.Operations = make(map[string]*asyncapi.Operation)
	for k, op := range spec.Operations {
		// Channels are reverse on application side
		if op.Action.IsSend() && side == SideIsApplication {
			gen.Operations[k] = op
		} else if op.Action.IsReceive() && side == SideIsUser {
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

// Generate will generate the listener code.
func (asg ListenerGenerator) Generate() (string, error) {
	tmplt, err := loadTemplate(
		listenerTemplatePath,
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
