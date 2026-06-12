package generators

import (
	"fmt"
	"strconv"
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

	if shouldAddOmitEmpty(schema) {
		directives = append(directives, "omitempty")
	}

	return fmt.Sprintf("json:\"%s\"", strings.Join(directives, ","))
}

func shouldAddOmitEmpty[T any](schema asyncapi.Validations[T]) bool {
	if schema.IsRequired {
		return schema.ShouldOmitEmpty != nil && *schema.ShouldOmitEmpty
	}

	return schema.ShouldOmitEmpty == nil || *schema.ShouldOmitEmpty
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
		if cStr, ok := scalarToValidatorValue(schema.Const); ok {
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
	if len(schema.Enum) == 0 {
		return directives
	}

	enumsStr := make([]string, 0, len(schema.Enum))
	for _, e := range schema.Enum {
		v, ok := scalarToValidatorValue(e)
		if !ok {
			// If any element is not a supported scalar, the whole enum is unsupported.
			return directives
		}

		// String values are single-quoted to handle values with spaces.
		if _, isStr := e.(string); isStr {
			v = singleQuote(v)
		}
		enumsStr = append(enumsStr, v)
	}

	return append(directives, fmt.Sprintf("oneof=%s", strings.Join(enumsStr, " ")))
}

// scalarToValidatorValue returns the validator literal for a scalar enum/const
// value. It reports false when the value is not a supported scalar (object,
// array, ...), in which case no validation directive should be generated.
func scalarToValidatorValue(v any) (string, bool) {
	switch t := v.(type) {
	case string:
		return t, true
	case bool:
		return strconv.FormatBool(t), true
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), true
	case int:
		return strconv.Itoa(t), true
	case int64:
		return strconv.FormatInt(t, 10), true
	default:
		return "", false
	}
}
