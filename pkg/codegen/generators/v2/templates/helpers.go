package templates

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"

	asyncapi "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v2"
	"github.com/lerenn/asyncapi-codegen/pkg/codegen/generators"
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

// ReferenceToTypeName will convert a reference to a type name in the form of
// golang conventional type names.
func ReferenceToTypeName(ref string) string {
	parts := strings.Split(ref, "/")
	return templateutil.Namify(parts[3])
}

// ChannelToMessage will convert a channel to its message, based on publish/subscribe.
//
//nolint:cyclop
func ChannelToMessage(ch asyncapi.Channel, direction string) *asyncapi.Message {
	switch {
	case ch.Publish != nil && ch.Subscribe == nil:
		return ch.Publish.Message.Follow()
	case ch.Subscribe != nil && ch.Publish == nil:
		return ch.Subscribe.Message.Follow()
	case direction == "publish":
		if ch.Publish == nil {
			panic("ChannelToMessage: channel has no publish operation")
		}
		return ch.Publish.Message.Follow()
	case direction == "subscribe":
		if ch.Subscribe == nil {
			panic("ChannelToMessage: channel has no subscribe operation")
		}
		return ch.Subscribe.Message.Follow()
	case ch.Subscribe == nil && ch.Publish == nil:
		panic("ChannelToMessage: channel has no publish or subscribe operation")
	default:
		panic("direction must be either 'publish' or 'subscribe'")
	}
}

// IsRequired will check if a field is required in a asyncapi struct.
func IsRequired(schema asyncapi.Schema, field string) bool {
	return schema.IsFieldRequired(field)
}

// GenerateChannelPath will generate a channel path with the given channel.
func GenerateChannelPath(ch asyncapi.Channel) string {
	// If there is no parameter, then just return the path
	if ch.Parameters == nil {
		return fmt.Sprintf("%q", ch.Path)
	}

	parameterRegexp := regexp.MustCompile("{[^{}]*}")

	matches := parameterRegexp.FindAllString(ch.Path, -1)
	format := parameterRegexp.ReplaceAllString(ch.Path, "%v")

	sprint := fmt.Sprintf("fmt.Sprintf(%q, ", format)
	for _, m := range matches {
		sprint += fmt.Sprintf("params.%s,", templateutil.Namify(m))
	}

	return sprint[:len(sprint)-1] + ")"
}

// OperationName returns `operationId` value from Publish or Subscribe operation if any.
// If no `operationID` exists — return provided default value (`name`).
func OperationName(channel asyncapi.Channel) string {
	var name string

	switch {
	case channel.Publish != nil && channel.Publish.OperationID != "":
		name = channel.Publish.OperationID
	case channel.Subscribe != nil && channel.Subscribe.OperationID != "":
		name = channel.Subscribe.OperationID
	default:
		name = channel.Name
	}

	return templateutil.Namify(name)
}

// HelpersFunctions returns the functions that can be used as helpers
// in a golang template.
func HelpersFunctions() template.FuncMap {
	return template.FuncMap{
		"getChildrenObjectSchemas":       GetChildrenObjectSchemas,
		"channelToMessage":               ChannelToMessage,
		"isRequired":                     IsRequired,
		"generateChannelPath":            GenerateChannelPath,
		"referenceToStructAttributePath": ReferenceToStructAttributePath,
		"operationName":                  OperationName,
		"referenceToTypeName":            ReferenceToTypeName,
		"generateValidateTags":           generators.GenerateValidateTags[asyncapi.Schema],
	}
}
