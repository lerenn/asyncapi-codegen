package generatorv3

import asyncapi "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"

// ActionOperations contains operations based on their action.
type ActionOperations struct {
	SendCount    uint
	ReceiveCount uint
	Receive      map[string]*asyncapi.Operation
	Send         map[string]*asyncapi.Operation
}

// NewActionOperations will create a struct with operations based on their action.
func NewActionOperations(side Side, spec asyncapi.Specification) ActionOperations {
	var ao ActionOperations

	// Get action count based on action
	sendCount, receiveCount := spec.GetOperationCountByAction()
	if side == SideIsApplication {
		ao.SendCount = sendCount
		ao.ReceiveCount = receiveCount
	} else {
		ao.SendCount = receiveCount
		ao.ReceiveCount = sendCount
	}

	// Get channels based on send/receive
	ao.Receive = make(map[string]*asyncapi.Operation)
	ao.Send = make(map[string]*asyncapi.Operation)
	for name, op := range spec.Operations {
		// Add channel to receive channels based on operation action and side
		if isControllerReceiveOperation(side, op) {
			ao.Receive[name] = op
		}

		// Add channel to send channels based on operation action and side
		if isControllerSendOperation(side, op) {
			ao.Send[name] = op
		}

		// Add a artificial operation for reply if the controller should respond to a reply
		if shouldControllerRespondToReply(side, op) {
			ch := op.Reply.Channel.Follow()
			ao.Send[ch.Name] = op.ReplyIs
		}
	}

	return ao
}
