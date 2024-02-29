package generatorv3

import (
	"bytes"

	asyncapi "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
)

// ControllerGenerator is a code generator for controllers that will turn an
// asyncapi specification into controller golang code.
type ControllerGenerator struct {
	SendOperationsCount    uint
	ReceiveOperationsCount uint
	ReceiveOperations      map[string]*asyncapi.Operation
	SendOperations         map[string]*asyncapi.Operation
	Prefix                 string
	Version                string
}

// NewControllerGenerator will create a new controller code generator.
func NewControllerGenerator(side Side, spec asyncapi.Specification) ControllerGenerator {
	var gen ControllerGenerator

	// Get action count based on action
	sendCount, receiveCount := spec.GetOperationCountByAction()
	if side == SideIsApplication {
		gen.SendOperationsCount = sendCount
		gen.ReceiveOperationsCount = receiveCount
	} else {
		gen.SendOperationsCount = receiveCount
		gen.ReceiveOperationsCount = sendCount
	}

	// Get channels based on send/receive
	gen.ReceiveOperations = make(map[string]*asyncapi.Operation)
	gen.SendOperations = make(map[string]*asyncapi.Operation)
	for name, op := range spec.Operations {
		// Add channel to receive channels based on operation action and side
		if isControllerReceiveOperation(side, op) {
			gen.ReceiveOperations[name] = op
		}

		// Add channel to send channels based on operation action and side
		if isControllerSendOperation(side, op) {
			gen.SendOperations[name] = op
		}

		// Add a artificial operation for reply if the controller should respond to a reply
		if shouldControllerRespondToReply(side, op) {
			ch := op.Reply.Channel.Follow()
			gen.SendOperations[ch.Name] = &asyncapi.Operation{
				Name:      "ReplyTo" + op.Name,
				Channel:   ch,
				IsReplyTo: op,
			}
		}
	}

	// Set generation name
	if side == SideIsApplication {
		gen.Prefix = "App"
	} else {
		gen.Prefix = "User"
	}

	// Set version
	gen.Version = spec.Info.Version

	return gen
}

func shouldControllerRespondToReply(side Side, op *asyncapi.Operation) bool {
	if op.Reply == nil || op.Reply.Channel == nil {
		return false
	}

	switch {
	case side == SideIsApplication && op.Action.IsReceive():
		return true
	case side == SideIsUser && op.Action.IsSend():
		return true
	default:
		return false
	}
}

func isControllerReceiveOperation(side Side, op *asyncapi.Operation) bool {
	switch {
	case side == SideIsApplication && op.Action.IsReceive():
		return true
	case side == SideIsUser && op.Action.IsSend():
		return true
	default:
		return false
	}
}

func isControllerSendOperation(side Side, op *asyncapi.Operation) bool {
	switch {
	case side == SideIsApplication && op.Action.IsSend():
		return true
	case side == SideIsUser && op.Action.IsReceive():
		return true
	default:
		return false
	}
}

// Generate will generate the controller code.
func (asg ControllerGenerator) Generate() (string, error) {
	tmplt, err := loadTemplate(
		controllerTemplatePath,
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
