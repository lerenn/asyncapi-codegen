package asyncapiv3

import (
	"strings"

	"github.com/lerenn/asyncapi-codegen/pkg/utils/template"
)

// OperationReplyAddress is a representation of the corresponding asyncapi object
// filled from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#operationReplyAddressObject
type OperationReplyAddress struct {
	// --- AsyncAPI fields -----------------------------------------------------

	Description string `json:"description"`
	Location    string `json:"location"`
	Reference   string `json:"$ref"`

	// --- Non AsyncAPI fields -------------------------------------------------

	Name             string                 `json:"-"`
	ReferenceTo      *OperationReplyAddress `json:"-"`
	LocationRequired bool                   `json:"-"`
}

// Process processes the OperationReplyAddress to make it ready for code generation.
func (ora *OperationReplyAddress) Process(name string, op *Operation, spec Specification) error {
	// Prevent modification if nil
	if ora == nil {
		return nil
	}

	// Set name
	ora.Name = template.Namify(name)

	// Add pointer to reference if there is one
	if ora.Reference != "" {
		refTo, err := spec.ReferenceOperationReplyAddress(ora.Reference)
		if err != nil {
			return err
		}
		ora.ReferenceTo = refTo
	}

	// Get location to schema
	locRequired, err := ora.isLocationRequired(op)
	if err != nil {
		return err
	}
	ora.LocationRequired = locRequired

	return nil
}

func (ora OperationReplyAddress) isLocationRequired(op *Operation) (bool, error) {
	if ora.Location == "" {
		return false, nil
	}

	msg, err := op.Follow().GetMessage()
	if err != nil {
		return false, err
	}

	locationParent := msg.createTreeUntilLocation(ora.Location)
	path := strings.Split(ora.Location, "/")
	return locationParent.IsFieldRequired(path[len(path)-1]), nil
}
