package generators

import (
	"fmt"
	"strings"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
	"github.com/lerenn/asyncapi-codegen/pkg/utils/template"
)

func appendDirectiveIfDefined(directives []string, tagName string, value float64) []string {
	if value != 0 {
		return append(directives, fmt.Sprintf("%s=%g", tagName, value))
	}
	return directives
}

// GenerateJSONTags returns the "json" tag for a given field in a struct, based on the asyncapi contract.
func GenerateJSONTags[T any](schema asyncapi.Validations[T], field string) string {
	directives := []string{
		template.ConvertKey(field),
	}

	if !schema.IsRequired {
		directives = append(directives, "omitempty")
	}

	return fmt.Sprintf("json:\"%s\"", strings.Join(directives, ","))
}

// GenerateValidateTags returns the "validate" tag for a given field in a struct, based on the asyncapi contract.
// This tag can then be used by go-playground/validator/v10 to validate the struct's content.
func GenerateValidateTags[T any](schema asyncapi.Validations[T], isPointer bool, schemaType string) string {
	var directives []string
	if schema.IsRequired && (isPointer || schemaType == "array") {
		directives = append(directives, "required")
	}

	directives = appendDirectiveIfDefined(directives, "min", float64(schema.MinLength))
	directives = appendDirectiveIfDefined(directives, "max", float64(schema.MaxLength))
	directives = appendDirectiveIfDefined(directives, "gte", schema.Minimum)
	directives = appendDirectiveIfDefined(directives, "lte", schema.Maximum)
	directives = appendDirectiveIfDefined(directives, "gt", schema.ExclusiveMinimum)
	directives = appendDirectiveIfDefined(directives, "lt", schema.ExclusiveMaximum)

	if schema.UniqueItems {
		directives = append(directives, "unique")
	}

	directives = appendEnumDirectives(schema, directives)
	if schema.Const != nil {
		if cStr, ok := schema.Const.(string); ok {
			// Only generate enum if the elements is a string, otherwise this is unsupported
			directives = append(directives, fmt.Sprintf("eq=%s", cStr))
		}
	}

	if len(directives) > 0 {
		if !schema.IsRequired {
			directives = append([]string{"omitempty"}, directives...)
		}

		return fmt.Sprintf(" validate:\"%s\"", strings.Join(directives, ","))
	} else {
		return ""
	}
}

// singleQuote prepends and appends a single quote to the provided string.
func singleQuote(s string) string {
	return "'" + s + "'"
}

func appendEnumDirectives[T any](schema asyncapi.Validations[T], directives []string) []string {
	if len(schema.Enum) > 0 {
		var enumsStr []string
		for _, e := range schema.Enum {
			if eStr, ok := e.(string); ok {
				// single quotes are mandatory in order to handle values with spaces
				enumsStr = append(enumsStr, singleQuote(eStr))
			}
		}

		// Only generate enum if all elements are string, otherwise this is unsupported
		if len(schema.Enum) == len(enumsStr) {
			directives = append(directives, fmt.Sprintf("oneof=%s", strings.Join(enumsStr, " ")))
		}
	}
	return directives
}
