package asyncapiv3

import (
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/utils/template"
)

// generateFullName generates a full name for a struct.
// If number is nil, it will not be added to the name.
func generateFullName(parentName, name, typeName string, number *int) string {
	// Namify all the strings
	parentName = template.Namify(parentName)
	name = template.Namify(name)
	typeName = template.Namify(typeName)

	// If number is nil, add number to type
	if number != nil {
		typeName += fmt.Sprintf("_%d", *number)
	}

	// If there is a parent name, prefix it with a "From"
	if parentName != "" {
		parentName = "From_" + parentName
	}

	// Return the name with the number
	return template.Namify(name + "_" + typeName + "_" + parentName)
}
