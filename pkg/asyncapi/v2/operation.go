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

// generateMetadata generate metadata for the operation and its children.
func (op *Operation) generateMetadata(name string) error {
	op.Name = template.Namify(name)

	// Get message name
	var msgName string
	if op.Message.Reference != "" {
		msgName = strings.Split(op.Message.Reference, "/")[3]
	} else {
		msgName = op.Name
	}

	// Generate message metadata
	return op.Message.generateMetadata(msgName + MessageSuffix)
}

// setDependencies set dependencies for the operation and its children from specification.
func (op *Operation) setDependencies(spec Specification) error {
	// Set message dependencies
	return op.Message.setDependencies(spec)
}
