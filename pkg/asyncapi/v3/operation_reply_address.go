package asyncapiv3

import (
	"strings"

	"github.com/lerenn/asyncapi-codegen/pkg/utils"
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
	LocationTo       *Schema                `json:"-"`
	LocationRequired bool                   `json:"-"`
}

// Process processes the OperationReplyAddress to make it ready for code generation.
func (ora *OperationReplyAddress) Process(name string, op *Operation, spec Specification) {
	// Prevent modification if nil
	if ora == nil {
		return
	}

	// Set name
	ora.Name = utils.UpperFirstLetter(name)

	// Add pointer to reference if there is one
	if ora.Reference != "" {
		ora.ReferenceTo = spec.ReferenceOperationReplyAddress(ora.Reference)
	}

	// Get location to schema
	ora.LocationTo = spec.ReferenceSchema(ora.Location)
	ora.LocationRequired = ora.isLocationRequired(op)
}

func (ora OperationReplyAddress) isLocationRequired(op *Operation) bool {
	if ora.Location == "" {
		return false
	}

	locationParent := op.Follow().GetMessage().createTreeUntilLocation(ora.Location)
	path := strings.Split(ora.Location, "/")
	return locationParent.IsFieldRequired(path[len(path)-1])
}
