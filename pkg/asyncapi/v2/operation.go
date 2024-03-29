package asyncapiv2

import (
	"strings"

	"github.com/lerenn/asyncapi-codegen/pkg/utils/template"
)

// Operation is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v2.6.0#operationObject
type Operation struct {
	// --- AsyncAPI fields -----------------------------------------------------

	OperationID string  `json:"operationId"`
	Message     Message `json:"message"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name string `json:"-"`
}

// Process processes the Operation structure to make it ready for code generation.
func (op *Operation) Process(name string, spec Specification) error {
	op.Name = template.Namify(name)

	// Get message name
	var msgName string
	if op.Message.Reference != "" {
		msgName = strings.Split(op.Message.Reference, "/")[3]
	} else {
		msgName = op.Name
	}

	// Process message
	return op.Message.Process(msgName+MessageSuffix, spec)
}
