package templates

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"

	asyncapi "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
	"github.com/lerenn/asyncapi-codegen/pkg/utils"
	templateutil "github.com/lerenn/asyncapi-codegen/pkg/utils/template"
)

// GetChildrenObjectSchemas will return all the children object schemas of a
// schema, only from first level and without AnyOf, AllOf and OneOf.
func GetChildrenObjectSchemas(s asyncapi.Schema) []*asyncapi.Schema {
	allSchemas := utils.MapToList(s.Properties)

	if s.Items != nil {
		allSchemas = append(allSchemas, s.Items)
	}

	if s.AdditionalProperties != nil {
		allSchemas = append(allSchemas, s.AdditionalProperties)
	}

	// Only keep object schemas
	filteredSchemas := make([]*asyncapi.Schema, 0, len(allSchemas))
	for _, schema := range allSchemas {
		if schema.Type == asyncapi.SchemaTypeIsObject.String() {
			filteredSchemas = append(filteredSchemas, schema)
		} else if schema.Type == asyncapi.SchemaTypeIsArray.String() &&
			schema.Items != nil &&
			schema.Items.Type == asyncapi.SchemaTypeIsObject.String() {
			filteredSchemas = append(filteredSchemas, schema.Items)
		}
	}

	return filteredSchemas
}

// referenceToSlicePath will convert a reference to a slice where each element is a
// step of the path.
func referenceToSlicePath(ref string) []string {
	ref = strings.ReplaceAll(ref, ".", "/")
	ref = strings.ReplaceAll(ref, "#", "")
	return strings.Split(ref, "/")[1:]
}

// ReferenceToStructAttributePath will convert a reference to a struct attribute
// path in the form of "a.b.c" where a, b and c are struct attributes in the
// form of golang conventional type names.
func ReferenceToStructAttributePath(ref string) string {
	path := referenceToSlicePath(ref)

	for k, v := range path {
		// If this is concerning the header, then it will be named "headers"
		if v == asyncapi.MessageFieldIsHeader.String() {
			v = "headers"
		}

		path[k] = templateutil.Namify(v)
	}

	return strings.Join(path, ".")
}

// ChannelToMessageTypeName will convert a channel to a message type name in the
// form of golang conventional type names.
func ChannelToMessageTypeName(ch asyncapi.Channel) string {
	msg, err := ch.Follow().GetMessage()
	if err != nil {
		panic(err)
	}
	return templateutil.Namify(msg.Follow().Name)
}

// OpToMsgTypeName will convert an operation to a message type name in the
// form of golang conventional type names.
func OpToMsgTypeName(op asyncapi.Operation) string {
	msg, err := op.Follow().GetMessage()
	if err != nil {
		panic(err)
	}
	return templateutil.Namify(msg.Follow().Name)
}

// OpToChannelTypeName will convert an operation to a channel type name in the
// form of golang conventional type names.
func OpToChannelTypeName(op asyncapi.Operation) string {
	ch := op.Channel.Follow()
	return templateutil.Namify(ch.Name)
}

// IsRequired will check if a field is required in a asyncapi struct.
func IsRequired(schema asyncapi.Schema, field string) bool {
	return schema.IsFieldRequired(field)
}

// GenerateChannelAddrFromOp will generate a channel path with the given operation.
func GenerateChannelAddrFromOp(op asyncapi.Operation) string {
	ch := op.Channel.Follow()
	return GenerateChannelAddr(ch)
}

// GenerateChannelAddr will generate a channel path with the given channel.
func GenerateChannelAddr(ch *asyncapi.Channel) string {
	// Be sure this is the final channel, not a proxy
	ch = ch.Follow()

	// If there is no parameter, then just return the path
	if ch.Parameters == nil {
		return fmt.Sprintf("%q", ch.Address)
	}

	parameterRegexp := regexp.MustCompile("{[^{}]*}")

	matches := parameterRegexp.FindAllString(ch.Address, -1)
	format := parameterRegexp.ReplaceAllString(ch.Address, "%s")

	sprint := fmt.Sprintf("fmt.Sprintf(%q, ", format)
	for _, m := range matches {
		sprint += fmt.Sprintf("params.%s,", templateutil.Namify(m))
	}

	return sprint[:len(sprint)-1] + ")"
}

func appendDirectiveIfDefined(directives []string, tagName string, value float64) []string {
	if value != 0 {
		return append(directives, fmt.Sprintf("%s=%g", tagName, value))
	}
	return directives
}

// GenerateValidateTags returns the "validate" tag for a given field in a struct, based on the asyncapi contract.
// This tag can then be used by go-playground/validator/v10 to validate the struct's content.
func GenerateValidateTags(schema asyncapi.Schema) string {
	var directives []string
	if schema.IsRequired {
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
	if len(schema.Enum) > 0 {
		var enumsStr []string
		for _, e := range schema.Enum {
			if eStr, ok := e.(string); ok {
				enumsStr = append(enumsStr, eStr)
			}
		}

		// Only generate enum if all elements are string, otherwise this is unsupported
		if len(schema.Enum) == len(enumsStr) {
			directives = append(directives, fmt.Sprintf("oneof=%s", strings.Join(enumsStr, " ")))
		}
	}
	if schema.Const != nil {
		if cStr, ok := schema.Const.(string); ok {
			// Only generate enum if the elements is a string, otherwise this is unsupported
			directives = append(directives, fmt.Sprintf("eq=%s", cStr))
		}
	}

	if len(directives) > 0 {
		return fmt.Sprintf(" validate:\"%s\"", strings.Join(directives, ","))
	} else {
		return ""
	}
}

// HelpersFunctions returns the functions that can be used as helpers
// in a golang template.
func HelpersFunctions() template.FuncMap {
	return template.FuncMap{
		"getChildrenObjectSchemas":       GetChildrenObjectSchemas,
		"channelToMessageTypeName":       ChannelToMessageTypeName,
		"opToMsgTypeName":                OpToMsgTypeName,
		"opToChannelTypeName":            OpToChannelTypeName,
		"isRequired":                     IsRequired,
		"generateChannelAddr":            GenerateChannelAddr,
		"generateChannelAddrFromOp":      GenerateChannelAddrFromOp,
		"referenceToStructAttributePath": ReferenceToStructAttributePath,
		"generateValidateTags":           GenerateValidateTags,
	}
}
